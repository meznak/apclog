module github.com/bdalpe/apclog

go 1.14

require (
	apclog/splunk v0.0.0-00010101000000-000000000000
	github.com/brianvoe/gofakeit v3.17.0+incompatible
	github.com/spf13/pflag v1.0.0
)

replace apclog/splunk => ./splunk
