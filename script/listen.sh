#!/usr/bin/env sh

TARGET_BROKER=localhost:9092
TARGET_TOPIC=produce-topic

kcat -b $TARGET_BROKER -t $TARGET_TOPIC -C -f '\nKey (%K string): %k Value (%S string): %s Timestamp: %T Partition: %p Offset: %o Headers: %h \n'
