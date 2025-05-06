package model

import (
	"time"
)

const windowTTL = time.Second * 60 * 3

type Window struct {
	start time.Time
	count int
}

func PutEvent(m map[string]Window, ip string, t time.Time) int {
	w, ok := m[ip]
	if !ok {
		w = Window{start: t}
	}

	if t.Sub(w.start) > windowTTL {
		w.count = 0
		w.start = t
	}
	w.count++
	m[ip] = w

	return w.count
}
