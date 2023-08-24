#!/usr/bin/env sh

TARGET_BROKER=localhost:9092
TARGET_TOPIC=word
TARGET_TOPIC_JOIN=word-type
JOIN_INTERMEDIATE_TOPIC=word-join-intermediate
JOIN_OUTPUT_TOPIC=word-join

kafka-topics --bootstrap-server $TARGET_BROKER --delete --topic $TARGET_TOPIC
kafka-topics --bootstrap-server $TARGET_BROKER --delete --topic $TARGET_TOPIC_JOIN
kafka-topics --bootstrap-server $TARGET_BROKER --delete --topic $TARGET_TOPIC-count
kafka-topics --bootstrap-server $TARGET_BROKER --delete --topic $TARGET_TOPIC-updated
kafka-topics --bootstrap-server $TARGET_BROKER --delete --topic $TARGET_TOPIC-collect
kafka-topics --bootstrap-server $TARGET_BROKER --delete --topic $JOIN_INTERMEDIATE_TOPIC
kafka-topics --bootstrap-server $TARGET_BROKER --delete --topic $JOIN_OUTPUT_TOPIC

kafka-topics --bootstrap-server $TARGET_BROKER --create --partitions 10 --topic $TARGET_TOPIC
kafka-topics --bootstrap-server $TARGET_BROKER --create --partitions 10 --topic $TARGET_TOPIC_JOIN
kafka-topics --bootstrap-server $TARGET_BROKER --create --partitions 10 --topic $TARGET_TOPIC-count
kafka-topics --bootstrap-server $TARGET_BROKER --create --partitions 10 --topic $TARGET_TOPIC-updated
kafka-topics --bootstrap-server $TARGET_BROKER --create --partitions 10 --topic $TARGET_TOPIC-collect
kafka-topics --bootstrap-server $TARGET_BROKER --create --partitions 10 --topic $JOIN_INTERMEDIATE_TOPIC
kafka-topics --bootstrap-server $TARGET_BROKER --create --partitions 10 --topic $JOIN_OUTPUT_TOPIC

# 10 apples
# 5 pizzas
# 2 beers

echo "apple:apple 1" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
echo "apple:apple 2" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
echo "apple:apple 3" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
echo "apple:apple 4" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
echo "apple:apple 5" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P

echo "pizza:pizza 1" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
echo "pizza:pizza 2" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P

echo "beer:beer 1" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P

echo "apple:apple 6" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
echo "apple:apple 7" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
echo "apple:apple 8" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
echo "apple:apple 9" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
echo "apple:apple 10" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P

echo "pizza:pizza 3" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
echo "pizza:pizza 4" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P

echo "beer:beer 2" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P

echo "pizza:pizza 5" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P

echo "apple:fruit" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC_JOIN -Z -K: -P
echo "pizza:carbs" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC_JOIN -Z -K: -P
echo "beer:drink" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC_JOIN -Z -K: -P

PGPASSWORD=postgres psql -U postgres -h localhost -p 5432 -d postgres -f script/create.sql