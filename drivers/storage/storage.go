// The encryption drivers are used to encrypt and decrypt passwords stored for a user.
// Standard encryption drivers are sha512, bcrypt and plaintext (not suitable for
// production systems)
//
// All drivers need to call Register in order to be usable by the system. Failure to select
// an encryption driver will cause a panic during runtime.
//
// Drivers are selected by:
//    crypt := encryption.Select( driverName ).Setup( options )
// The call to Setup() is optional and driver specific.

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

// Pick a registered driver for use in the system. Only one driver can be selected at a time.
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
