# HTTP Check

[![Build Status](https://api.travis-ci.org/ggiamarchi/http-check.png?branch=master)](https://travis-ci.org/ggiamarchi/http-check)

HTTP Check is a simple tool that exposes any system command with a single HTTP endpoint. HTTP response code
differs depending command status code.

## Configuration

Here is an sample configuration file :

```yaml
---

checks:
  - name: "ping-10.1.1.2"
    command:
      executable: "/sbin/ping %s %d %s"
      args: ["-c", 2, "10.100.0.1"]
    status:
      failure: 223
      success: 200
  - name: "ping-google"
    command:
      executable: "/sbin/ping %s %d %s"
      args: ["-c", 2, "8.8.8.8"]
    status:
      failure: 503
      success: 204

server:
  port: 2323
```

## Installation

1. Download http-check binary from github release
2. Create configuration file. Default location is `/etc/http-check/http-check.yml`. You can specificy custom location on the command line when running the server

## Run the server

Basically run the binary

```
$ http-check server
```

Or, with a custom configuration file location

```
$ http-check server --config /etc/http-check.yml
```

## Run checks

Clients use a single endpoint to run any check.

__Request__

```
GET /v1/check/{name}
```

__Response__

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 29 Aug 2018 13:18:13 GMT
Content-Length: 350

{
    "error":"",
    "stderr":"",
    "stdout":"PING 8.8.8.8 (8.8.8.8): 56 data bytes\n64 bytes from 8.8.8.8: icmp_seq=0 ttl=121 time=12.676 ms\n64 bytes from 8.8.8.8: icmp_seq=1 ttl=121 time=20.267 ms\n\n--- 8.8.8.8 ping statistics ---\n2 packets transmitted, 2 packets received, 0.0% packet loss\nround-trip min/avg/max/stddev = 12.676/16.471/20.267/3.796 ms\n"
}
```

__NB.__ Status code depends the YAML check configuration.

# License

Everything in this repository is published under the MIT license.
