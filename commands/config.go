package commands

import (
	"L3afMe/Krul/kdgr"

	"github.com/sirkon/go-format"
)

func loadConfigCommands(r *kdgr.Route) {
	r.On("prefix", func(c *kdgr.Context) {
		if len(c.Args) == 0 {
			c.ReplyAutoHandle(kdgr.NewMessage("Prefix").
				Desc(format.Formatp("Current prefix: `${}`", c.Route.Config.Prefix)))

			return
		}

		oldPrefix := c.Route.Config.Prefix
		c.Route.Config.Prefix = c.Args.All("")
		c.Route.Config.Save()

		c.ReplyAutoHandle(kdgr.NewMessage("Prefix").
			Desc(format.Formatp("Old prefix: `${}`\nNew prefix: `${}`",
				oldPrefix, c.Route.Config.Prefix)))
	}).
		Desc("Set the prefix to execute commands."+
			"If no args are given then the current prefix will be displayed").
		Arg("prefix", "New prefix to set", false, kdgr.RouteArgString).
		Example("", "~")
}

func LoadConfig(r *kdgr.Route) {
	r.Group(func(r *kdgr.Route) {
		r.Cat("Config")

		loadConfigCommands(r)
	})
}
