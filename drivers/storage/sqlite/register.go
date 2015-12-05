package sqlite

import (
	"github.com/cgentry/gdriver"
)

const (
	DRIVER_NAME 		= "sqlite"
	STORAGE_IDENTITY	= "Sqlite"
	SHORT_HELP			= "Simple SQLite3 driver. Only suitable for testing purposes."
	HELP_TEMPLATE		= `

   This is a lightweight driver meant for testing and debugging systems.
   It provides a full database testing system and stores the data in a
   standard Sqlite3 file, accessable by the command line tool.

   DSN: This is a simple string that defines the path where to store the
        database file. The directory must be writable.

   Options: This is just a string that defines the table to store the data in.
         If nothing is passed, the default is "User".

   `
}

type registerDriver struct{}

// Register is a simple wrapper to make sure registration occurs properly
func Register() {
	gdriver.Register(encryption.DRIVER_GROUP, &registerDriver{})
}

func (r *registerDriver) New() interface{} {
	return New()
}

func (r *registerDriver) Identity(id int) string {
	switch id {
	case gdriver.IDENT_NAME:
		return STORAGE_IDENTITY
	case gdriver.IDENT_SHORT:
		return 	SHORT_HELP
	case gdriver.IDENT_LONG:
		return HELP_TEMPLATE
	}
	return gdriver.IDENT_UNKNOWN
}
