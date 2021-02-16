package kdgr

import (
	"L3afMe/Krul/utils"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirkon/go-format"
)

var (
	reID      = regexp.MustCompile(`\d{18}`)
	reUser    = regexp.MustCompile(`<@!?\d{18}>`)
	reChannel = regexp.MustCompile(`<#\d{18}>`)

	ErrIDRegexNotFound = errors.New("unable to get ID from argument")
)

// ArgError is a helper that makes returning the index and route easier if there are invalid args
type ArgError struct {
	Index  int
	Reason string
}

// Arg is a helper type to make using args easier
type Arg string

// IsUser returns if the argument is a mention or user id
func (a Arg) IsUser() bool {
	return (len(a.AsString()) == 18 && reID.MatchString(a.AsString())) ||
		(len(a.AsString()) <= 22 && reUser.MatchString(a.AsString()))
}

// IsChannel returns if the argument is a channel or channel id
func (a Arg) IsChannel() bool {
	return (len(a.AsString()) == 18 && reID.MatchString(a.AsString())) ||
		(len(a.AsString()) <= 22 && reChannel.MatchString(a.AsString()))
}

// IsInteger returns if the value can be parsed to an int
func (a Arg) IsInteger() bool {
	_, err := strconv.Atoi(a.AsString())
	return err == nil
}

// AsString returns the argument as a string
func (a Arg) AsString() string {
	return string(a)
}

// AsUser returns the discordgo.User that was mentioned
func (a Arg) AsUser(ses *discordgo.Session) (*discordgo.User, error) {
	channelID := reID.Find([]byte(a))

	if len(channelID) == 0 {
		return nil, ErrIDRegexNotFound
	}

	return ses.User(string(channelID))
}

// AsChannel returns the discordgo.Channel that was mentioned
func (a Arg) AsChannel(ses *discordgo.Session) (*discordgo.Channel, error) {
	channelID := reID.Find([]byte(a))

	if len(channelID) == 0 {
		return nil, ErrIDRegexNotFound
	}

	return ses.Channel(string(channelID))
}

// AsInteger returns the value parsed to an int
func (a Arg) AsInteger() int {
	val, _ := strconv.Atoi(a.AsString())
	return val
}

func (a Arg) GetEstType() RouteArgumentType {
	if a.IsUser() {
		return RouteArgUser
	} else if a.IsChannel() {
		return RouteArgChannel
	}
	return RouteArgString
}

// Args is an array of Arguments
type Args []Arg

// Get returns the argument at index n
func (a Args) Get(n int) Arg {
	if n >= 0 && n < len(a) {
		return a[n]
	}
	return ""
}

func (a Args) AsStrings() []string {
	strArgs := []string{}
	for _, arg := range a {
		strArgs = append(strArgs, arg.AsString())
	}

	return strArgs
}

func (a Args) All(separator string) string {
	return strings.Join(a.AsStrings(), separator)
}

func (a Args) Before(index int, separator string) string {
	if index >= 0 && index < len(a) {
		strArgs := []string{}
		for _, arg := range a[:index] {
			strArgs = append(strArgs, arg.AsString())
		}
		return strings.Join(strArgs, separator)
	}
	return ""
}

func (a Args) From(index int, separator string) string {
	if index >= 0 && index < len(a) {
		strArgs := []string{}
		for _, arg := range a[index:] {
			strArgs = append(strArgs, arg.AsString())
		}
		return strings.Join(strArgs, separator)
	}
	return ""
}

func ParseArgsNoCheck(content string, rt *Route) Args {
	return newArgs(strings.Split(content, rt.Separator)...)
}

// ParseArgs parses command arguments
func ParseArgs(content string, rt *Route) (Args, *ArgError) {
	args := newArgs(strings.Split(content, rt.Separator)...)

	minArgs := 0
	maxArgs := 0
	for i, r := range rt.Args {
		if r.Required {
			minArgs++
		}

		maxArgs++
		if i+1 == len(rt.Args) && strings.HasSuffix(r.Name, "...") {
			// Basically unlimited as Discord limits messages to 2000 characters
			maxArgs = 2000
		}
	}

	if minArgs > len(args) {
		argWord := "arguments"
		if minArgs == 1 {
			argWord = "argument"
		}

		return nil, &ArgError{-1, format.Formatp("Minimum of `${}` ${}, got `${}`", minArgs, argWord, len(args))}
	}

	if maxArgs < len(args) {
		argWord := "arguments"
		if maxArgs == 1 {
			argWord = "argument"
		}

		return nil, &ArgError{-1, format.Formatp("Maximum of `${}` ${}, got `${}`", maxArgs, argWord, len(args))}
	}

	for i, arg := range args {
		rtArg := rt.Args[utils.Min(i, len(rt.Args)-1)]

		if !checkArg(arg, rtArg) {
			return nil, &ArgError{i, format.Formatp("Expected type `${}`, got `${}`", rtArg.Type, arg.GetEstType())}
		}
	}

	return args, nil
}

func newArgs(strArgs ...string) Args {
	args := []Arg{}

	for _, arg := range strArgs {
		if arg != "" {
			args = append(args, Arg(arg))
		}
	}

	return args
}

func checkArg(arg Arg, rtArg RouteArgument) (valid bool) {
	valid = false

	switch rtArg.Type {
	case RouteArgString:
		valid = true
	case RouteArgUser:
		valid = arg.IsUser()
	case RouteArgChannel:
		valid = arg.IsChannel()
	case RouteArgInteger:
		valid = arg.IsInteger()
	}
	return
}
