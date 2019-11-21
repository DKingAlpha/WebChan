package utils

import (
	"strings"
)

func GetPerm(perm string) *ChanPerm {
	cp := ChanPerm{
		R: 0,
		W: 0,
		D: 0,
	}
	for _, p := range perm {
		if p == 'r' || p == 'R' {
			cp.R = 1
		}
		if p == 'w' || p == 'W' {
			cp.W = 1
		}
		if p == 'd' || p == 'D' {
			cp.D = 1
		}
	}
	return &cp
}

func GetUrlArgs(args string) map[string]string {
	q := map[string]string{}
	for _, kv :=  range strings.Split(args, "&") {
		kva := strings.SplitN(kv, "=", 2)
		if len(kva) != 2 {
			continue
		}
		q[kva[0]] = kva[1]
	}
	return q
}
