package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/sensiblecodeio/tiny-ssl-reverse-proxy/proxyprotocol"
)


/*
*	tinyssl
*/
// Version number
const Version = "0.22.0"

var message = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>
Backend Unavailable
</title>
<style>
body {
	font-family: fantasy;
	text-align: center;
	padding-top: 20%;
	background-color: #f1f6f8;
}
</style>
</head>
<body>
<h1>503 Backend Unavailable</h1>
<p>Sorry, we&lsquo;re having a brief problem. You can retry.</p>
<p>If the problem persists, please get in touch.</p>
</body>
</html>`

type ConnectionErrorHandler struct{ http.RoundTripper }

func (c *ConnectionErrorHandler) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.RoundTripper.RoundTrip(req)
	if err != nil {
		log.Printf("Error: backend request failed for %v: %v",
			req.RemoteAddr, err)
	}
	if _, ok := err.(*net.OpError); ok {
		r := &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Body:       ioutil.NopCloser(bytes.NewBufferString(message)),
		}
		return r, nil
	}
	return resp, err
}

func main() {
	var (
		listen, cert, key, where           string
		useTLS, useLogging, behindTCPProxy bool
		flushInterval                      time.Duration
	)
	flag.StringVar(&listen, "listen", ":443", "Bind address to listen on")
	flag.StringVar(&key, "key", "/etc/ssl/private/key.pem", "Path to PEM key")
	flag.StringVar(&cert, "cert", "/etc/ssl/private/cert.pem", "Path to PEM certificate")
	flag.StringVar(&where, "where", "http://localhost:80", "Place to forward connections to")
	flag.BoolVar(&useTLS, "tls", true, "accept HTTPS connections")
	flag.BoolVar(&useLogging, "logging", true, "log requests")
	flag.BoolVar(&behindTCPProxy, "behind-tcp-proxy", false, "running behind TCP proxy (such as ELB or HAProxy)")
	flag.DurationVar(&flushInterval, "flush-interval", 0, "minimum duration between flushes to the client (default: off)")
	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%v version %v\n\n", os.Args[0], Version)
		oldUsage()
	}
	flag.Parse()

	url, err := url.Parse(where)
	if err != nil {
		log.Fatalln("Fatal parsing -where:", err)
	}

	httpProxy := httputil.NewSingleHostReverseProxy(url)
	httpProxy.Transport = &ConnectionErrorHandler{http.DefaultTransport}
	httpProxy.FlushInterval = flushInterval

	var handler http.Handler

	handler = httpProxy

	originalHandler := handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/_version" {
			w.Header().Add("X-Tiny-SSL-Version", Version)
		}
		r.Header.Set("X-Forwarded-Proto", "https")
		originalHandler.ServeHTTP(w, r)
	})

	if useLogging {
		handler = &LoggingMiddleware{handler}
	}

	server := &http.Server{Addr: listen, Handler: handler}

	switch {
	case useTLS && behindTCPProxy:
		err = proxyprotocol.BehindTCPProxyListenAndServeTLS(server, cert, key)
	case behindTCPProxy:
		err = proxyprotocol.BehindTCPProxyListenAndServe(server)
	case useTLS:
		err = server.ListenAndServeTLS(cert, key)
	default:
		err = server.ListenAndServe()
	}

	log.Fatalln(err)
}

















type LoggingMiddleware struct {
	http.Handler
}

type ResponseRecorder struct {
	ResponseWriter http.ResponseWriter
	response       int
	*WriteCounter
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{w, 0, &WriteCounter{w, 0}}
}

func (r *ResponseRecorder) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *ResponseRecorder) WriteHeader(n int) {
	r.ResponseWriter.WriteHeader(n)
	r.response = n
}

func (r *ResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("Not a Hijacker: %T", r.ResponseWriter)
	}
	return hijacker.Hijack()
}

func (r *ResponseRecorder) Flush() {
	flusher, ok := r.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}
	flusher.Flush()
}

type WriteCounter struct {
	io.Writer
	nBytes int
}

func (r *WriteCounter) Write(bs []byte) (n int, err error) {
	if r.Writer != nil {
		n, err = r.Writer.Write(bs)
	} else {
		n = len(bs)
	}
	r.nBytes += n
	return n, err
}

func (x *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recorder := NewResponseRecorder(w)

	uploaded := &WriteCounter{Writer: ioutil.Discard}
	r.Body = struct {
		io.Reader
		io.Closer
	}{io.TeeReader(r.Body, uploaded), r.Body}

	start := time.Now()
	x.Handler.ServeHTTP(recorder, r)
	duration := time.Since(start)

	log.Printf("%21v %3d %10d %10d %7.1fms %4v %v%-30v %v",
		r.RemoteAddr,
		recorder.response,
		uploaded.nBytes,
		recorder.nBytes,
		duration.Seconds()*1000,
		r.Method,
		r.URL.Host,
		r.URL.EscapedPath(),
		r.Header.Get("User-Agent"))
}