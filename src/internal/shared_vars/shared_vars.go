package shared_vars

import "time"

var CurrentTime = time.Now().Unix()

const (
	AntiSpamPeriod      = 30
	AntiSpamPeriodLimit = 100
	ActivityTimeout     = 1*24*60*60
	DumpDBPath          = "db.json"
	DumpActivityPath    = "activity.json"
	AdminPassword       = "DKingAlpha"
)
