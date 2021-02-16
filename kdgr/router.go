package kdgr

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirkon/go-format"
)

var (
	ErrCouldNotFindRoute  = errors.New("could not find route")
	ErrRouteAlreadyExists = errors.New("route already exists")
)

type HandlerFunc func(*Context)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type ExecBeforeFunc func(*Context) bool

type ExecAfterFunc func(*Context)

func (r *Route) On(name string, handler HandlerFunc) (rt *Route) {
	rt = r.OnMatch(name, nil, handler)
	rt.Matcher = NewNameMatcher(rt)
	return
}

func (r *Route) OnMatch(name string, matcher func(string) bool, handler HandlerFunc) (rt *Route) {
	if rt = r.Find(name); rt != nil {
		return
	}

	for _, v := range r.Middleware {
		handler = v(handler)
	}

	rt = &Route{
		Name:         name,
		Description:  "No description provided",
		Category:     r.Category,
		Handler:      handler,
		Matcher:      matcher,
		Availability: RouteBoth,
		Separator:    " ",
		Config:       r.Config,
	}

	_ = r.AddRoute(rt)

	return
}

func (r *Route) Use(fn ...MiddlewareFunc) *Route {
	r.Middleware = append(r.Middleware, fn...)
	return r
}

func (r *Route) Find(name string) (rt *Route) {
	for _, subRt := range r.Routes {
		if subRt.Matcher(name) {
			rt = subRt
			return
		}
	}

	return
}

func (r *Route) FindFull(args ...string) (route *Route, depth int) {
	route = r
	for _, arg := range args {
		if rt := route.Find(arg); rt != nil {
			route = rt
			depth++
		} else {
			return
		}
	}

	return
}

func mention(id string) string {
	return "<@" + id + ">"
}

func nickMention(id string) string {
	return "<@!" + id + ">"
}

func (r *Route) Group(fn func(r *Route)) *Route {
	rt := New(r.Config)
	fn(rt)
	for _, v := range rt.Routes {
		r.AddRoute(v)
	}
	return r
}

func (r *Route) FindAndExecute(s *discordgo.Session, prefix, botID string, m *discordgo.Message) {
	var pf string

	if r.Default != nil && m.Content == mention(botID) || r.Default != nil && m.Content == nickMention(botID) {
		r.Default.Handler(NewContext(s, m, newArgs(""), r.Default))
		return
	}

	bmention := mention(botID) + " "
	nmention := nickMention(botID) + " "

	switch {
	case prefix != "" && strings.HasPrefix(m.Content, prefix):
		pf = prefix
	case strings.HasPrefix(m.Content, bmention):
		pf = bmention
	case strings.HasPrefix(m.Content, nmention):
		pf = nmention
	default:
		return
	}

	command := strings.TrimPrefix(m.Content, pf)
	args := strings.Split(command, " ")
	if rt, depth := r.FindFull(args...); depth > 0 {
		strArgs := strings.Join(args[depth:], " ")
		args, err := ParseArgs(strArgs, rt)
		if err != nil {
			args = ParseArgsNoCheck(strArgs, rt)
			context := NewContext(s, m, args, rt)

			if r.ExecBefore(context) {
				context.ReplyInvalidArg(err.Index, err.Reason)
				r.ExecAfter(context)
			}
		} else {
			context := NewContext(s, m, args, rt)

			if r.ExecBefore(context) {
				if rt.Availability == RouteGuild && m.GuildID == "" ||
					rt.Availability == RouteDM && m.GuildID != "" {
					context.ReplyAutoHandle(NewError(format.
						Formatp("This command can only be used in ${}", rt.Availability)))
				} else {
					rt.Handler(context)
				}

				r.ExecAfter(context)
			}
		}
	}
}

func (r *Route) AddRoute(route *Route) error {
	if rt := r.Find(route.Name); rt != nil {
		return ErrRouteAlreadyExists
	}

	route.Parent = r
	r.Routes = append(r.Routes, route)
	return nil
}

func (r *Route) RemoveRoute(route *Route) error {
	for i, v := range r.Routes {
		if v == route {
			r.Routes = append(r.Routes[:i], r.Routes[i+1:]...)
			return nil
		}
	}
	return ErrCouldNotFindRoute
}
