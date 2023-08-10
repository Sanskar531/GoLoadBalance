package structs

import (
	"time"
	r "math/rand"
)

type Request struct {
	data int
	resp chan float64
}

func createAndRequest(req chan Request) {
	resp := make(chan float64)

	for {
		time.Sleep(time.Microsecond)
		req <- Request {
			data: r.Intn(2000),
			resp: resp,
		}
		<-resp
	}
}
