#!/bin/bash

if [ "$#" -ne 7 ]; then
    echo "Usage: $0 gb_per_day days workers batchsize server token index"
    exit 1
fi

GB_PER_DAY=$1
DAYS_TO_GEN=$2
WORKERS=$3
BATCH=$4
SERVER=$5
TOKEN=$6
INDEX=$7

SERVERS=()
SERVERS+=("server1")
SERVERS+=("server2")
SERVERS+=("server3")

for server in ${SERVERS[@]}; do
        ./gen_hec.sh $GB_PER_DAY $DAYS_TO_GEN $WORKERS $BATCH $server $TOKEN $INDEX &
done