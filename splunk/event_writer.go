package splunk

import (
	"sync"
	"time"
)

// Writer is a threadsafe, aysnchronous splunk writer.
// It implements io.Writer for usage in logging libraries, or whatever you want to send to splunk :)
// Writer.Client's configuration determines what source, sourcetype & index will be used for events
// Example for logrus:
//    splunkWriter := &splunk.Writer {Client: client}
//    logrus.SetOutput(io.MultiWriter(os.Stdout, splunkWriter))
type EventWriter struct {
	Client     *Client
	BufferSize int
	// How often the write buffer should be flushed to splunk
	FlushInterval time.Duration
	// How many Write()'s before buffer should be flushed to splunk
	FlushThreshold int
	// Max number of retries we should do when we flush the buffer
	MaxRetries int
	dataChan   chan *eventMessage
	errors     chan error
	once       sync.Once
}

type eventMessage struct {
	event Event
}

// Writer asynchronously writes to splunk in batches
func (w *EventWriter) Write(e Event) (int, error) {
	// only initialize once. Keep all of our buffering in one thread
	w.once.Do(func() {
		// synchronously set up dataChan
		w.dataChan = make(chan *eventMessage, w.BufferSize)
		// Spin up single goroutine to listen to our writes
		w.errors = make(chan error, w.BufferSize)
		go w.listen()
	})
	// Send the data to the channel
	w.dataChan <- &eventMessage{
		event: e,
	}
	// We don't know if we've hit any errors yet, so just say we're good
	return 1, nil
}

// Errors returns a buffered channel of errors. Might be filled over time, might not
// Useful if you want to record any errors hit when sending data to splunk
func (w *EventWriter) Errors() <-chan error {
	return w.errors
}

// listen for messages
func (w *EventWriter) listen() {
	if w.FlushInterval <= 0 {
		w.FlushInterval = defaultInterval
	}
	if w.FlushThreshold == 0 {
		w.FlushThreshold = defaultThreshold
	}
	ticker := time.NewTicker(w.FlushInterval)
	buffer := make([]*eventMessage, 0)
	//Define function so we can flush in several places
	flush := func() {
		// Go send the data to splunk
		go w.send(buffer, w.MaxRetries)
		// Make a new array since the old one is getting used by the splunk client now
		buffer = make([]*eventMessage, 0)
	}
	for {
		select {
		case <-ticker.C:
			if len(buffer) > 0 {
				flush()
			}
		case e := <-w.dataChan:
			buffer = append(buffer, e)
			if len(buffer) > w.FlushThreshold {
				flush()
			}
		}
	}
}

// send sends data to splunk, retrying upon failure
func (w *EventWriter) send(messages []*eventMessage, retries int) {
	// Create events from our data so we can send them to splunk
	events := make([]*Event, len(messages))
	for i, m := range messages {
		// Use the configuration of the Client for the event
		events[i] = &m.event
	}
	// Send the events to splunk
	err := w.Client.LogEvents(events)
	// If we had any failures, retry as many times as they requested
	if err != nil {
		for i := 0; i < retries; i++ {
			// retry
			err = w.Client.LogEvents(events)
			if err == nil {
				return
			}
		}
		// if we've exhausted our max retries, let someone know via Errors()
		// might not have retried if retries == 0
		select {
		case w.errors <- err:
		// Don't block in case no one is listening or our errors channel is full
		default:
		}
	}
}
