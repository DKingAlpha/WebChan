package utils

import (
	"fmt"
	"internal/shared_vars"
	"sync"
	"time"
)

func NewTimeoutQueue(channelId string, cap int, timeout int64, key string, perm ChanPerm) *TimeoutQueue {
	return &TimeoutQueue{
		lock:    sync.Mutex{},
		ChanId:  channelId,
		Cap:     cap,
		Timeout: timeout,
		Key:     key,
		Perm:    perm,
		Msgs:    nil,
	}
}

func (tq *TimeoutQueue) Empty() bool {
	tq.lock.Lock()
	defer tq.lock.Unlock()
	if tq.Msgs == nil || len(tq.Msgs) == 0 {
		return true
	} else {
		return false
	}
}

func (tq *TimeoutQueue) CleanTimeout() {
	tq.lock.Lock()
	defer tq.lock.Unlock()
	if tq.Timeout == -1 {
		return
	}
	if tq.Msgs == nil {
		return
	}
	earliestTime := shared_vars.CurrentTime - tq.Timeout
	validIndex := 0
	for {
		if validIndex >= len(tq.Msgs) {
			break
		}
		if tq.Msgs[validIndex].T < earliestTime {
			// earlier than allowed, Timeout
			// this one is invalid, check next
			validIndex += 1
		} else {
			// the earliest not Timeout, others wont be possible
			break
		}
	}
	if validIndex != 0 {
		tq.Msgs = tq.Msgs[validIndex:]
	}
}

func (tq *TimeoutQueue) AntiSpam() {
	tq.lock.Lock()
	defer tq.lock.Unlock()
	if tq.Msgs == nil {
		return
	}
	if tq.Cap < shared_vars.AntiSpamPeriodLimit {
		// Cap is less than rate limit
		dt := tq.Msgs[len(tq.Msgs) - 1].T - tq.Msgs[0].T
		if int(dt) * shared_vars.AntiSpamPeriodLimit <  shared_vars.AntiSpamPeriod * tq.Cap {
			// wait enough time for a new elem
			need2WaitT := time.Second * shared_vars.AntiSpamPeriod / shared_vars.AntiSpamPeriodLimit
			time.Sleep(need2WaitT)
		}
	} else {
		evaluateSince := len(tq.Msgs) - shared_vars.AntiSpamPeriodLimit
		if evaluateSince >= 0 {
			dt := tq.Msgs[len(tq.Msgs) - 1].T - tq.Msgs[evaluateSince].T
			need2WaitT := shared_vars.AntiSpamPeriod - dt
			if need2WaitT > 0 {
				time.Sleep(time.Duration(need2WaitT) * time.Second)
			}
		}
	}
}

func (tq *TimeoutQueue) Enqueue(chanMsg *ChanMessage) string {
	tq.AntiSpam()
	tq.CleanTimeout()
	tq.lock.Lock()
	defer tq.lock.Unlock()
	if tq.Msgs == nil {
		tq.Msgs = []*Message{}
	}
	if tq.Cap != -1 && len(tq.Msgs) > tq.Cap- 1{
		tq.Msgs = append(tq.Msgs[1:], &Message{
			T: chanMsg.T,
			M: chanMsg.M,
		})
	} else {
		tq.Msgs = append(tq.Msgs, &Message{
			T: chanMsg.T,
			M: chanMsg.M,
		})
	}
	if len(tq.Msgs) == tq.Cap {
		return "OK(Full)"
	} else {
		return "OK"
	}
}

func (tq *TimeoutQueue) GetData(showTime bool) string {
	tq.lock.Lock()
	defer tq.lock.Unlock()
	if tq.Msgs == nil {
		return ""
	}
	s := ""
	for _, i := range tq.Msgs {
		if showTime {
			s += fmt.Sprintf("%d:%s\n", i.T, i.M)
		} else {
			s += fmt.Sprintln(i.M)
		}
	}
	return s
}

func (tq *TimeoutQueue) GetDataFrom(from int64, showTime bool) string {
	tq.lock.Lock()
	defer tq.lock.Unlock()
	if tq.Msgs == nil {
		return ""
	}
	s := ""
	for _, i := range tq.Msgs {
		if i.T >= from {
			if showTime {
				s += fmt.Sprintf("%d:%s\n", i.T, i.M)
			} else {
				s += fmt.Sprintln(i.M)
			}
		}
	}
	return s
}

func (tq *TimeoutQueue) GetDataFromTo(from int64, to int64, showTime bool) string {
	tq.lock.Lock()
	defer tq.lock.Unlock()
	if tq.Msgs == nil {
		return ""
	}
	s := ""
	for _, i := range tq.Msgs {
		if i.T >= from && i.T <= to {
			if showTime {
				s += fmt.Sprintf("%d:%s\n", i.T, i.M)
			} else {
				s += fmt.Sprintln(i.M)
			}
		}
	}
	return s
}