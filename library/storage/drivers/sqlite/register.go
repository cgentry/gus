package sqlite

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/storage"
)

const (
	// DriverName Specifies the specific identity of this driver within a group
	DriverName      = "sqlite"
	IdentityStorage = "Sqlite"
	HelpShort       = "Simple SQLite3 driver. Only suitable for testing purposes."
	HelpTemplate    = `

   This is a lightweight driver meant for testing and debugging systems.
   It provides a full database testing system and stores the data in a
   standard Sqlite3 file, accessable by the command line tool.

   DSN: This is a simple string that defines the path where to store the
        database file. The directory must be writable.

   Options: This is just a string that defines the table to store the data in.
         If nothing is passed, the default is "User".

   `
)

type registerDriver struct{}

// Register is a simple wrapper to make sure registration occurs properly
func Register() {
	gdriver.Register(storage.DriverGroup, &registerDriver{})
}

func SetDefault() {
	gdriver.Default(storage.DriverGroup,DriverName)
}

// New() will return the resutls of the Sqlite New() function. You must cast
// this on return to the proper type (StorageDriver)
func (r *registerDriver) New() interface{} {
	return NewSqliteDriver()
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
