package kdgr

import (
	"L3afMe/Krul/utils"
	"strings"

	"github.com/sirkon/go-format"
)

func (c *Context) ReplyInvalidArg(idx int, reason string) {
	msg := NewError(format.Formatp("Invalid arguments: ${}", reason)).
		AddField("Input", "```\n"+c.Route.GetFullName()+" "+c.Args.From(0, " ")+"\n```", false)

	if len(c.Route.Args) > 0 {
		args := c.Route.GetFullName() + " "
		btmStr := ""

		errPos := 0
		var errArg RouteArgument
		for i, arg := range c.Route.Args {
			if i == idx {
				errArg = arg
				errPos = len(args) + 1
			}

			if arg.Required {
				args = args + "<" + arg.Name + "> "
			} else {
				args = args + "[" + arg.Name + "] "
			}
		}

		if errPos > 0 {
			arrows := strings.Repeat(" ", errPos) + strings.Repeat("^", len(errArg.Name))

			//nolint:gocritic // Can't convert to switch statement
			if len(arrows+" "+errArg.Description) <= len(args) {
				btmStr = arrows + " " + errArg.Description
			} else if errPos+len(errArg.Description) <= len(args) {
				btmStr = arrows + "\n" + strings.Repeat(" ", errPos) + errArg.Description
			} else if len(errArg.Description) <= len(args) {
				btmStr = arrows + "\n" + strings.Repeat(" ", len(args)-len(errArg.Description)) + errArg.Description
			} else {
				btmStr = arrows + "\n  " + errArg.Description
			}
		}

		if len(strings.Trim(btmStr, " ")) > 0 {
			args = args + "\n" + btmStr
		}
		msg.AddField("Expected", "```\n"+args+"\n```", false)
	}

	c.ReplyAutoHandle(msg)
}

func SendHelp(c *Context) {
	r := c.Route.Parent

	cats := []string{}
	for _, v := range r.Routes {
		if !utils.SAContains(cats, v.Category) {
			cats = append(cats, v.Category)
		}
	}

	if len(c.Args) == 0 {
		helpMsg := NewMessage("Help").
			Desc("Use `help ?` for more information").
			AddField("Categories", strings.Join(cats, "\n"), false)
		c.ReplyAutoHandle(helpMsg)
		return
	}

	if rt, depth := r.FindFull(c.Args.AsStrings()...); depth > 0 {
		sendRouteHelp(c, rt)
		return
	}

	cat := c.Args.Get(0).AsString()
	if utils.SAContainsCI(cats, cat) {
		sendCategoryHelp(c, cat, r)
		return
	}

	c.ReplyAutoHandle(NewMessage("Help").
		Desc("Unable to find command or category named `" + c.Args.All(" ") + "`"))
}

func sendRouteHelp(c *Context, rt *Route) {
	argStr := "`"
	for _, arg := range rt.Args {
		if arg.Required {
			argStr += "<" + arg.Name + "> "
		} else {
			argStr += "[" + arg.Name + "] "
		}
	}
	argStr += "`"

	parents := ""
	parent := rt.Parent
	for {
		if parent == nil || parent.Category == MainGroup {
			break
		}

		parents = parents + parent.Name + " "
		parent = parent.Parent
	}

	exmpStr := ""
	if len(rt.Examples) == 0 {
		exmpStr = "`" + parents + rt.Name + "`"
	} else {
		exmpStr = "`" + parents + rt.Name + " " + strings.Join(rt.Examples, "`\n`"+parents+rt.Name+" ") + "`"
	}

	subStr := ""
	if len(rt.Routes) > 0 {
		subs := []string{}
		for _, rtt := range rt.Routes {
			subs = append(subs, rtt.Name)
		}

		subStr = "`" + strings.Join(subs, "`, `") + "`"
	}

	helpMsg := NewMessage("Help").
		AddField(rt.Name, rt.Description, false)

	if len(argStr) > 2 {
		helpMsg.AddField("Arguments", argStr, true)
	}

	helpMsg.AddField("Examples", exmpStr, true)

	if len(rt.Aliases) > 0 {
		helpMsg.AddField("Aliases", "`"+strings.Join(rt.Aliases, "`, `")+"`", true)
	}

	if rt.Parent.Name != MainGroup && len(strings.Trim(rt.Parent.Name, " ")) > 0 {
		helpMsg.AddField("Parent Command", rt.Parent.Name, true)
	}

	if len(subStr) > 0 {
		helpMsg.AddField("Subcommands", subStr, true)
	}

	if rt.Category != MainGroup {
		helpMsg.AddField("Category", rt.Category, true)
	}

	c.ReplyAutoHandle(helpMsg)
}

func sendCategoryHelp(c *Context, category string, route *Route) {
	routes := []string{}

	capCat := ""
	for _, v := range route.Routes {
		if strings.EqualFold(v.Category, category) {
			routes = append(routes, v.Name+" - "+strings.Split(v.Description, ".")[0])
			capCat = v.Category
		}
	}

	c.ReplyAutoHandle(NewMessage("Help").
		AddField(capCat, strings.Join(routes, "\n"), false))
}
