package storage

import (
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/tenant"
	"github.com/cgentry/gdriver"
)

// These are the names of fields we expect to occur in the database and will
// pass to database functions when performing UserFetch operations. You may
// map them in the driver-level routines in order to provide names that are
// more appropriate to the driver mechanism.
const (
	FIELD_EMAIL = `Email`
	FIELD_NAME  = `FullName`
	FIELD_GUID  = `Guid`
	FIELD_LOGIN = `LoginName`
	FIELD_TOKEN = `Token`
)
const driver_name = "Storage"

const MATCH_ANY_DOMAIN = "*"


/* ======================================================
 * 					Store functions
 * ======================================================
 */


func ( s *Store ) Id() string {
	return s.rawDriver.Identity( gdriver.IDENT_NAME )
}

func ( s *Store ) ShortHelp() string {
	return s.rawDriver.Identity( gdriver.IDENT_SHORT )
}

func ( s *Store ) LongHelp() string {
	return s.rawDriver.Identity( gdriver.IDENT_LONG )
}


// Return the actual connection to the database for low-level access.
// This should be avoided unless you are coding for a very non-portable
// function
func (s *Store) GetStorageConnector() Conn {
	return s.connection
}

/*
 * The following functions are provided by this class and are not
 * encapsulated
 */
// Return the last known error condition that was given by a call
func (s *Store) LastError() error {
	return s.lastError
}

func (s *Store) SetLastError(err error) *Store {
	s.lastError = err
	return s
}

// IsOpen will return the the open status of the connection
func (s *Store) IsOpen() bool {
	return s.isOpen
}

func ( s *Store ) saveLastError( err error ) error {
	s.lastError = err
	return err
}

func ( s *Store ) ClearErrors(){
	s.lastError = nil
}


/* ======================================================
 * 					Mandatory functions
 *		StorageDriver
 * ======================================================
 */

// Open a connection to the storage mechanism and return both a storage
// structure and an error status of the open
func (s *Store) Open(connect string, extraDriverOptions string) error {

	if s.isOpen == true {
		return s.saveLastError(ErrAlreadyOpen)
	}
	s.isOpen = false
	s.lastError = nil
	s.connectString = connect
	s.connection, s.lastError = s.driver.Open( connect, extraDriverOptions )

	return s.lastError
}

func (s *Store) UserUpdate(user *tenant.User) error {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return ErrNotOpen
	}
	return s.saveLastError(s.connection.UserUpdate(user))
}

func (s *Store) UserInsert(user *tenant.User) error {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return ErrNotOpen
	}
	return s.saveLastError(s.connection.UserInsert(user))
}

// Fetch a user's record using the domain, a field name and the field value. There will only be one record
// returned. If you pass MATCH_ANY_DOMAIN as the domain, this will only be valid for a small number of
// key-types (e.g. enforced unique keys.)
func (s *Store) UserFetch(domain, lookupKey, lookkupValue string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	if domain == MATCH_ANY_DOMAIN {
		if lookupKey != FIELD_GUID || lookupKey != FIELD_TOKEN {
			return nil, ErrMatchAnyNotSupported
		}
	}
	rec, err := s.connection.UserFetch(domain, lookupKey, lookkupValue)
	s.lastError = err
	return rec, err
}

/* ======================================================
 * 					Optional functions
 * If not provided, they should return a 'good' result
 * rather than an error
 * ======================================================
 */

// Reset any errors or intermediate conditions
func (s *Store) Reset() {
	s.ClearErrors()
	if reseter, found := s.connection.(Reseter); found {
		reseter.Reset()
	}
	return
}

// Release any locks or memory
func (s *Store) Release() error {
	s.ClearErrors()
	if release, found := s.connection.(Releaser); found {
		s.SetLastError(release.Release())
	}
	return s.LastError()
}

// Close the connection to the storage mechanism. If there is no close routine
// ignore the call
func (s *Store) Close() error {
	if s.isOpen != true {
		return s.saveLastError(ErrNotOpen)
	}
	s.isOpen = false
	s.ClearErrors()
	if closer, found := s.connection.(Closer); found {
		s.lastError = closer.Close()
	}
	return s.lastError
}

// If implemented, create the basic storage. If not implemented, an error will be returned.
func (s *Store) CreateStore() error {
	if s.isOpen != true {
		return s.saveLastError(ErrNotOpen)
	}

	if creater, found := s.connection.(Creater); found {
		return s.saveLastError(creater.CreateStore())
	}
	return ErrNoSupport
}

func (s *Store) Ping() error {
	s.ClearErrors()
	if pinger, found := s.connection.(Pinger); found {
		s.lastError = pinger.Ping()
	}
	return s.lastError
}


/* ------------------------ THE FOLLOWING ARE 'CONVENIENCE' FUNCTIONS ***********************/

// Fetch a user by the GUID. No domains are required as this is the primary (or unique) key
func (s *Store) FetchUserByGuid(guid string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(MATCH_ANY_DOMAIN, FIELD_GUID, guid)
	s.lastError = err
	return rec, err
}

// Fetch a user by the logged-in token. If the user is not logged in, a 'User not found' error is returned.
func (s *Store) FetchUserByToken(token string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(MATCH_ANY_DOMAIN, FIELD_TOKEN, token)
	s.lastError = err
	return rec, err
}

// Fetch a user by the email. Emails are not unique, except within a domain.
func (s *Store) FetchUserByEmail(domain, email string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(domain, FIELD_EMAIL, email)
	s.lastError = err
	return rec, err
}

// Fetch the user record by the login string. Login names are only unique within the domain
func (s *Store) FetchUserByLogin(domain, loginName string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(domain, FIELD_LOGIN, loginName)
	s.lastError = err
	return rec, err
}
