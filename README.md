# Ping Machine

> Simple HTTP server that allows to ping Internet adresses

```
# Build it
$ go build -o ping-yeah main.go
# Allow it to ping
$ sudo setcap cap_net_raw+ep ping-yeah
# Run it
$ ./ping-yeah
```
