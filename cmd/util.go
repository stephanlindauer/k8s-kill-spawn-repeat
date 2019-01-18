package cmd

import (
	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	"github.com/sirupsen/logrus"
)

func must(err error) {
	if err != nil {
		logrus.Debugf("%+v", err)
		logrus.Error(err)
		logrus.Warn("Something failed.")
		cmdutil.Exit(1)
	}
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
