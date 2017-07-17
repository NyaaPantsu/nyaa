// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pprof serves via its HTTP server runtime profiling data
// in the format expected by the pprof visualization tool.
//
// The package is typically only imported for the side effect of
// registering its HTTP handlers.
// The handled paths all begin with /debug/pprof/.
//
// To use pprof, link this package into your program:
//	import _ "net/http/pprof"
//
// If your application is not already running an http server, you
// need to start one. Add "net/http" and "log" to your imports and
// the following code to your main function:
//
// 	go func() {
// 		log.Println(http.ListenAndServe("localhost:6060", nil))
// 	}()
//
// Then use the pprof tool to look at the heap profile:
//
//	go tool pprof http://localhost:6060/debug/pprof/heap
//
// Or to look at a 30-second CPU profile:
//
//	go tool pprof http://localhost:6060/debug/pprof/profile
//
// Or to look at the goroutine blocking profile, after calling
// runtime.SetBlockProfileRate in your program:
//
//	go tool pprof http://localhost:6060/debug/pprof/block
//
// Or to collect a 5-second execution trace:
//
//	wget http://localhost:6060/debug/pprof/trace?seconds=5
//
// To view all available profiles, open http://localhost:6060/debug/pprof/
// in your browser.
//
// For a study of the facility in action, visit
//
//	https://blog.golang.org/2011/06/profiling-go-programs.html
//
package pprofController

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Cmdline responds with the running program's
// command line, with arguments separated by NUL bytes.
// The package initialization registers it as /debug/pprof/cmdline.
func PprofCmdline(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(c.Writer, strings.Join(os.Args, "\x00"))
}

func sleep(w http.ResponseWriter, d time.Duration) {
	var clientGone <-chan bool
	if cn, ok := w.(http.CloseNotifier); ok {
		clientGone = cn.CloseNotify()
	}
	select {
	case <-time.After(d):
	case <-clientGone:
	}
}

// Profile responds with the pprof-formatted cpu profile.
// The package initialization registers it as /debug/pprof/profile.
func PprofProfile(c *gin.Context) {
	sec, _ := strconv.ParseInt(c.Request.FormValue("seconds"), 10, 64)
	if sec == 0 {
		sec = 30
	}

	// Set Content Type assuming StartCPUProfile will work,
	// because if it does it starts writing.
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	if err := pprof.StartCPUProfile(c.Writer); err != nil {
		// StartCPUProfile failed, so no writes yet.
		// Can change header back to text content
		// and send error code.
		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Could not enable CPU profiling: %s\n", err)
		return
	}
	sleep(c.Writer, time.Duration(sec)*time.Second)
	pprof.StopCPUProfile()
}

// Trace responds with the execution trace in binary form.
// Tracing lasts for duration specified in seconds GET parameter, or for 1 second if not specified.
// The package initialization registers it as /debug/pprof/trace.
func PprofTrace(c *gin.Context) {
	sec, err := strconv.ParseFloat(c.Request.FormValue("seconds"), 64)
	if sec <= 0 || err != nil {
		sec = 1
	}

	// Set Content Type assuming trace.Start will work,
	// because if it does it starts writing.
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	if err := trace.Start(c.Writer); err != nil {
		// trace.Start failed, so no writes yet.
		// Can change header back to text content and send error code.
		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Could not enable tracing: %s\n", err)
		return
	}
	sleep(c.Writer, time.Duration(sec*float64(time.Second)))
	trace.Stop()
}

// Symbol looks up the program counters listed in the request,
// responding with a table mapping program counters to function names.
// The package initialization registers it as /debug/pprof/symbol.
func PprofSymbol(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// We have to read the whole POST body before
	// writing any output. Buffer the output here.
	var buf bytes.Buffer

	// We don't know how many symbols we have, but we
	// do have symbol information. Pprof only cares whether
	// this number is 0 (no symbols available) or > 0.
	fmt.Fprintf(&buf, "num_symbols: 1\n")

	var b *bufio.Reader
	if c.Request.Method == "POST" {
		b = bufio.NewReader(c.Request.Body)
	} else {
		b = bufio.NewReader(strings.NewReader(c.Request.URL.RawQuery))
	}

	for {
		word, err := b.ReadSlice('+')
		if err == nil {
			word = word[0 : len(word)-1] // trim +
		}
		pc, _ := strconv.ParseUint(string(word), 0, 64)
		if pc != 0 {
			f := runtime.FuncForPC(uintptr(pc))
			if f != nil {
				fmt.Fprintf(&buf, "%#x %s\n", pc, f.Name())
			}
		}

		// Wait until here to check for err; the last
		// symbol will have an err because it doesn't end in +.
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(&buf, "reading request: %v\n", err)
			}
			break
		}
	}

	c.Writer.Write(buf.Bytes())
}

// Handler returns an HTTP handler that serves the named profile.
func Handler(name string) http.Handler {
	return handler(name)
}

type handler string

func (name handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	debug, _ := strconv.Atoi(r.FormValue("debug"))
	p := pprof.Lookup(string(name))
	if p == nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Unknown profile: %s\n", name)
		return
	}
	gc, _ := strconv.Atoi(r.FormValue("gc"))
	if name == "heap" && gc > 0 {
		runtime.GC()
	}
	p.WriteTo(w, debug)
	return
}

// Index responds with the pprof-formatted profile named by the request.
// For example, "/debug/pprof/heap" serves the "heap" profile.
// Index responds to a request for "/debug/pprof/" with an HTML page
// listing the available profiles.
func PprofIndex(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/debug/pprof/") {
		name := strings.TrimPrefix(c.Request.URL.Path, "/debug/pprof/")
		if name != "" {
			handler(name).ServeHTTP(c.Writer, c.Request)
			return
		}
	}

	profiles := pprof.Profiles()
	if err := indexTmpl.Execute(c.Writer, profiles); err != nil {
		log.Print(err)
	}
}

var indexTmpl = template.Must(template.New("index").Parse(`<html>
  <head>
  <title>/debug/pprof/</title>
  </head>
  <body>
  /debug/pprof/<br>
  <br>
  profiles:<br>
  <table>
  {{range .}}
  <tr><td align=right>{{.Count}}<td><a href="{{.Name}}?debug=1">{{.Name}}</a>
  {{end}}
  </table>
  <br>
  <a href="goroutine?debug=2">full goroutine stack dump</a><br>
  </body>
  </html>
  `))
