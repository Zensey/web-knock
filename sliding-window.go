package main

import "time"

const windowTTL = time.Second * 60

type window struct {
	start time.Time
	count int
}

func putEvent(m map[string]window, ip string, t time.Time) int {
	w, ok := m[ip]
	if !ok {
		w = window{start: t}
	}

	if t.Sub(w.start) > windowTTL {
		w.count = 0
	}
	w.count++
	m[ip] = w

	return w.count
}
