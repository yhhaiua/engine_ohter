package perf

import (
	"net/http"
	_ "net/http/pprof"
)

func init()  {
	go startPerf()
}

func startPerf()  {
	_ = http.ListenAndServe("localhost:6060", nil)
}
