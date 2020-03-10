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

LOG_FORMAT=bluecoat
#LOG_FORMAT=apache_common

ALOG=../output/hec

BYTES_PER_LINE=()
# added 1 byte more to accommodate for multiple every entry
BYTES_PER_LINE[apache_common]=106
BYTES_PER_LINE[apache_combined]=233
BYTES_PER_LINE[bluecoat]=1250
# Changed to 234 after updating random functions

BYTES_PER_GB=1024*1024*1024
NS_PER_SEC=1000*1000*1000
SEC_PER_DAY=24*60*60

# computed constants
LINES_PER_DAY=$(( GB_PER_DAY * BYTES_PER_GB / BYTES_PER_LINE[$LOG_FORMAT] ))
#LINES_PER_DAY=$(( LINES_PER_DAY / 1000000 )) # FIXME
LINES_PER_WORKER=$(( LINES_PER_DAY / WORKERS ))
TIME_PER_LINE=$(( SEC_PER_DAY * NS_PER_SEC / LINES_PER_DAY ))
TIME_PER_WORKER=$(( SEC_PER_DAY / WORKERS ))

START_DAY=$(date --date="today" +"%m/%d/%Y")

gen_day() {
    DAY=$1
    START_TS=$(date --date="$START_DAY -$DAY day" +"%s")
    echo "Generating day $DAY -- started at `date`"

    JOBS=()
    for W in $(seq 1 $WORKERS); do
        WIDX=$(( 1 + (W-1) * LINES_PER_WORKER ))
        WSTART=$(( START_TS + TIME_PER_WORKER * (W-1) ))
        $ALOG -f $LOG_FORMAT -c $WSTART -n $LINES_PER_WORKER -i $WIDX -s $TIME_PER_LINE -b $BATCH -u $SERVER -t $TOKEN -r $INDEX &
        JOBS[$W]=$!
    done

    for PID in ${JOBS[*]}; do
        echo "waiting for job $PID"
        wait $PID
    done
}

for DAY in $(seq $DAYS_TO_GEN -1 1); do
    gen_day $DAY
done
