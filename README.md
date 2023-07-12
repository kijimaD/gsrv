# gsrv

```
$ go run cmd/main.go .
> GET /dummy.txt HTTP/1.0

HTTP/1.0 200 OK
Date: Wed, 12 Jul 2023 23:48:35 JST
Server: gsrv/1.0.0
Connection: close
Content-Length: 188
Content-Type: text/plain

this is dummy text

> GET /this_is_not_exists HTTP/1.0

HTTP/1.0 404 Not Found
Date: Wed, 12 Jul 2023 23:55:19 JST
Server: gsrv/1.0.0
Connection: close
Content-Type: text/html
```
