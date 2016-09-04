package jsonfile

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/storage"
)

const (
	// DriverName Specifies the specific identity of this driver within a group
	DriverName      = "jsonfile"
	IdentityStorage = "JsonFile"
	HelpShort       = "Simple JSON File. Store data JSON encoded into a file"
	HelpTemplate    = `

   This is a testing-only system that will read/write to json files. It is very
   slow and only meant for testing purposes.

   DSN: This is a simple string that defines the path where to store the
        JSON file. The directory must be writable.

   Options: This is unused..

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
	return New()
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
