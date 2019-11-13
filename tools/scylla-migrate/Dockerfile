
ARG drepo=eu.gcr.io/noomi-vnext-ci
FROM ubuntu:16.04

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app
COPY scylla-migrate /app

CMD [ "/bin/bash", "-c", "/app/scylla-migrate -kafka_url=confluent-platform-cp-kafka-headless.default.svc.cluster.local:9092 -cassandra_hosts=scylladb.default.svc.cluster.local"]

