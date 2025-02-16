package inspectareq

import (
	"net/http"
	"strings"
)

type curl struct {
	r *Runner
}

// CurlEnv is the environment variable key to enable debugging with curl.
const CurlEnv string = "DEBUG_CURL"

// curl returns a curl compatible command string.
func (c curl) Print(req *http.Request) (string, error) {
	cmd := []string{
		"curl",
		"-X", req.Method,
	}

	for _, v := range c.r.headers(req.Header) {
		cmd = append(cmd, "-H", bashEscape(v[0]+": "+v[1]))
	}

	cmd = append(cmd, bashEscape(req.URL.String()))

	if req.Body != nil {
		body, err := readBody(&req.Body)
		if err != nil {
			return "", err
		}
		cmd = append(cmd, "-d", bashEscape(body))
	}

	return strings.Join(cmd, " ") + "\n", nil
}

// WithCurl enables the curl debugger.
func WithCurl() Option {
	return func(r *Runner) {
		r.debuggers = append(r.debuggers, curl{r: r})
	}
}
