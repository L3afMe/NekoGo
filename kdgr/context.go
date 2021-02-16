package kdgr

import (
	"sync"

	"github.com/op/go-logging"

	"github.com/bwmarrin/discordgo"
)

var MainGroup = "main"
var log = logging.MustGetLogger("Command")

type Context struct {
	Route *Route
	Msg   *discordgo.Message
	Ses   *discordgo.Session

	Args Args
	Log  *logging.Logger

	vmu  sync.RWMutex
	Vars map[string]interface{}
}

func (c *Context) Set(key string, d interface{}) {
	c.vmu.Lock()
	c.Vars[key] = d
	c.vmu.Unlock()
}

func (c *Context) Get(key string) interface{} {
	if c, ok := c.Vars[key]; ok {
		return c
	}
	return nil
}

func NewContext(s *discordgo.Session, m *discordgo.Message, args Args, route *Route) *Context {
	return &Context{
		Route: route,
		Msg:   m,
		Ses:   s,
		Args:  args,
		Log:   log,
		Vars:  map[string]interface{}{},
	}
}
