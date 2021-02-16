package kdgr

import (
	"L3afMe/Krul/config"
)

type RouteArgumentType string
type RouteAvailabilityType string

const (
	RouteArgString  RouteArgumentType = "string"
	RouteArgInteger RouteArgumentType = "integer"
	RouteArgUser    RouteArgumentType = "user"
	RouteArgChannel RouteArgumentType = "channel"

	RouteBoth  RouteAvailabilityType = "DMs and Guilds"
	RouteDM    RouteAvailabilityType = "DMs"
	RouteGuild RouteAvailabilityType = "Guilds"
)

type RouteArgument struct {
	Name        string
	Description string
	Required    bool
	Type        RouteArgumentType
}

type Route struct {
	Routes       []*Route
	Name         string
	Aliases      []string
	Description  string
	Args         []RouteArgument
	Examples     []string
	Category     string
	Permissions  int
	Availability RouteAvailabilityType
	Separator    string
	Matcher      func(string) bool
	Handler      HandlerFunc
	Default      *Route
	Parent       *Route
	Middleware   []MiddlewareFunc
	ExecBefore   ExecBeforeFunc
	ExecAfter    ExecAfterFunc
	Config       *config.Config
}

func New(config *config.Config) *Route {
	return &Route{
		Routes:   []*Route{},
		Category: "",
		Config:   config,
	}
}

func (r *Route) Desc(description string) *Route {
	r.Description = description
	return r
}

func (r *Route) Cat(category string) *Route {
	r.Category = category
	return r
}

func (r *Route) Alias(aliases ...string) *Route {
	r.Aliases = append(r.Aliases, aliases...)
	return r
}

func (r *Route) Arg(name, description string, required bool, argType RouteArgumentType) *Route {
	r.Args = append(r.Args, RouteArgument{name, description, required, argType})
	return r
}

func (r *Route) Example(examples ...string) *Route {
	r.Examples = append(r.Examples, examples...)
	return r
}

func (r *Route) Perms(permissions ...int) *Route {
	for _, perm := range permissions {
		r.Permissions = r.Permissions | perm
	}

	return r
}

func (r *Route) Where(where RouteAvailabilityType) *Route {
	r.Availability = where
	return r
}

func (r *Route) Sep(separator string) *Route {
	r.Separator = separator
	return r
}

func (r *Route) Before(before ExecBeforeFunc) *Route {
	r.ExecBefore = before
	return r
}

func (r *Route) After(after ExecAfterFunc) *Route {
	r.ExecAfter = after
	return r
}

func (r *Route) GetFullName() string {
	parent := r
	fullName := parent.Name
	for {
		if parent.Parent == nil || parent.Parent.Name == "" {
			break
		}

		parent = parent.Parent
		fullName = parent.Name + " " + fullName
	}

	return fullName
}

func (r *Route) GetRootParent() *Route {
	parent := r
	for {
		if parent.Parent == nil || parent.Parent.Name == "" {
			break
		}

		parent = parent.Parent
	}

	return parent
}
