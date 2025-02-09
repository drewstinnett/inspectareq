package inspectareq

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestBodyRead(t *testing.T) {
	req := mustNewRequest("GET", "https://www.example.com", bytes.NewReader([]byte("hello world")), nil)
	h := httpie{}
	if _, err := h.Print(req); err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}
	if string(body) == "" {
		t.Errorf("body was emptied out by HTTPPie")
	}
}

func TestEnvs(t *testing.T) {
	for desc, tt := range map[string]struct {
		envs   map[string]string
		given  *http.Request
		expect string
	}{
		"httpie": {
			envs:   map[string]string{HTTPieEnv: "1"},
			given:  mustNewRequest("GET", "http://www.example.com", nil, nil),
			expect: "http GET 'http://www.example.com'\n",
		},
		"curl": {
			envs:   map[string]string{CurlEnv: "1"},
			given:  mustNewRequest("GET", "http://www.example.com", nil, nil),
			expect: "curl -X GET 'http://www.example.com'\n",
		},
		"httpie-and-curl": {
			envs: map[string]string{
				HTTPieEnv: "1",
				CurlEnv:   "1",
			},
			given:  mustNewRequest("GET", "http://www.example.com", nil, nil),
			expect: "curl -X GET 'http://www.example.com'\nhttp GET 'http://www.example.com'\n",
		},
	} {
		os.Clearenv()
		for k, v := range tt.envs {
			t.Setenv(k, v)
		}

		// Test a custom Runner
		b := bytes.NewBufferString("")
		c := New(WithWriter(b), WithEnvironment())
		// Set the default client to the new one we just created
		Set(c)
		if err := c.Print(tt.given); err != nil {
			t.Fatalf("got unexpected error: %v", err)
		}

		if b.String() != tt.expect {
			t.Errorf("%v:\nexpected: %v\n but got: %v\n", desc, tt.expect, b.String())
		}

		// Test the default debugger
		b = bytes.NewBufferString("")
		c = Get()
		c.writer = b
		// c.Enable() // The envs are reset so just enable it here
		if err := c.Print(tt.given); err != nil {
			t.Fatalf("\n%v:\ngot unexpected error: %v", desc, err)
		}

		if b.String() != tt.expect {
			t.Errorf("\n%v:\nexpected: %v\n but got: %v\n", desc, tt.expect, b.String())
		}
	}
}

func TestEnableDisable(t *testing.T) {
	b := bytes.NewBufferString("")
	d := New(WithWriter(b), WithDebugger(httpie{}))

	req := mustNewRequest("GET", "http://www.example.com", nil, nil)

	if err := d.Print(req); err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	if b.String() != "" {
		t.Error("got a non empty string back from a fresh debugger: " + b.String())
	}

	// Enable printing
	d.Enable()
	if err := d.Print(req); err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	expect := "http GET 'http://www.example.com'\n"
	if b.String() != expect {
		t.Errorf("\nexpected: %v\n but got: %v", expect, b.String())
	}

	// Disable printing
	d.Disable()
	b = bytes.NewBufferString("")
	if err := d.Print(req); err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	if b.String() != "" {
		t.Error("got a non empty string back from a fresh debugger: " + b.String())
	}
}

func TestHeaders(t *testing.T) {
	given := http.Header{}
	given.Add("Key2", "Value1")
	given.Add("Key1", "Value1")
	given.Add("Key1", "Value2")
	r := New()
	got := r.headers(given)

	if len(got) != 3 {
		t.Errorf("expected length of 3 but got: %v", len(got))
	}
	if got[1][0] != "Key1" {
		t.Errorf("expected key of: 'Key1', but got: %v", got[1][0])
	}
	if got[1][1] != "Value2" {
		t.Errorf("expected value of: 'Value2', but got: %v", got[1][1])
	}
}

// Panic if given an invalid request. Really try hard not do that here in the tests
func mustNewRequest(method string, url string, body io.Reader, headers map[string][]string) *http.Request {
	req, err := http.NewRequest(method, url, body)
	for k, v := range headers {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	if err != nil {
		panic(err)
	}
	return req
}
