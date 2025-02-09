package inspectareq

import (
	"bytes"
	"net/http"
	"testing"
)

func TestCurl(t *testing.T) {
	for desc, tt := range map[string]struct {
		given  *http.Request
		expect string
	}{
		"simple-get": {
			given:  mustNewRequest("GET", "http://www.example.com", nil, nil),
			expect: "curl -X GET 'http://www.example.com'\n",
		},
		"headers": {
			given: mustNewRequest("GET", "http://www.example.com", nil, map[string][]string{
				"Foo":           {"Bar", "Bin"},
				"Baz":           {"Bazinga"},
				"Authorization": {"Bearer my-token"},
			}),
			expect: "curl -X GET -H 'Authorization: REDACTED' -H 'Baz: Bazinga' -H 'Foo: Bar' -H 'Foo: Bin' 'http://www.example.com'\n",
		},
		"post-body": {
			given: mustNewRequest("POST", "http://www.example.com", bytes.NewReader([]byte("foo bar")), nil),
			expect: `curl -X POST 'http://www.example.com' -d 'foo bar'
`,
		},
		"post-body-shell-escape": {
			given: mustNewRequest("POST", "http://www.example.com", bytes.NewReader([]byte("foo ' bar")), nil),
			expect: `curl -X POST 'http://www.example.com' -d 'foo '\'' bar'
`,
		},
	} {
		got, err := curl{r: New()}.Print(tt.given)
		if err != nil {
			t.Fatalf("expected no error, but got: %v\n", err)
		}

		if got != tt.expect {
			t.Errorf("%v:\nexpected: %v\n but got: %v\n", desc, tt.expect, got)
		}
	}
}
