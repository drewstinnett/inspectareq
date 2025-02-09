# Inspect a Req(uest)

This is meant to make debugging http requests in Go easier by printing out curl,
httpie, or possibly other types of shell commands that can be used to debug.

Ever working on an app that's making lots of web calls, and sometimes you just
wanna run curl or httpie to try to recreate whatever your app is doing? Just do
`export DEBUG_CURL=1` and build it right in to your app!

## Usage

```bash
go run ./examples/enable-env/
...
export DEBUG_CURL=1
go run ./examples/enable-env/
...
export DEBUG_HTTPIE=1
go run ./examples/enable-env/
...
```
