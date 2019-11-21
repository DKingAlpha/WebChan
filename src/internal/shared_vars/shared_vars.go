package shared_vars

import "time"

var CurrentTime = time.Now().Unix()

const (
	AntiSpamPeriod = 30
	AntiSpamPeriodLimit = 100
	DumpJSONPath = "db.json"
	AdminPassword = "DKingAlpha"
)
