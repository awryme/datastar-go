package ds

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/awryme/datastar-go/bufpool"
)

// Options contain options for datastar url actions.
//
// Refer to datastar docs for details: https://data-star.dev/reference/action_plugins#options
type Options struct {
	// Query allows to set additional url query to send to server.
	Query map[string]any

	// Sets contentType to 'form'.
	Form bool

	// Specifies a form to send when Form == true.
	FormSelector string

	// Specifies whether to include local signals (those beggining with an underscore) in the request.
	IncludeLocal bool

	// An map containing headers to send with the request.
	Headers map[string]string

	// Specifies whether to keep the connection open when the page is hidden. Useful for dashboards but can cause a drain on battery life and other resources when enabled.
	KeepOpen bool

	// The retry interval.
	RetryInterval time.Duration

	// A numeric multiplier applied to scale retry wait times.
	RetryScaler int

	// The maximum allowable wait time between retries.
	RetryMaxWait time.Duration

	// The maximum number of retry attempts.
	RetryMaxCount int
}

// Sends a GET request to the backend using fetch, and merges the response with the current DOM and signals. The URL can be any valid URL and the response must contain zero or more Datastar SSE events.
func Get(url string, opts ...Options) string {
	return actionURL("get", url, mergeOpts(opts))
}

// Works the same as Get but sends a POST request to the backend.
func Post(url string, opts ...Options) string {
	return actionURL("post", url, mergeOpts(opts))
}

// Works the same as Get but sends a PUT request to the backend.
func Put(url string, opts ...Options) string {
	return actionURL("put", url, mergeOpts(opts))
}

// Works the same as Get but sends a PATCH request to the backend.
func Patch(url string, opts ...Options) string {
	return actionURL("patch", url, mergeOpts(opts))
}

// Works the same as Get but sends a DELETE request to the backend.
func Delete(url string, opts ...Options) string {
	return actionURL("delete", url, mergeOpts(opts))
}

func mergeOpts(opts []Options) *Options {
	if len(opts) == 0 {
		// fast path
		return nil
	}
	mergedOpts := &Options{}
	for _, opt := range opts {
		if len(opt.Query) > 0 {
			maps.Copy(mergedOpts.Query, opt.Query)
		}

		if len(opt.Headers) > 0 {
			maps.Copy(mergedOpts.Headers, opt.Headers)
		}

		setNotZero(opt.Form, &mergedOpts.Form)
		setNotZero(opt.FormSelector, &mergedOpts.FormSelector)
		setNotZero(opt.IncludeLocal, &mergedOpts.IncludeLocal)
		setNotZero(opt.KeepOpen, &mergedOpts.KeepOpen)
		setNotZero(opt.RetryInterval, &mergedOpts.RetryInterval)
		setNotZero(opt.RetryScaler, &mergedOpts.RetryScaler)
		setNotZero(opt.RetryMaxWait, &mergedOpts.RetryMaxWait)
		setNotZero(opt.RetryMaxCount, &mergedOpts.RetryMaxCount)
	}
	return mergedOpts
}

func setNotZero[T comparable](value T, to *T) {
	var zero T
	if value == zero {
		return
	}
	if to == nil {
		return
	}
	*to = value
}

func actionURL(method string, url string, options *Options) string {
	if options == nil {
		// fast path
		return fmt.Sprintf("@%s(`%s`)", method, url)
	}
	dsOpts := make(map[string]string)
	if options.Form {
		dsOpts["contentType"] = "'form'"
	}
	if options.FormSelector != "" {
		dsOpts["selector"] = fmt.Sprintf("'%s'", options.FormSelector)
	}

	if options.IncludeLocal {
		dsOpts["includeLocal"] = "true"
	}

	// headers
	if len(options.Headers) > 0 {
		dsHeaders := make(map[string]string)
		for k, v := range options.Headers {
			dsHeaders[k] = fmt.Sprintf("'%s'", v)
		}
		dsOpts["headers"] = mapToJs(dsHeaders)
	}

	if options.KeepOpen {
		dsOpts["openWhenHidden"] = "true"
	}

	if options.RetryInterval != 0 {
		dsOpts["retryInterval"] = fmt.Sprint(options.RetryInterval.Milliseconds())
	}

	if options.RetryScaler != 0 {
		dsOpts["retryScaler"] = fmt.Sprint(options.RetryScaler)
	}

	if options.RetryMaxWait != 0 {
		dsOpts["retryMaxWaitMs"] = fmt.Sprint(options.RetryMaxWait.Milliseconds())
	}

	if options.RetryMaxCount != 0 {
		dsOpts["retryMaxCount"] = fmt.Sprint(options.RetryMaxCount)
	}

	jsopts := mapToJs(dsOpts)

	return fmt.Sprintf("@%s('%s%s', %s)", method, url, buildQuery(options.Query), jsopts)
}

func buildQuery(query map[string]any) string {
	if len(query) == 0 {
		// fast path
		return ""
	}
	buf := bufpool.GetBuffer()
	defer bufpool.PutBuffer(buf)

	buf.WriteString("?")
	for key, value := range query {
		fmt.Fprintf(buf, "%s=%v&", key, value)
	}
	qs := buf.String()
	return strings.TrimSuffix(qs, "&")
}

func mapToJs(m map[string]string) string {
	if len(m) == 0 {
		// fast path
		return ""
	}

	buf := bufpool.GetBuffer()
	defer bufpool.PutBuffer(buf)

	buf.WriteString("{")
	for idx, name := range slices.Sorted(maps.Keys(m)) {
		val := m[name]
		fmt.Fprintf(buf, "%s: %s", name, val)
		if idx < len(m)-1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteString("}")
	return buf.String()
}
