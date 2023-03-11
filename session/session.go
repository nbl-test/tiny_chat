package session

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
)

var DefaultSessionManager SessionManager

// SessionManager should manage connected sessions, and provide send method
type SessionManager struct {
	// stores  generated_id(string) -> session_interface(SessionInterface)
	sessions sync.Map
}

// get by name.
// nil means not found
func (s *SessionManager) Get(name string) SessionInterface {
	v, found := s.sessions.Load(name)
	if !found {
		return nil
	}
	vv, ok := v.(SessionInterface)
	if !ok {
		return nil
	}
	return vv
}

// create a name
func (s *SessionManager) CreateID() string {
	return uuid.New().String()
}

// set
func (s *SessionManager) Set(name string, sess SessionInterface) {
	if sess == nil {
		s.sessions.Delete(name)
		return
	}
	s.sessions.Store(name, sess)
}

// add
func (s *SessionManager) Add(sess SessionInterface) string {
	name := s.CreateID()
	s.Set(name, sess)
	return name
}

// for each item
func (s *SessionManager) ForEach(iterator func(name string, sess SessionInterface)) {
	if iterator == nil {
		return
	}
	s.sessions.Range(func(key, value any) bool {
		k, ok := key.(string)
		if !ok {
			return true
		}
		v, ok := value.(SessionInterface)
		if !ok {
			return true
		}
		iterator(k, v)
		return true
	})
}

// must implements these
type SessionInterface interface {
	GetID() string
	Created() time.Time
	Send([]byte)
	SendString(string)
}

func SendJsonErr(sess SessionInterface, err string) {
	if sess == nil {
		return
	}
	b, _ := json.Marshal(map[string]string{
		"error": err,
	})
	sess.Send(b)
}
