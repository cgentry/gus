package mock

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/storage"
)

const (
	DRIVER_NAME      = "mock"
	STORAGE_IDENTITY = "Mock"
	SHORT_HELP       = "Simple JSON File. Store data JSON encoded into a file"
	HELP_TEMPLATE    = `

   This is a dummy driver used for testing purposes.

   `
)

type registerDriver struct{}

// Register is a simple wrapper to make sure registration occurs properly
func Register() {
	gdriver.Register(storage.DRIVER_GROUP, &registerDriver{})
}

// New() will return the resutls of the jsonfile New() function. You must cast
// this on return to the proper type (StorageDriver)
func (r *registerDriver) New() interface{} {
	return NewMockDriver()
}

// Identity provides a simple identifying string to the caller.
func (r *registerDriver) Identity(id int) string {
	switch id {
	case gdriver.IDENT_SHORT:
		return SHORT_HELP

	case gdriver.IDENT_LONG:
		return HELP_TEMPLATE
	}
	return STORAGE_IDENTITY
}
