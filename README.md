# apclog
Golang Eventgen for Splunk HEC

* forked from https://github.com/rsomu/apclog
* based on https://github.com/mingrammer/flog

## Usage
### Build
```
make
```

### Backfill
```
cd utils
./hec.sh 1 1 1 100 servername {uuid-token} index
```

### Continuous Real-time
```
cd output
./hec -f format -x -n #eps
```
