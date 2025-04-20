package datastar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/awryme/datastar-go/bufpool"
	"github.com/awryme/sse-go/sseserver"
	"github.com/valyala/fastjson"
)

var parserPool fastjson.ParserPool

// Datastar is the main engine to handle datastart requests.
// It allows you to parse incoming signals or send events to client.
type Datastar struct {
	// Sets sse retry field
	// See https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#retry
	SSERetry time.Duration

	resp http.ResponseWriter
	req  *http.Request

	rawData    []byte
	jsonData   *fastjson.Value
	jsonParser *fastjson.Parser

	sse *sseserver.Server
}

// New creates a new Datastar instance.
// It uses fastjson.Parser to parse incoming signals.
//
// You should `defer release()` to reuse these parsers.
// If you don't - nothing will leak, but parsing signals will be less optimised.
func New(w http.ResponseWriter, r *http.Request) (ds *Datastar, release func()) {
	ds = &Datastar{
		resp: w,
		req:  r,
	}

	release = func() {
		if ds.jsonParser != nil {
			parserPool.Put(ds.jsonParser)
		}
	}

	return ds, release
}

// Request signals

// UnmarshalSignals unmarshals a signal (or multiple) into a provided value.
// It uses json.Unmarshal to do it, so regular Unmarshal rules apply.
// If path is provided it will find signal value at that path.
//
// Path can be separated by "." or be made of individual components, like: "my.data.value" or ("my", "data", "value").
func (ds *Datastar) UnmarshalSignals(value any, path ...string) error {
	// ensure we have at least raw data
	if err := ds.readRawData(); err != nil {
		return err
	}

	if len(path) == 0 {
		// fast path, just unmarshal the whole object
		return json.Unmarshal(ds.rawData, value)
	}

	// slower path, use fastjson to parse
	if err := ds.parseSignals(); err != nil {
		return err
	}

	fullpath := strings.Join(path, signalSeparator)
	keys := strings.Split(fullpath, signalSeparator)

	jsonValue := ds.jsonData.Get(keys...)
	if jsonValue == nil {
		return fmt.Errorf("signal %s not found", fullpath)
	}

	return json.Unmarshal(jsonValue.MarshalTo(nil), &value)
}

func (ds *Datastar) parseSignals() error {
	if ds.jsonData != nil {
		return nil
	}

	if err := ds.readRawData(); err != nil {
		return err
	}

	if ds.jsonParser == nil {
		ds.jsonParser = parserPool.Get()
	}

	var err error
	ds.jsonData, err = ds.jsonParser.ParseBytes(ds.rawData)
	return err
}

func (ds *Datastar) readRawData() (err error) {
	if ds.rawData != nil {
		return nil
	}

	if ds.req.Method == http.MethodGet {
		query := ds.req.URL.Query()

		data := query.Get("datastar")
		if data == "" {
			return fmt.Errorf("datastar query signals not found")
		}

		ds.rawData = []byte(data)
		return nil
	}

	ds.rawData, err = io.ReadAll(ds.req.Body)
	if err != nil {
		return fmt.Errorf("read request body signals: %w", err)
	}

	return nil
}

// SSE

// Send sends a datastar event to client.
// Events are created individually with respective functions or structs.
// Events are buffered, with reusable buffer pool.
func (ds *Datastar) Send(event Event) error {
	writer, release := ds.newEventWriter(event.Name())
	defer release()

	err := event.WriteEvent(writer, ds.req)
	if err != nil {
		return err
	}
	return ds.writeEvent(writer)
}

func (ds *Datastar) writeEvent(writer *sseserver.EventWriter) (err error) {
	if ds.sse != nil {
		return ds.sse.WriteEvent(writer)
	}

	ds.sse, err = sseserver.New(ds.resp, ds.req)
	if err != nil {
		return fmt.Errorf("make new sse server: %w", err)
	}

	return ds.sse.WriteEvent(writer)
}

func (ds *Datastar) newEventWriter(name string) (writer *sseserver.EventWriter, release func()) {
	buf := bufpool.GetBuffer()
	release = func() {
		bufpool.PutBuffer(buf)
	}

	writer = sseserver.NewEventWriter(buf, name, "", ds.SSERetry)
	return writer, release
}
