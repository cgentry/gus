package mock

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/storage"
)

const (
	// DriverName Specifies the specific identity of this driver within a group
	DriverName      = "mock"
	IdentityStorage = "Mock"
	HelpShort       = "Simple JSON File. Store data JSON encoded into a file"
	HelpTemplate    = `

   This is a dummy driver used for testing purposes.

   `
)

type registerDriver struct{}

// Register is a simple wrapper to make sure registration occurs properly
func Register() {
	gdriver.Register(storage.DriverGroup, &registerDriver{})
}

// New() will return the resutls of the jsonfile New() function. You must cast
// this on return to the proper type (StorageDriver)
func (r *registerDriver) New() interface{} {
	return NewMockDriver()
}

// Identity provides a simple identifying string to the caller.
func (r *registerDriver) Identity(id int) string {
	switch id {
	case gdriver.IdentityShort:
		return HelpShort

	case gdriver.IdentityLong:
		return HelpTemplate
	}
	return IdentityStorage
}
