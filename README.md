# serverlog

Small utility server to print out the request path and body

Useful to quickly tests webhooks

## Quick start

```
git clone git@github.com:plopezlpz/serverlog.git
cd serverlog
make build
./bin/serverlog -p 8080 -r '{"response": "ok"}'
```
