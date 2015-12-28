// The storage drivers are used to store identification for a user and nothing more.
// Standard encryption drivers are mock, sqlite, jsonfile and mysql.
//
//
// Drivers are selected by:
//    db := encryption.Select( driverName )
// To conform to the gdriver interface, all drivers must have New() and Identity()
// functions. To conform to the minimum storage driver, they must also have

package storage

import (
	"encoding/json"
	"github.com/cgentry/gdriver"
	"strings"
)

const DRIVER_GROUP = "storage"

// The interface gives the set of methods that a storage driver MAY implement.
// Because some are optional, see the methods for Store for ones that
// are required

type Storer interface {
	Open(connect string, extraDriverOptions string) (Storer, error)
	Close() error
	GetStorageConnector() Conn
	LastError() error
	IsOpen() bool
	Ping() error
	Release() error
	Reset()

	FetchUserByEmail(domain, email string) (*tenant.User, error)
	FetchUserByGuid(guid string) (*tenant.User, error)
	FetchUserByLogin(domain, loginName string) (*tenant.User, error)
	FetchUserByToken(token string) (*tenant.User, error)
	UserFetch(domain, lookupKey, lookkupValue string) (*tenant.User, error)
	UserInsert(user *tenant.User) error
	UserUpdate(user *tenant.User) error

	//  The following are wrappers for the gdriver routines.
	Id() string
	ShortHelp() string
	LongHelp() string
}

type storeMap map[string]Storer

// Store holds the state for any storage driver. It allows you to have
// consistent returns, such as getting the last error, discovering how
// a connection was made (connectString) or the name of the driver (name)
type Store struct {
	name		  string
	connectString string
	isOpen        bool
	lastError     error
	driver        StorageDriver
	connection    Conn
}

func NewStore( name string ) *Store {
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
func Select(name string) *Store {
	gdriver.Default(DRIVER_GROUP,name)
	selectedStore := NewStore( name )
	selectedStore.driver = gdriver.MustNew(DRIVER_GROUP, name).
	return selectedStore
}

func GetDriver() *Store {
	return gdriver.MustNewDefault( DRIVER_GROUP).( *Store )
}
