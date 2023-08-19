#!/usr/bin/env sh

TARGET_BROKER=localhost:9092

kafka-consumer-groups --bootstrap-server $TARGET_BROKER --delete --group "flows-word-count"
kafka-consumer-groups --bootstrap-server $TARGET_BROKER --delete --group "flows-word-remap"
kafka-consumer-groups --bootstrap-server $TARGET_BROKER --delete --group "flows-word-materialise"
kafka-consumer-groups --bootstrap-server $TARGET_BROKER --delete --group "flows-word-join"
