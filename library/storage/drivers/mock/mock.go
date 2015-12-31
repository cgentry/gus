// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package mock

import (
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/library/storage"
	"github.com/cgentry/gus/record/tenant"
)

type MockDriver struct{}

type MockConn struct {
	db      map[string]*tenant.User
	errList map[string]error
}

// Fetch a raw database Mock driver
func NewMockDriver() *MockDriver {
	return &MockDriver{}
}

// The main driver will call this function to get a connection to the SqlLite db driver.
// it then 'routes' calls through this connection.
func (t *MockDriver) Open(option1 string, extraDriverOptions string) (storage.Conn, error) {
	store := &MockConn{}
	store.db = make(map[string]*tenant.User)
	store.errList = make(map[string]error)
	return store, nil
}

// Return the raw database handle to the caller. This allows more flexible options
func (t *MockConn) GetRawHandle() interface{} {
	return t.db
}

// Close the connection to the database (if it is open)
func (t *MockConn) Close() error {
	return nil
}

func (t *MockConn) UserUpdate(user *tenant.User) error {
	if err, ok := t.errList[user.Guid]; ok {
		return err
	}
	t.db[user.Guid] = user
	return nil
}
func (t *MockConn) UserInsert(user *tenant.User) error {
	if err, ok := t.errList[user.Guid]; ok {
		return err
	}
	t.db[user.Guid] = user
	return nil
}

func (t *MockConn) UserFetch(domain, key, value string) (*tenant.User, error) {
	found := false
	for _, user := range t.db {

		if domain == storage.MATCH_ANY_DOMAIN || domain == user.Domain {
			switch key {
			case storage.FIELD_GUID:
				found = (value == user.Guid)
			case storage.FIELD_EMAIL:
				found = (value == user.Email)
			case storage.FIELD_LOGIN:
				found = (value == user.LoginName)
			case storage.FIELD_TOKEN:
				found = (value == user.Token)
			}
			if found {
				if err, ok := t.errList[user.Guid]; ok {
					return nil, err
				}
				return user, nil
			}
		}
	}
	return nil, ErrUserNotFound
}
