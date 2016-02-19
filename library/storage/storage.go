// The storage drivers are used to store identification for a user and nothing more.
// Standard encryption drivers are mock, sqlite, jsonfile and mysql.
//
//
// Drivers are selected by:
//    db := storage.Select( driverName )
// To conform to the gdriver interface, all drivers must have New() and Identity()
// functions.

package storage

import (
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gdriver"
)

const DRIVER_GROUP = "storage"



type storeMap map[string]Storer

// Store holds the state for any storage driver. It allows you to have
// consistent returns, such as getting the last error, discovering how
// a connection was made (connectString) or the name of the driver (name)
type Store struct {
	name		  string
	connectString string
	isOpen        bool
	lastError     error
	rawDriver	  gdriver.DriverInterface
	driver        StorageDriver			// From New(). This is not the db connection
	connection    Conn					// from Open(). this is the true DB connection
}

func NewStore( name string ) Storer {
	return &Store{
		name:          name,
		isOpen:        false,
		connectString: "",
		lastError:     ErrNoDriverFound,
	}
}

// Select will pick a registered driver for use in the system.
// Only one driver can be selected at a time. Calling GetDriver will
// return the current driver
// This will panic if no drivers have been registered
func SetDefault(name string) Storer {
	gdriver.Default(DRIVER_GROUP,name)
	return GetDriver(name)
}

// GetDriver will fetch the gdriver, call New and build up the storage driver
// interface. Store defines all possible calls and acts as a simple gateway in order to
// make any driver respond to any calls, even if it isn't implemented.
func GetDefaultDriver() Storer {
	name,_ := gdriver.GetDefaultName( DRIVER_GROUP )
	return GetDriver( name )
}

func GetDriver( name string ) Storer {
	st := NewStore( name )
	st.SetStorageDriver( gdriver.MustNew( DRIVER_GROUP , name ).(StorageDriver ) )
	drive, _ := gdriver.GetDriver( DRIVER_GROUP, name )
	st.SetDriverInterface(drive )
	return st
}

func Open( defaultDriverName, dsn, options string ) ( Storer , error ) {
	storeObject := SetDefault( defaultDriverName )
	return storeObject, storeObject.Open( dsn, options )
}

func ( s *Store ) SetDriverInterface( x gdriver.DriverInterface ){
	s.rawDriver = x
}
func ( s *Store ) GetDriverInterface() gdriver.DriverInterface {
	return s.rawDriver
}
func ( s * Store ) SetStorageDriver( x StorageDriver ) {
	s.driver = x
}
func ( s *Store ) GetStorageDriver( ) StorageDriver {
	return s.driver
}