package logger

import (
	cfg "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/config"
)

var config cfg.Config = nil

func init() {
	cfg.OnConfig(func(newConfig cfg.Config) {
		config = newConfig
		Reset() // reset log root upon config change.
	})
}
