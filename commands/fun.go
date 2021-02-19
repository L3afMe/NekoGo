package commands

import (
	"L3afMe/NekoGo/router"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/sirkon/go-format"
	"github.com/valyala/fasthttp"
)

type neko8BallResponse struct {
	URL      string `json:"url"`
	Response string `json:"response"`
}

func fun8Ball(c *router.Context) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURIBytes([]byte("https://nekos.life/api/v2/8ball"))
	req.Header.SetMethodBytes([]byte("GET"))

	if err := fasthttp.Do(req, resp); err != nil {
		msgErr := router.NewError(format.Formatp(
			"Unable to get image: ${}",
			err,
		))

		c.ReplyAutoHandle(msgErr)
	}

	bodyStr := resp.Body()
	fasthttp.ReleaseResponse(resp)

	var res *neko8BallResponse
	if err := json.Unmarshal(bodyStr, &res); err != nil {
		msgErr := router.NewError(format.Formatp(
			"Unable to unmarshal response: ${}",
			err,
		))
		c.ReplyAutoHandle(msgErr)
	}

	msg := router.NewMessage("8 Ball")
	msg.Desc(format.Formatp(
		"You ask \"${}\" and the 8ball responds with \"${}\"",
		c.Args.All(" "),
		res.Response,
	))
	msg.Thumbnail(res.URL)

	c.ReplyAutoHandle(msg)
}

func funDick(c *router.Context) {
	usr := c.Msg.Author
	if len(c.Args) == 1 {
		var err error
		if usr, err = c.Args.Get(0).AsUser(c.Ses); err != nil {
			c.ReplyInvalidArg(0, "Invalid user specified.")

			return
		}
	}

	usrID, err := strconv.Atoi(usr.ID)
	if err != nil {
		msgErr := router.NewError("Unable to convert user ID to int.")
		c.ReplyAutoHandle(msgErr)

		return
	}

	size := (usrID/75)%14 + 4
	header := format.Formatp(
		"${} has a ${}\" dick",
		usr.Mention(),
		size,
	)
	dick := format.Formatp(
		"8${}D",
		strings.Repeat("=", size),
	)

	msg := router.NewMessage("Dick")
	msg.Desc(format.Formatp(
		"${}\n${}",
		header,
		dick,
	))

	c.ReplyAutoHandle(msg)
}

func funShip(c *router.Context) {
	usr1, err := c.Args.Get(0).AsUser(c.Ses)
	if err != nil {
		c.ReplyInvalidArg(0, "Invalid user specified.")

		return
	}

	usr2 := c.Msg.Author
	if len(c.Args) == 2 {
		if usr2, err = c.Args.Get(1).AsUser(c.Ses); err != nil {
			c.ReplyInvalidArg(1, "Invalid user specified.")

			return
		}
	}

	usr1ID, err1 := strconv.Atoi(usr1.ID)
	usr2ID, err2 := strconv.Atoi(usr2.ID)
	if err1 != nil || err2 != nil {
		c.ReplyAutoHandle(router.NewError("Unable to convert user ID to int."))
		return
	}

	compat := ((usr1ID + usr2ID) / 58) % 100
	shipName := usr1.Username[:len(usr1.Username)/2] + usr2.Username[len(usr2.Username)/2:]

	header := format.Formatp(
		"${} has ${}% compatibility",
		shipName,
		compat,
	)
	bar := format.Formatp(
		"[${}${}]",
		strings.Repeat("=", compat/5),
		strings.Repeat("-", 20-compat/5),
	)

	msg := router.NewMessage("Compatibility")
	msg.Desc(format.Formatp(
		"${}\n${}",
		header,
		bar,
	))

	c.ReplyAutoHandle(msg)
}

func LoadFun(r *router.Route) {
	r.Group(func(r *router.Route) {
		r.Cat("Fun")

		c8Ball := r.On("8ball", fun8Ball)
		c8Ball.Desc("Ask the magic 8ball a question.")
		c8Ball.Example("Should I get out of bed?")
		c8Ball.Arg("question...", "The question to ask the 8ball", true, router.ArgString)

		cDick := r.On("dick", funDick)
		cDick.Desc("Check the size of a user.")
		cDick.Example("@l3af#0001")
		cDick.Arg("user", "User to check the size of", false, router.ArgUser)

		cShip := r.On("compatibility", funShip)
		cShip.Alias("ship", "compat")
		cShip.Desc("Check the compatibility between two users.")
		cShip.Example("@l3af#0001")
		cShip.Arg("user", "User to check compatibility with", true, router.ArgUser)
		cShip.Arg("user", "User to check compatibility with", false, router.ArgUser)
	})
}
