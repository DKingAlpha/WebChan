package shared_vars

import "time"

var CurrentTime = time.Now().Unix()

const (
	AntiSpamPeriod      = 30			// in 30s
	AntiSpamPeriodLimit = 300			// 300 max
	ActivityTimeout     = 1*24*60*60	// 1 day
	DumpDBPath          = "db.json"
	DumpActivityPath    = "activity.json"
	AdminPassword       = "DKingAlpha"
)
