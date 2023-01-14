#!/usr/bin/env sh

TARGET_BROKER=localhost:9092

kafka-consumer-groups.sh --bootstrap-server $TARGET_BROKER --delete --group "test"
