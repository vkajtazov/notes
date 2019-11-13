kubectl exec -it zebra-platform-cp-kafka-0 -c cp-kafka-broker -- /usr/bin/kafka-topics --zookeeper $KAFKA_ZOOKEEPER_CONNECT --create --topic airstream-table --partitions 1 --replication-factor 2 -config min.insync.replicas=1 --config cleanup.policy=compact --config segment.ms=1000 --config delete.retention.ms=5000
kubectl exec -it zebra-platform-cp-kafka-0 -c cp-kafka-broker -- /usr/bin/kafka-topics --zookeeper $KAFKA_ZOOKEEPER_CONNECT --create --topic sherlock-table --partitions 1 --replication-factor 2 -config min.insync.replicas=1 --config cleanup.policy=compact --config segment.ms=1000 --config delete.retention.ms=5000

kafka-topics --zookeeper $KAFKA_ZOOKEEPER_CONNECT --topic airstream-table --describe