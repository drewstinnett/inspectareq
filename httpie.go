package inspectareq

import (
	"fmt"
	"net/http"
	"strings"
)

// HTTPieEnv is the environment variable to enable httpie debugging
const HTTPieEnv string = "DEBUG_HTTPIE"

type httpie struct {
	r *Runner
}

// httpie returns an httpie compatible command string
func (h httpie) Print(req *http.Request) (string, error) {
	cmd := []string{
		"http",
		req.Method,
		bashEscape(req.URL.String()),
	}
	for _, v := range h.r.headers(req.Header) {
		cmd = append(cmd, bashEscape(v[0]+":"+v[1]))
	}
	if req.Body != nil {
		body, err := readBody(&req.Body)
		if err != nil {
			return "", err
		}
		cmd = append(cmd, fmt.Sprintf("data=%v", bashEscape(body)))
	}

	return strings.Join(cmd, " ") + "\n", nil
}

// WithHTTPie enables the httpie debugger
func WithHTTPie() Option {
	return func(r *Runner) {
		r.debuggers = append(r.debuggers, httpie{r: r})
	}
}
