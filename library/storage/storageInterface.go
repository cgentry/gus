package storage

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/record/tenant"
)

// StorageDriver interface defines very general, high level operations for retrieval and storage of
// data. The back-storage can be a flat file, database or document store.
// The interfaces specify NO sql methods and flatten out operations
type StorageDriver interface {
	Open(connect string, extraDriverOptions string) (Conn, error)
}

// Storer interface gives the set of methods that a storage driver MAY implement.
// Because some are optional, see the methods for Store for ones that
// are required.
type Storer interface {
	// Required wrappers for required StorageDriver functions
	Open(connect string, extraDriverOptions string) error
	UserFetch(domain, lookupKey, lookkupValue string) (*tenant.User, error)
	UserInsert(user *tenant.User) error
	UserUpdate(user *tenant.User) error

	// Optional device connection functions
	CreateStore() error
	Close() error
	GetStorageConnector() Conn
	LastError() error
	IsOpen() bool
	Ping() error
	Release() error
	Reset()

	// The following are convienence wrapper functions for the UserXXX functions
	FetchUserByEmail(domain, email string) (*tenant.User, error)
	FetchUserByGUID(guid string) (*tenant.User, error)
	FetchUserByLogin(domain, loginName string) (*tenant.User, error)
	FetchUserByToken(token string) (*tenant.User, error)

	//  The following are wrappers for the gdriver routines.
	Id() string
	ShortHelp() string
	LongHelp() string

	// Internal routines
	SetDriverInterface(x gdriver.DriverInterface)
	GetDriverInterface() gdriver.DriverInterface

	SetStorageDriver(x StorageDriver)
	GetStorageDriver() StorageDriver
}

// Conn has a minimum call set that every driver is required to implement
type Conn interface {
	UserUpdate(user *tenant.User) error
	UserInsert(user *tenant.User) error

	UserFetch(domain, key, value string) (*tenant.User, error)
}

// Opener is an interface that will open up the storage driver for stroing data
type Opener interface {
	Open(connect string, extraDriverOptions string) (Conn, error)
}

// Creater is an optional Storge Creation interface
type Creater interface {
	CreateStore() error
}

// Closer is an optional interface. If this isn't implemented, no error is reported.
type Closer interface {
	Close() error
}

//Reseter is an optional interface. This will reset any errors and cleanup any intermediate results
type Reseter interface {
	Reset()
}

// Pinger is an optional database 'ping' interface. This will check the database connection
type Pinger interface {
	Ping() error
}

//Releaser optional interface. This will release any locks/resources that a driver may have set
//For example, the MySQL will do a SELECT...FOR UPDATE for all of the FetchXXX calls. The
//release will cause an explicit commit. This, in the code, will be called by a 'defer' call after
//any fetch/insert operation. For other drivers, it can be ignored or perform any other operation
//required.
// Note that SQLITE doesn't do anything at this stage as it isn't really considered a robust, fully
// hardened storage mechanism. Document-style interfaces will probably not use it either.
type Releaser interface {
	Release() error
}
