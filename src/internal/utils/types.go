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
type TimeoutQueue struct {
	lock    sync.Mutex
	ChanId  string
	Cap     int
	Timeout int64
	Msgs    []*Message
}
