package utils

import (
	"sync"
)

type ChanMessage struct {
	Ch string
	T int64
	M string
}

type Message struct {
	T int64
	M string
}

type ChanPerm struct {
	// using int instead of bool to save space in JSON
	R int
	W int
	D int
}

type TimeoutQueue struct {
	lock    sync.Mutex
	ChanId  string
	Cap     int
	Timeout int64
	Key     string
	Perm    ChanPerm
	Msgs    []*Message
}
