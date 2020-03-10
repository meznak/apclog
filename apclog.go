package main

import (
	"apclog/splunk"
	"crypto/tls"
	"math"
	"net/http"
	"time"
)

var source = map[string]string{
	"apache_common":   "apache",
	"apache_combined": "apache",
	"bluecoat":        "bluecoat",
}

var sourcetype = map[string]string{
	"apache_common":   "apache_common",
	"apache_combined": "apache_combined",
	"bluecoat":        "bluecoat:proxysg:access:syslog",
}

// Generate generates the logs with given options
func Generate(option *Option) error {
	// Create new Splunk client
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	var httpClient = &http.Client{Timeout: time.Second * 20, Transport: tr}

	s := splunk.NewClient(
		httpClient,
		"https://"+option.Server+":8088/services/collector",
		option.Token,
		source[option.Format],
		sourcetype[option.Format],
		option.SplunkIndex,
	)

	w := splunk.EventWriter{
		Client:         s,
		FlushInterval:  60 * time.Second,
		FlushThreshold: int(option.Batch),
		MaxRetries:     2,
		BufferSize:     int(option.Batch),
	}

	frac_ns := 1e9 * (option.Created - int64(math.Floor(float64(option.Created))))
	created := time.Unix(int64(option.Created), int64(frac_ns))
	ival := time.Duration(option.Sleep)

	index := option.Index

	if option.Continuous {
		limiter := time.Tick(time.Duration(1000000000/option.Number) * time.Nanosecond)
		for {
			every := DeriveEvery(index)
			evt := NewLog(option.Format, time.Now(), every)
			e := s.NewEventWithTime(time.Now().UnixNano(), evt, source[option.Format], sourcetype[option.Format], s.Index)
			<-limiter
			_, _ = w.Write(*e)
			index++
		}
	} else {
		// Generates the logs until the certain number of lines is reached
		for line := 0; line < option.Number; line++ {
			every := DeriveEvery(index)
			evt := NewLog(option.Format, created, every)
			e := s.NewEventWithTime(created.Unix(), evt, source[option.Format], sourcetype[option.Format], s.Index)

			_, _ = w.Write(*e)

			created = created.Add(ival)
			index++
		}
	}

	return nil
}

// NewLog creates a log for given format
func NewLog(format string, t time.Time, every string) string {
	switch format {
	case "apache_common":
		return NewApacheCommonLog(t, every)
	case "apache_combined":
		return NewApacheCombinedLog(t, every)
	case "bluecoat":
		return NewBluecoatLog(t, every)
	default:
		return ""
	}
}

func DeriveEvery(count int64) string {
	if count%1e12 == 0 {
		return "every1t"
	} else if count%1e11 == 0 {
		return "every1 every10 every100 every1k every10k every10b every1m every10m every100m every1b every10b every100b"
	} else if count%1e10 == 0 {
		return "every1 every10 every100 every1k every10k every10b every1m every10m every100m every1b every10b"
	} else if count%1e9 == 0 {
		return "every1 every10 every100 every1k every10k every100k every1m every10m every100m every1b"
	} else if count%1e8 == 0 {
		return "every1 every10 every100 every1k every10k every100k every1m every10m every100m"
	} else if count%1e7 == 0 {
		return "every1 every10 every100 every1k every10k every100k every1m every10m"
	} else if count%1e6 == 0 {
		return "every1 every10 every100 every1k every10k every100k every1m"
	} else if count%1e5 == 0 {
		return "every1 every10 every100 every1k every10k every100k"
	} else if count%1e4 == 0 {
		return "every1 every10 every100 every1k every10k"
	} else if count%1e3 == 0 {
		return "every1 every10 every100 every1k"
	} else if count%1e2 == 0 {
		return "every1 every10 every100"
	} else if count%1e1 == 0 {
		return "every1 every10"
	} else {
		return "every1"
	}
}
