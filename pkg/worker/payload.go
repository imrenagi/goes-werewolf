package worker

import "time"

type Payload struct {
	NameName  string
	Delay time.Duration
	Fn func() error
}

type Job struct {
	Payload Payload
}
