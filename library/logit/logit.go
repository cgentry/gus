// Package logit handles the setup and logging for the program. Like other driver-based packages,
// you select the logger you want and pass it configuration information. Then you simply log
// messages and they will be formatted by the driver

package logit

import (
//"strings"
)

const (
	DRIVER_GROUP = "logging"
)
// The interface gives the set of methods that an encryption driver must implement.
type LogitDriver interface {
	Open() LogitDriver
	Write(level int, logval ...interface{})
	Close()

	Id() string
	ShortHelp() string
	LongHelp() string
}


func Select( name string ) LogitDriver {
	gdriver.Default( DRIVER_GROUP , name )
	return GetDriver()
}

// This will panic if no drivers have been registered
func GetDriver() LogitDriver {
	return gdriver.MustNewDefault(DRIVER_GROUP).(LogitDriver)
}
