package session

import (
"fmt"
"log"

"github.com/gin-gonic/gin"
)

type SessionMgrType string

const (
	// SessionID在cookie里面的名字
	SessionCookieName = "session_id"
	// Session对象在Context里面的名字
	SessionContextName                = "session"
	Memory             SessionMgrType = "memory"
	Redis              SessionMgrType = "redis"
)

// Session 接口
type Session interface {
	// 获取Session对象的ID
	ID() string
	// 加载redis数据到 session data
	Load() error
	// 获取key对应的value值
	Get(string) (interface{}, error)
	// 设置key对应的value值
	Set(string, interface{})
	// 删除key对应的value值
	Del(string)
	// 落盘数据到redis
	Save()
	// 设置Redis数据过期时间,内存版本无效
	SetExpired(int)
}

// SessionMgr Session管理器对象
type SessionMgr interface {
	// 初始化Redis数据库连接
	Init(addr string, options ...string) error
	// 通过SessionID获取已经初始化的Session对象
	GetSession(string) (Session, error)
	// 创建一个新的Session对象
	CreateSession() Session
	// 使用SessionID清空一个Session对象
	Clear(string)
}

// Options Cookie对应的相关选项
type Options struct {
	Path   string
	Domain string
	// Cookie中的SessionID存活时间
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

func CreateSessionMgr(name SessionMgrType, addr string, options ...string) (sm SessionMgr, err error) {
	switch name {
	case Memory:
		sm = NewMemSessionMgr()
	case Redis:
		sm = NewRedisSessionMgr()
	default:
		err = fmt.Errorf("unsupported %v\n", name)
		return
	}
	err = sm.Init(addr, options...)
	return
}

func SessionMiddleware(sm SessionMgr, options Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		var session Session
		// 尝试从cookie获取session ID
		sessionID, err := c.Cookie(SessionCookieName)
		if err != nil {
			log.Printf("get session_id from cookie failed, err:%v\n", err)
			session = sm.CreateSession()
			sessionID = session.ID()
		} else {
			log.Printf("SessionId: %v\n", sessionID)
			session, err = sm.GetSession(sessionID)
			if err != nil {
				log.Printf("Get session by %s failed, err: %v\n", sessionID, err)
				session = sm.CreateSession()
				sessionID = session.ID()
			}
		}

		session.SetExpired(options.MaxAge)
		c.Set(SessionContextName, session)
		c.SetCookie(SessionCookieName, sessionID, options.MaxAge, options.Path, options.Domain, options.Secure, options.HttpOnly)
		defer sm.Clear(sessionID)
		c.Next()
	}
}
