package storage

import (
	"github.com/cgentry/gdriver"
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/tenant"
)

// These are the names of fields we expect to occur in the database and will
// pass to database functions when performing UserFetch operations. You may
// map them in the driver-level routines in order to provide names that are
// more appropriate to the driver mechanism.
const (
	FieldEmail = `Email`
	FieldName  = `FullName`
	FieldGUID  = `Guid`
	FieldLogin = `LoginName`
	FieldToken = `Token`
)

// MatchAnyDomain is a special character that should be used to search ALL domains.
const MatchAnyDomain = "*"

/* ======================================================
 * 					Store functions
 * ======================================================
 */

// Id will return the identity of the storage driver
func (s *Store) Id() string {
	return s.rawDriver.Identity(gdriver.IdentityName)
}

// ShortHelp will return a brief description of the storage driver.
func (s *Store) ShortHelp() string {
	return s.rawDriver.Identity(gdriver.IdentityShort)
}

// LongHelp will return a long description of the storage driver. It should
// give a fair amount of detail
func (s *Store) LongHelp() string {
	return s.rawDriver.Identity(gdriver.IdentityLong)
}

// GetStorageConnector will return the actual connection to the database for low-level access.
// This should be avoided unless you are coding for a very non-portable
// function
func (s *Store) GetStorageConnector() Conn {
	return s.connection
}

/*
 * The following functions are provided by this class and are not
 * encapsulated
 */

// LastError returns the last known error condition that was given by a call
func (s *Store) LastError() error {
	return s.lastError
}

// SetLastError will save the last error that occured in this class
func (s *Store) SetLastError(err error) *Store {
	s.lastError = err
	return s
}

// IsOpen will return the the open status of the connection
func (s *Store) IsOpen() bool {
	return s.isOpen
}

// saveAndReturnError is used internally to save the erro but also to return it back to caller
func (s *Store) saveAndReturnError(err error) error {
	s.lastError = err
	return err
}

// ClearErrors simply clears the last stored error
func (s *Store) ClearErrors() {
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
		return s.saveAndReturnError(ErrAlreadyOpen)
	}
	s.isOpen = false
	s.lastError = nil
	s.connectString = connect
	s.connection, s.lastError = s.driver.Open(connect, extraDriverOptions)

	return s.lastError
}

// UserUpdate will update the 'tenant record in the database. It makes sure the
// database is open to stop any problems with low-level drivers
func (s *Store) UserUpdate(user *tenant.User) error {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return ErrNotOpen
	}
	return s.saveAndReturnError(s.connection.UserUpdate(user))
}

// UserInsert will attempt to insert a new 'tenant' record in the datase. It makes sure
// the database is open and will save any error code that occurs.
func (s *Store) UserInsert(user *tenant.User) error {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return ErrNotOpen
	}
	return s.saveAndReturnError(s.connection.UserInsert(user))
}

// UserFetch will find a tenant using the domain, a field name and the field value. There will only be one record
// returned. If you pass MatchAnyDomain as the domain, this will only be valid for a small number of
// key-types (e.g. enforced unique keys.)
func (s *Store) UserFetch(domain, lookupKey, lookkupValue string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	if domain == MatchAnyDomain {
		if lookupKey != FieldGUID || lookupKey != FieldToken {
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
		return s.saveAndReturnError(ErrNotOpen)
	}
	s.isOpen = false
	s.ClearErrors()
	if closer, found := s.connection.(Closer); found {
		s.lastError = closer.Close()
	}
	return s.lastError
}

// CreateStore , if implemented, initialises any storage. If not implemented, an error will be returned.
func (s *Store) CreateStore() error {
	if s.isOpen != true {
		return s.saveAndReturnError(ErrNotOpen)
	}

	if creater, found := s.connection.(Creater); found {
		return s.saveAndReturnError(creater.CreateStore())
	}
	return ErrNoSupport
}

// Ping will test to see if the storage system is alive. This is an optional routine. If it
// doesn't exist, a nil return occurs (no error)
func (s *Store) Ping() error {
	s.ClearErrors()
	if pinger, found := s.connection.(Pinger); found {
		s.lastError = pinger.Ping()
	}
	return s.lastError
}

/* ------------------------ THE FOLLOWING ARE 'CONVENIENCE' FUNCTIONS ***********************/

// FetchUserByGUID No domains are required as this is the primary (or unique) key
func (s *Store) FetchUserByGUID(guid string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(MatchAnyDomain, FieldGUID, guid)
	s.lastError = err
	return rec, err
}

// FetchUserByToken If the user is not logged in, a 'User not found' error is returned.
func (s *Store) FetchUserByToken(token string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(MatchAnyDomain, FieldToken, token)
	s.lastError = err
	return rec, err
}

// FetchUserByEmail Emails are not unique, except within a domain.
func (s *Store) FetchUserByEmail(domain, email string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(domain, FieldEmail, email)
	s.lastError = err
	return rec, err
}

// FetchUserByLogin Login names are only unique within the domain
func (s *Store) FetchUserByLogin(domain, loginName string) (*tenant.User, error) {
	if !s.isOpen {
		s.lastError = ErrNotOpen
		return nil, ErrNotOpen
	}
	rec, err := s.connection.UserFetch(domain, FieldLogin, loginName)
	s.lastError = err
	return rec, err
}
