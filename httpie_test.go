package inspectareq

import (
	"bytes"
	"net/http"
	"testing"
)

func TestHTTPie(t *testing.T) {
	for desc, tt := range map[string]struct {
		given  *http.Request
		expect string
	}{
		"simple-get": {
			given:  mustNewRequest("GET", "http://www.example.com", nil, nil),
			expect: "http GET 'http://www.example.com'\n",
		},
		"headers": {
			given: mustNewRequest("GET", "http://www.example.com", nil, map[string][]string{
				"Foo": {"Bar"},
				"Baz": {"Bazinga"},
			}),
			expect: "http GET 'http://www.example.com' 'Baz:Bazinga' 'Foo:Bar'\n",
		},
		"post-body": {
			given: mustNewRequest("POST", "http://www.example.com", bytes.NewReader([]byte("foo bar")), nil),
			expect: `http POST 'http://www.example.com' data='foo bar'
`,
		},
		"post-body-shell-escape": {
			given: mustNewRequest("POST", "http://www.example.com", bytes.NewReader([]byte("foo ' bar")), nil),
			expect: `http POST 'http://www.example.com' data='foo '\'' bar'
`,
		},
	} {
		got, err := httpie{r: New()}.Print(tt.given)
		if err != nil {
			t.Fatalf("expected no error, but got: %v\n", err)
		}

		if got != tt.expect {
			t.Errorf("%v:\nexpected: %q\n but got: %q\n", desc, tt.expect, got)
		}
	}
}

func TestHTTPieErr(t *testing.T) {
	if err, want := New(WithHTTPie()).Print(mustNewRequest("POST", "http://www.example.com", badBody{}, nil)), "error reading"; err.Error() != want {
		t.Errorf("expected error: %q, but got %q", want, err)
	}
}
