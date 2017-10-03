package common

import "time"

const (
	// Builtins
	INT_TYPE        = "$int"
	FLOAT_TYPE      = "$float"
	STRING_TYPE     = "$str"
	DATE_TYPE       = "$date"
	DICT_TYPE       = "$dict"
	BOOL_TYPE       = "$bool"
	ENUM_TYPE       = "$enum"
	SERIAL_TYPE     = "$incr"
	UID_TYPE        = "$uid"
	UNIQUE_INT_TYPE = "$uniqint"

	// distributions
	DIST_TYPE    = "distribution"
	NORMAL_DIST  = "*normal"
	WEIGHT_DIST  = "*weight"
	PERCENT_DIST = "*percent"
)

var UNIX_EPOCH time.Time
var NOW time.Time

func init() {
	UNIX_EPOCH, _ = time.Parse("2006-01-02", "1970-01-01")
	NOW = time.Now()
}
