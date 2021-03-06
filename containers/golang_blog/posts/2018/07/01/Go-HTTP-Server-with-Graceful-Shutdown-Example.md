Go's examples of an http server usually do not exit gracefully, so I created a hello world that will exit 0.

# GoLang HTTP Graceful Shutdown

From the examples I have come across online, they usually don't shutdown correctly (return exit code 0). I wanted to get some more experience and learn how to shutdown correctly. I did find examples online, but they also usually lack commenting for newbs.

Below I have given some examples of what I find in the golang readme's along with an example I have created with a simple Hello World that I hope is documented enough.

## Basic HTTP Server Examples

From Effective Go

```go
package main

import (
    "flag"
    "html/template"
    "log"
    "net/http"
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18

var templ = template.Must(template.New("qr").Parse(templateStr))

func main() {
    flag.Parse()
    http.Handle("/", http.HandlerFunc(QR))
    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

func QR(w http.ResponseWriter, req *http.Request) {
    templ.Execute(w, req.FormValue("s"))
}

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .}}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/" name=f method="GET"><input maxLength=1024 size=70
name=s value="" title="Text to QR Encode"><input type=submit
value="Show QR" name=qr>
</form>
</body>
</html>
`
```

From Golang - Writing Web Applications

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Graceful Shutdown Example

```go
package main

import (
        "context"
        "fmt"
        "log"
        "net/http"
        "os"
        "os/signal"
        "syscall"
)

var address = ""
var port = "6060"

// rootHandler is answering on / ("root")
func rootHandler(w http.ResponseWriter, r *http.Request) {

        // display to browser request
        fmt.Fprintln(w, "Hello World!")

        // log who requests root
        log.Println(r.RemoteAddr, "requested", r.RequestURI)
}

func main() {

        // A Server defines parameters for running an HTTP server.
        // The zero value for Server is a valid configuration.
        log.Printf("creating new http server \"%v:%v\"\n", address, port)
        srv := http.Server{
                Addr: address + ":" + port,
        }

        log.Println("add rootHandler handler to http server.")
        http.HandleFunc("/", rootHandler)

        log.Println("make channel")
        ch := make(chan os.Signal, 1)
        log.Println("use channel to wait for signals")
        signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

        log.Println("starting http server in goroutine")
        go func() {
                // ListenAndServe starts an HTTP server with a given address and handler.
                // The handler is usually nil, which means to use DefaultServeMux.
                // Handle and HandleFunc add handlers to DefaultServeMux:
                // ie. http.ListenAndServe(":6060", nil)
                err := srv.ListenAndServe()
                if err != nil {
                        log.Fatal("ListenAndServe Error:", err)
                }
        }()

        log.Println("waiting for signal... (channels block ~ Todd McLeod)")
        s := <-ch
        log.Println("Got SIGNAL:", s)

        log.Println("close channel")
        close(ch)

        log.Println("shutdown http server")
        err := srv.Shutdown(context.Background())
        if err != nil {
                // Error from closing listeners, or context timeout
                log.Fatal("Shutdown error:", err)
        }

        log.Println("Graceful Shutdown")
}
```

## Dockerfile

```none
FROM golang as binaryBuilder
WORKDIR /go/src/app
COPY . .
RUN env GOOS=linux GOARCH=386 go build -v -o main main.go

FROM scratch
WORKDIR /opt/helloworld/
COPY --from=binaryBuilder /go/src/app/main .
CMD ["./main"]
```

## Output

```none
[vagrant@centos7 ]$ docker ps -a
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS                    NAMES
bd7c1678bf06        helloworld:latest   "./main"            24 seconds ago      Up 23 seconds       0.0.0.0:6060->6060/tcp   helloworld
```

```none
[vagrant@centos7 ]$ curl -L localhost:6060
Hello World!
```

```none
[vagrant@centos7 ]$ time docker stop helloworld
helloworld

real    0m0.182s
user    0m0.053s
sys     0m0.008s
```

```none
[vagrant@centos7 ]$ docker ps -a
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS                     PORTS               NAMES
bd7c1678bf06        helloworld:latest   "./main"            53 seconds ago      Exited (0) 8 seconds ago                       helloworld
```

```none
[vagrant@centos7 ]$ docker logs helloworld
2018/07/01 16:54:22 creating new http server ":6060"
2018/07/01 16:54:22 add rootHandler handler to http server.
2018/07/01 16:54:22 make channel
2018/07/01 16:54:22 use channel to wait for signals
2018/07/01 16:54:22 starting http server in goroutine
2018/07/01 16:54:22 waiting for signal... (channels block ~ Todd McLeod)
2018/07/01 16:54:50 172.17.0.1:33622 requested /
2018/07/01 16:55:06 Got SIGNAL: terminated
2018/07/01 16:55:06 close channel
2018/07/01 16:55:06 shutdown http server
2018/07/01 16:55:06 Supposed Graceful Shutdown
```

## References

- [https://golang.org/pkg/net/http/](https://golang.org/pkg/net/http/)
- [https://golang.org/doc/effective_go.html#web_server](https://golang.org/doc/effective_go.html#web_server)
- [https://golang.org/doc/articles/wiki/](https://golang.org/doc/articles/wiki/)
