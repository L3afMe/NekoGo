package commands

import (
	"L3afMe/NekoGo/router"

	"github.com/sirkon/go-format"
)

func cfgPrefix(c *router.Context) {
	msg := router.NewMessage("Prefix")

	if len(c.Args) == 0 {
		msg.Desc(format.Formatp(
			"Current prefix: `${}`",
			c.Route.Config.Prefix,
		))
	} else {
		oldPrefix := c.Route.Config.Prefix
		c.Route.Config.Prefix = c.Args.Get(0).AsString()
		c.Route.Config.Save()

		msg.Desc(format.Formatp(
			"Old prefix: `${}`\nNew prefix: `${}`",
			oldPrefix, c.Route.Config.Prefix,
		))
	}

	c.ReplyAutoHandle(msg)
}

func LoadConfig(r *router.Route) {
	r.Group(func(r *router.Route) {
		r.Cat("Config")

		cPrefix := r.On("prefix", cfgPrefix)
		cPrefix.Desc("Set the prefix to execute commands." +
			"If no args are given then the current prefix will be displayed.")
		cPrefix.Arg("prefix", "New prefix to set", false, router.ArgString)
		cPrefix.Example("", "~")
	})
}
