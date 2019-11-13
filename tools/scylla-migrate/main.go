package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/AiflooAB/cloud-engine/pkg/actor/persistence"
	"github.com/AiflooAB/cloud-engine/pkg/actor/persistence/pkafka"
	"github.com/AiflooAB/cloud-engine/pkg/log"
	"github.com/AiflooAB/cloud-engine/pkg/messaging/kafkaclient"
	"github.com/Shopify/sarama"
	"github.com/gocql/gocql"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

var (
	kafkaUrl       string
	cassandraHosts string
	logger         log.Logger
	serviceName    = "scylla-migration"
)

func initFlags() {

	flag.StringVar(&kafkaUrl, "kafka_url", "localhost:9092", "dns address of kafka broker")
	flag.StringVar(&cassandraHosts, "cassandra_hosts", "localhost", "dns address of cassandra")

	flag.Parse()
}

func sendKafkaSnapshot(writer kafkaclient.KafkaSyncProducer, topic string, snap *persistence.Snapshot, logger log.Logger) {
	snapEnvelopePayload, err := proto.Marshal(snap)
	if err != nil {
		logger.Error("KafkaPersistenceProvider marshal err", zap.Error(err))
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(snap.Id),
		Value: sarama.ByteEncoder(snapEnvelopePayload),
	}
	_, _, err = writer.SendMessage(msg)
	if err != nil {
		logger.Error("KafkaPersistenceProvider send err", zap.Error(err))
	}
}

func main() {
	logger = log.GetFactory().Bg()
	logger.Info("Starting scylla-migrate..")

	initFlags()

	writerCfg := sarama.NewConfig()
	writerCfg.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	writerCfg.Producer.Retry.Max = 20                   // Retry up to 10 times to produce the message
	writerCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{kafkaUrl}, writerCfg)
	if err != nil {
		logger.Fatal("Failed to start Sarama producer:", zap.Error(err))
	}

	c := gocql.NewCluster(cassandraHosts)

	c.Timeout = 15 * time.Second
	c.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
		NumRetries: 50,
		Min:        50 * time.Millisecond,
		Max:        3 * time.Second,
	}

	s, err := c.CreateSession()
	if err != nil {
		panic(err)
	}
	defer s.Close()
	sinner, err := c.CreateSession()
	if err != nil {
		panic(err)
	}
	defer sinner.Close()

	idmap := make(map[string]int)
	pIDIter := s.Query(`select distinct persistence_id from journal.snapshots`).Iter()
	var pID string
	for pIDIter.Scan(&pID) {

		iter := sinner.Query(`SELECT sequence_nr, snapshot FROM journal.snapshots WHERE persistence_id = ?`, pID).Iter()
		var payload []byte
		var seqnr int
		var snp *persistence.Snapshot

		for iter.Scan(&seqnr, &payload) {
			//can be nil if snapshot is does not exist anymore
			if &payload != nil && len(payload) > 0 {
				if snp != nil && snp.EventIndex >= int32(seqnr) {
					continue
				} else {
					snp = &persistence.Snapshot{
						Id:         pID,
						EventIndex: int32(seqnr),
						Payload:    payload,
					}
					idmap[pID] = seqnr
				}
				if snp != nil {
					//logger.Info(fmt.Sprintf("Processing %s, seq: %d", pId, ver))
					if strings.HasPrefix(pID, "wbDeviceActor") || strings.HasPrefix(pID, "clipActor") {
						sendKafkaSnapshot(producer, "airstream-table", snp, logger)
					} else {
						sendKafkaSnapshot(producer, "sherlock-table", snp, logger)
					}
				}
			}
		}
		if err = iter.Close(); err != nil {
			logger.Warn("GetSnapshot", zap.Error(err))
		}
	}
	if err := pIDIter.Close(); err != nil {
		logger.Warn("GetDistinct", zap.Error(err))
	}
	producer.Close()

	logger.Info(fmt.Sprintf("Processed %d snapshot", len(idmap)))

	//compose table of persistent_id - highest_version and check it after migration

	logger.Info("MIGRATION DONE")

	airstreamProvider := pkafka.NewKafkaPersistenceProvider(kafkaUrl, "airstream-table", "migration")
	defer airstreamProvider.Close()
	sherlockProvider := pkafka.NewKafkaPersistenceProvider(kafkaUrl, "sherlock-table", "migrations")
	defer sherlockProvider.Close()

	time.Sleep(5 * time.Second)

	for pID, ver := range idmap {
		if strings.HasPrefix(pID, "wbDeviceActor") || strings.HasPrefix(pID, "clipActor") {
			_, rev, ok := airstreamProvider.GetSnapshot(pID)
			if !ok {
				logger.Info(fmt.Sprintf("No snapshot for %s", pID))
			}
			if rev != ver {
				logger.Info(fmt.Sprintf("Snapshot version mismatch %s, seq: %d has: %d", pID, ver, rev))
			}
		} else {
			_, rev, ok := sherlockProvider.GetSnapshot(pID)
			if !ok {
				logger.Info(fmt.Sprintf("No snapshot for %s", pID))
			}
			if rev != ver {
				logger.Info(fmt.Sprintf("Snapshot version mismatch %s, seq: %d has: %d", pID, ver, rev))
			}
		}
	}

	logger.Info("VALIDATION DONE")

	// read all snapshots from scylla
	// and copy to the kafka-topics

}
