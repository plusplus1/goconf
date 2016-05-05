package goconf

import (
	"strconv"
)

func (cfg Configuration) Sub(key string) Configuration {
	temp := cfg.Get(key)
	if nil != temp {
		if ret, isOk := temp.(Configuration); isOk {
			return ret
		}
	}
	return nil
}

func (cfg Configuration) String(key string) string {
	temp := cfg.Get(key)
	if nil != temp {
		if ret, isOk := temp.(string); isOk {
			return ret
		}
	}
	return ""
}

func (cfg Configuration) Integer(key string) int {
	if ret, e := strconv.Atoi(cfg.String(key)); nil == e {
		return ret
	}
	return 0
}

// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
// Any other value returns false.
func (cfg Configuration) Bool(key string) bool {
	if ret, e := strconv.ParseBool(cfg.String(key)); nil == e {
		return ret
	}
	return false
}

func (cfg Configuration) Float64(key string) float64 {
	if ret, e := strconv.ParseFloat(cfg.String(key), 64); nil == e {
		return ret
	}
	return 0.0
}
