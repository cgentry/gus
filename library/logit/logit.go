// Package logit handles the setup and logging for the program. Like other driver-based packages,
// you select the logger you want and pass it configuration information. Then you simply log
// messages and they will be formatted by the driver

package logit

import (
	"github.com/cgentry/gdriver"
)

const (
	DriverGroup = "logging"
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


func SetDefault( name string ) LogitDriver {
	gdriver.Default( DriverGroup , name )
	return GetDriver( name )
}

// This will panic if no drivers have been registered
func GetDefaultDriver() LogitDriver {
	return gdriver.MustNewDefault(DriverGroup).(LogitDriver)
}

// This will panic if no drivers have been registered
func GetDriver(name string ) LogitDriver {
	return gdriver.MustNew(DriverGroup, name ).(LogitDriver)
}
