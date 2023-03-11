package chat

import (
	"encoding/json"
	"time"

	"github.com/BeanLiu1994/tiny_chat/session"

	lru "github.com/hashicorp/golang-lru/v2"
)

var DefaultChatManager = NewChatManager(&session.DefaultSessionManager)

func NewChatManager(sessMgr *session.SessionManager) *ChatManager {
	l, _ := lru.New[string, any](10)
	return &ChatManager{
		sessMgr: sessMgr,
		cache:   l,
	}
}

type ChatManager struct {
	sessMgr *session.SessionManager
	cache   *lru.Cache[string, any]
}

func (c *ChatManager) OnConnect(who string) {
	list := c.getCache()
	if len(list) == 0 {
		return
	}
	sess := c.sessMgr.Get(who)
	if sess == nil {
		return
	}
	for _, v := range list {
		sess.SendString(v)
	}
}
func (c *ChatManager) OnDisconnect(who string) {
}

func (c *ChatManager) Say(from, to, what string) {
	m := Message{
		Sender:   from,
		SentTime: time.Now(),
		Content:  what,
	}
	if to == "" {
		c.Broadcast(m)
		return
	}
	b, _ := json.Marshal(m)
	toSess := c.sessMgr.Get(to)
	fromSess := c.sessMgr.Get(from)
	if toSess != nil && fromSess != nil {
		toSess.Send(b)
		fromSess.Send(b)
	} else {
		session.SendJsonErr(fromSess, "target user is offline")
	}
}

func (c *ChatManager) Broadcast(m Message) {
	b, _ := json.Marshal(m)
	c.addCache(string(b))
	c.sessMgr.ForEach(func(name string, sess session.SessionInterface) {
		sess.Send(b)
	})
}

func (c *ChatManager) addCache(msg string) {
	if c.cache == nil {
		return
	}
	c.cache.Add(msg, nil)
}

func (c *ChatManager) clearCache() {
	if c.cache == nil {
		return
	}
	c.cache.Purge()
}

func (c *ChatManager) getCache() []string {
	if c.cache == nil {
		return nil
	}
	return c.cache.Keys()
}

// cache public message only
type Message struct {
	Sender   string    `json:"sender"`
	SentTime time.Time `json:"sent_time"`
	Content  string    `json:"content"`
}
