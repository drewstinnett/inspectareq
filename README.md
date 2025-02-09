# Inspect a Req(uest)

This is meant to make debugging http requests in Go easier by printing out curl,
httpie, or possibly other types of shell commands that can be used to debug.

Ever working on an app that's making lots of web calls, and sometimes you just
wanna run curl or httpie to try to recreate whatever your app is doing? Just do
`export DEBUG_CURL=1` and build it right in to your app!

## Usage

```bash
❯ go run ./examples/enable-env/
❯ DEBUG_CURL=1 go run ./examples/enable-env/
curl -X POST -H 'Authorization: REDACTED' -H 'Content-Type: application/json' -H 'X-Debug: true' 'https://pie.dev/anything' -d '{"username": "alice", "password": "secret"}'
❯ DEBUG_HTTPIE=1 go run ./examples/enable-env/
http POST 'https://pie.dev/anything' 'Authorization:REDACTED' 'Content-Type:application/json' 'X-Debug:true' data='{"username": "alice", "password": "secret"}'
❯ DEBUG_CURL=1 DEBUG_HTTPIE=1 go run ./examples/enable-env/
curl -X POST -H 'Authorization: REDACTED' -H 'Content-Type: application/json' -H 'X-Debug: true' 'https://pie.dev/anything' -d '{"username": "alice", "password": "secret"}'
http POST 'https://pie.dev/anything' 'Authorization:REDACTED' 'Content-Type:application/json' 'X-Debug:true' data='{"username": "alice", "password": "secret"}'
```
