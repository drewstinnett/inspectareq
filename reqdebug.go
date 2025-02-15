/*
Package inspectareq is used to debug http.Request items
*/
package inspectareq

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"
)

// Runner is the container struct for things printin' out strings
type Runner struct {
	enabled       bool
	writer        io.Writer
	debuggers     DebuggerList
	redact        bool
	redactWith    string
	redactHeaders []string
	mu            *sync.RWMutex
}

// Debugger is an interface that defines how a debugger behaves
type Debugger interface {
	Print(*http.Request) (string, error)
}

// DebuggerList is a list of fmt.Stringer objects that print out the debug command to a string
type DebuggerList []Debugger

// NoRedactEnv is the environment variable to disable redactions
const NoRedactEnv string = "NO_REDACT"

// New returns a new Debugger using functional options
func New(opts ...Option) *Runner {
	d := &Runner{
		writer:        os.Stderr,
		redact:        true,
		redactWith:    "REDACTED",
		enabled:       true,
		redactHeaders: []string{"Authorization"},
		mu:            &sync.RWMutex{},
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Option is a functional option to the Debugger object
type Option func(*Runner)

// WithEnvironment allows the Debugger to look at environment variables when
// making decisions about how to format the output, and if to output at all
func WithEnvironment() Option {
	return func(r *Runner) {
		// Assume off, if using the environment to configure
		r.enabled = false

		if os.Getenv(CurlEnv) != "" {
			r.enabled = true
			r.debuggers = append(r.debuggers, curl{r: r})
		}
		if os.Getenv(HTTPieEnv) != "" {
			r.enabled = true
			r.debuggers = append(r.debuggers, httpie{r: r})
		}
		if os.Getenv(NoRedactEnv) != "" {
			r.redact = false
		}
	}
}

// WithoutRedact turns redaction off
func WithoutRedact() Option {
	return func(r *Runner) {
		r.redact = false
	}
}

// WithDebugger enables the httpie debugger
func WithDebugger(d Debugger) Option {
	return func(r *Runner) {
		r.debuggers = append(r.debuggers, d)
	}
}

// WithWriter sets the writer for the Debugger client
func WithWriter(w io.Writer) Option {
	return func(r *Runner) {
		r.writer = w
	}
}

// Enable enables the default debug printer
func Enable() { defaultR.Enable() }

// Enable enables the debug printer
func (r *Runner) Enable() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.enabled = true
}

// Disable diables the default debug printer
func Disable() { defaultR.Disable() }

// Disable disables the debug printer
func (r *Runner) Disable() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.enabled = false
}

// Enabled returns the state of the default Debugger
func Enabled() bool { return defaultR.Enabled() }

// Enabled returns the enabled state
func (r Runner) Enabled() bool {
	return r.enabled
}

// Print prints using the default Debugger
func Print(req *http.Request) error { return defaultR.Print(req) }

// Print will print the req statement out to the writer, if enabled
func (r Runner) Print(req *http.Request) error {
	if !r.enabled {
		return nil
	}
	for _, m := range r.debuggers {
		got, err := m.Print(req)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprint(r.writer, got); err != nil {
			return err
		}
	}
	return nil
}

func headerKeys(h http.Header) []string {
	keys := make([]string, 0, len(h)) // Use capacity to avoid unnecessary allocations
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// bashEscape takes a string and returns a safely escaped version for use in Bash.
func bashEscape(input string) string {
	// Single quotes prevent Bash from interpreting special characters.
	// If the input contains single quotes, we need to escape them.
	if strings.Contains(input, "'") {
		// Replace ' with ''' (closes, escapes, and reopens single-quoted string)
		return "'" + strings.ReplaceAll(input, "'", "'\\''") + "'"
	}
	return "'" + input + "'"
}

// readBody reads an io.ReadCloser and immediately puts the original data back
// in place, so it can be used again later. A copy of the body is returned as a
// string
func readBody(b *io.ReadCloser) (string, error) {
	// https://blog.flexicondev.com/read-go-http-request-body-multiple-times#heading-the-simplest-solution
	body, err := io.ReadAll(*b)
	if err != nil {
		return "", err
	}
	// Replace the body with a new reader after reading from the original
	*b = io.NopCloser(bytes.NewBuffer(body))
	return string(body), nil
}

// defaultR is the default Runner, created at init() using WithEnvironment()
var defaultR *Runner

func init() {
	defaultR = New(WithEnvironment())
}

// Get returns the default debugger
func Get() *Runner {
	return defaultR
}

// Set sets the given Debugger to the default
func Set(v *Runner) {
	defaultR = v
}

// headers returns the keys and values from an http.Header as a slice
func (r *Runner) headers(h http.Header) [][2]string {
	ret := [][2]string{}
	for _, k := range headerKeys(h) {
		values := []string{}
		for _, vv := range h[k] {
			headerValue := r.redactHeader(k, vv)
			values = append(values, headerValue)
		}
		sort.Strings(values)
		for _, value := range values {
			ret = append(ret, [2]string{k, value})
		}
	}
	return ret
}

// redact Headers replaces sensitive header values with the string in Runner.redactWith
func (r *Runner) redactHeader(k, v string) string {
	if !r.redact {
		return v
	}
	if slices.Contains(r.redactHeaders, k) {
		v = r.redactWith
	}
	return v
}
