// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package jsonfile

import (
	"encoding/json"
	"fmt"
	. "github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/library/storage"
	"github.com/cgentry/gus/record/tenant"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type JsonFileDriver struct{}

func New() *JsonFileDriver {
	return &JsonFileDriver{}
}

type JsonFileConn struct {
	filename  string
	filemod   time.Time
	filesize  int
	busy      sync.Mutex
	isdirty   bool
	isMonitor bool

	userlist map[string]*tenant.User

	messages chan *jsonMessage
}

func NewJsonFileConn(name string) *JsonFileConn {
	store := &JsonFileConn{
		filename:  name,
		isdirty:   false,
		isMonitor: true}
	store.messages = make(chan *jsonMessage, 10)
	store.userlist = make(map[string]*tenant.User)
	return store
}

// Fetch a raw database JsonFile driver
func NewJsonFileDriver() *JsonFileDriver {
	return &JsonFileDriver{}
}

const (
	CmdTimer = iota
	CmdNew
	CmdLoad
)

type jsonMessage struct {
	Command    int
	Parameters string
}

// The main driver will call this function to get a connection to the SqlLite db driver.
// it then 'routes' calls through this connection.
func (t *JsonFileDriver) Open(jsonfile string, extraDriverOptions string) (storage.Conn, error) {
	store := NewJsonFileConn(jsonfile)
	store.messages <- &jsonMessage{Command: CmdLoad}
	go store.Monitor() // Start the MONITOR in the background
	return store, nil
}

// Monitor will check the file for changes every minute. This will change to an notify routine once the
// functions are merged from experimental to main. This routine should work with any OS (rather than just one
// or two)
func (t *JsonFileConn) Monitor() {
	for {
		msg := <-t.messages
		t.busy.Lock()

		if msg.Command == CmdNew {
			if len(t.userlist) < 1 {
				continue
			}
			if buff, err := json.MarshalIndent(t.userlist, "", "  "); err == nil {
				err := ioutil.WriteFile(t.filename, buff, 0600)
				if err != nil {
					fmt.Println("ERROR! JsonFileConn: ", err.Error())
				}
				t.isdirty = false
			} else {
				fmt.Println("ERROR! JsonFileConn: ", err.Error())
			}
		} else if msg.Command == CmdLoad {
			if buff, err := ioutil.ReadFile(t.filename); err == nil {
				if len(buff) > 0 {
					err := json.Unmarshal(buff, &t.userlist)
					if err != nil {
						fmt.Println("ERROR! JsonFileConn: ", err.Error())
					}
				}
			}
		} else if msg.Command == CmdTimer {
			finfo, err := os.Stat(t.filename)
			if err == nil {
				if finfo.ModTime().After(t.filemod) {
					t.messages <- &jsonMessage{Command: CmdLoad}
				}
			}
		}
		if stat, err := os.Stat(t.filename); err == nil {
			t.filemod = stat.ModTime()
		}
		t.busy.Unlock()
	}
}

// Close the connection to the database (if it is open)
func (t *JsonFileConn) Close() error {
	return nil
}

func (t *JsonFileConn) UserUpdate(userRecord *tenant.User) error {
	t.busy.Lock()
	defer t.busy.Unlock()

	t.userlist[userRecord.Guid] = userRecord
	t.isdirty = true
	t.messages <- &jsonMessage{Command: CmdNew}
	return nil
}
func (t *JsonFileConn) UserInsert(userRecord *tenant.User) error {
	t.busy.Lock()
	defer t.busy.Unlock()
	userRecord.Id = len(t.userlist)
	t.userlist[userRecord.Guid] = userRecord
	t.isdirty = true
	t.messages <- &jsonMessage{Command: CmdNew}
	return nil
}

func (t *JsonFileConn) UserFetch(domain, key, value string) (*tenant.User, error) {
	found := false
	t.busy.Lock()
	defer t.busy.Unlock()

	for _, userRecord := range t.userlist {

		if domain == storage.MatchAnyDomain || domain == userRecord.Domain {
			switch key {
			case storage.FieldGUID:
				found = (value == userRecord.Guid)
			case storage.FieldEmail:
				found = (value == userRecord.Email)
			case storage.FieldLogin:
				found = (value == userRecord.LoginName)
			case storage.FieldToken:
				found = (value == userRecord.Token)
			}
			if found {
				return userRecord, nil
			}
		}
	}
	return nil, ErrUserNotFound
}
