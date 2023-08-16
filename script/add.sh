#!/usr/bin/env sh

TARGET_BROKER=localhost:9092
TARGET_TOPIC=word
TARGET_TOPIC_JOIN=word-type
JOIN_INTERMEDIATE_TOPIC=word-join-intermediate
JOIN_OUTPUT_TOPIC=word-join

# 10 apples
# 5 pizzas
# 2 beers

# echo "chicken:chicken 1" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
# echo "apple:apple 2" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
# echo "pizza:pizza 2" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
# echo "beer:beer 1" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
# echo "ramen:ramen 1" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
# echo "sushi:sushi 1" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P
# echo "burger:burger 1" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC -Z -K: -P

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

echo "apple:great fruit" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC_JOIN -Z -K: -P
echo "pizza:dirty carbs" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC_JOIN -Z -K: -P
echo "beer:alcoholic drink" | kcat -b $TARGET_BROKER -t $TARGET_TOPIC_JOIN -Z -K: -P
