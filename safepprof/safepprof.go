package safepprof

import (
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"

	"golang.org/x/net/netutil"
)

// Run serves pprof endpoint on addr :6060. It setups a dedicated http.Server
// with sane defaults and enable additional profiling.
func Run() error {
	// Enable additional profiles. Values taken from CockroachDB.
	// See: https://github.com/cockroachdb/cockroach/pull/30233/
	runtime.SetMutexProfileFraction(1_000)  // 1 sample per 1000 mutex contention events
	runtime.SetBlockProfileRate(10_000_000) // 1 sample per 10 milliseconds spent blocking

	// Register handlers, same as pprof.init()
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := &http.Server{
		Handler: mux,

		// Sane defaults, to make sure that you can't slow down the service for a
		// long period of time.
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	// Disable Keep-Alives
	srv.SetKeepAlivesEnabled(false)

	ln, err := net.Listen("tcp", ":6060")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	// Accept only 1 connection at a given time
	ln = netutil.LimitListener(ln, 1)

	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
