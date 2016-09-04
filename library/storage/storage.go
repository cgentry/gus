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
	"github.com/cgentry/gdriver"
	. "github.com/cgentry/gus/ecode"
)

// DriverGroup defines a logical grouping for the drivers
const DriverGroup = "storage"

type storeMap map[string]Storer

// Store holds the state for any storage driver. It allows you to have
// consistent returns, such as getting the last error, discovering how
// a connection was made (connectString) or the name of the driver (name)
type Store struct {
	name          string
	connectString string
	isOpen        bool
	lastError     error
	rawDriver     gdriver.DriverInterface
	driver        StorageDriver // From New(). This is not the db connection
	connection    Conn          // from Open(). this is the true DB connection
}

// NewStore will return an address of the Store structure.
func NewStore(name string) Storer {
	return &Store{
		name:          name,
		isOpen:        false,
		connectString: "",
		lastError:     ErrNoDriverFound,
	}
}

// SetDefault sets the name of the driver that should be the default driver.
func SetDefault(name string) Storer {
	gdriver.Default(DriverGroup, name)
	return GetDriver(name)
}

// GetDefaultDriver fetches the name of the default driver and then activates the driver
func GetDefaultDriver() Storer {
	name, _ := gdriver.GetDefaultName(DriverGroup)
	return GetDriver(name)
}

// GetDriver will pick the driver from the map and return a fully initialised storage driver to the caller.
func GetDriver(name string) Storer {
	st := NewStore(name)
	st.SetStorageDriver(gdriver.MustNew(DriverGroup, name).(StorageDriver))
	drive, _ := gdriver.GetDriver(DriverGroup, name)
	st.SetDriverInterface(drive)
	return st
}

// Open will set the default name and then open the driver for storage activity
func Open(defaultDriverName, dsn, options string) (Storer, error) {
	storeObject := SetDefault(defaultDriverName)
	return storeObject, storeObject.Open(dsn, options)
}

//SetDriverInterface is a low level function and should never be called. It will force the
// storage driver to use the passed driver instead of the one defined.
func (s *Store) SetDriverInterface(x gdriver.DriverInterface) {
	s.rawDriver = x
}

// GetDriverInterface is a low level function that replaces the raw storage driver with the driver
// passed.
func (s *Store) GetDriverInterface() gdriver.DriverInterface {
	return s.rawDriver
}

// SetStorageDriver will set the higher-level driver to the storagedriver passed.
func (s *Store) SetStorageDriver(x StorageDriver) {
	s.driver = x
}

// GetStorageDriver will return the objects storage driver. This is for internal use only
func (s *Store) GetStorageDriver() StorageDriver {
	return s.driver
}
