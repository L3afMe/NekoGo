package commands

import (
	"L3afMe/Krul/kdgr"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirkon/go-format"
	"github.com/valyala/fasthttp"
)

type neko8BallResponse struct {
	URL      string `json:"url"`
	Response string `json:"response"`
}

func fun8Ball(c *kdgr.Context) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURIBytes([]byte("https://nekos.life/api/v2/8ball"))
	req.Header.SetMethodBytes([]byte("GET"))

	err := fasthttp.Do(req, resp)
	if err != nil {
		c.ReplyAutoHandle(kdgr.NewError(format.Formatp("Unable to get image: ${}", err)))
	}
	bodyStr := resp.Body()

	fasthttp.ReleaseResponse(resp)
	var res *neko8BallResponse
	err = json.Unmarshal(bodyStr, &res)
	if err != nil {
		c.ReplyAutoHandle(kdgr.NewError(format.Formatp("Unable to unmarshal response: ${}", err)))
	}

	msg := kdgr.NewMessage("8 Ball").
		Desc(format.Formatp("You ask \"${}\" and the 8ball responds with \"${}\"", c.Args.All(" "), res.Response)).
		Thumbnail(res.URL)

	c.ReplyAutoHandle(msg)
}

func funDick(ctx *kdgr.Context) {
	var usr *discordgo.User
	if len(ctx.Args) == 0 {
		usr = ctx.Msg.Author
	} else {
		var err error
		usr, err = ctx.Args.Get(0).AsUser(ctx.Ses)
		if err != nil {
			ctx.ReplyInvalidArg(0, "Invalid user specified")
			return
		}
	}

	usrID, err := strconv.Atoi(usr.ID)
	if err != nil {
		ctx.ReplyAutoHandle(kdgr.NewError("Unable to convert user ID to int"))
		return
	}

	size := ((usrID / 75) % 14) + 4
	ctx.Log.Info(usrID, size)
	header := format.Formatp("${} has a ${}\" dick", usr.Mention(), size)
	dick := format.Formatp("8${}D", strings.Repeat("=", size))

	ctx.ReplyAutoHandle(kdgr.NewMessage("Compatibility").Desc(header + "\n" + dick))
}

func funShip(ctx *kdgr.Context) {
	usr1, err := ctx.Args.Get(0).AsUser(ctx.Ses)
	if err != nil {
		ctx.ReplyInvalidArg(0, "Invalid user specified")
		return
	}

	var usr2 *discordgo.User
	if len(ctx.Args) == 1 {
		usr2 = ctx.Msg.Author
	} else {
		usr2, err = ctx.Args.Get(1).AsUser(ctx.Ses)
		if err != nil {
			ctx.ReplyInvalidArg(0, "Invalid user specified")
			return
		}
	}

	usr1ID, err1 := strconv.Atoi(usr1.ID)
	usr2ID, err2 := strconv.Atoi(usr2.ID)
	if err1 != nil || err2 != nil {
		ctx.ReplyAutoHandle(kdgr.NewError("Unable to convert user ID to int"))
		return
	}

	compat := ((usr1ID + usr2ID) / 58) % 100
	shipName := usr1.Username[:len(usr1.Username)/2] + usr2.Username[len(usr2.Username)/2:]

	header := format.Formatp("${} has ${}% compatibility", shipName, compat)
	bar := format.Formatp("[${}${}]", strings.Repeat("=", compat/5), strings.Repeat("-", 20-compat/5))

	ctx.ReplyAutoHandle(kdgr.NewMessage("Compatibility").Desc(header + "\n" + bar))
}

func loadFunCommands(r *kdgr.Route) {
	r.On("8ball", fun8Ball).
		Desc("Ask the magic 8ball a question").
		Example("Should I get out of bed?").
		Arg("question...", "The question to ask the 8ball", false, kdgr.RouteArgString)

	r.On("dick", funDick).
		Desc("Check the size of a user").
		Example("@l3af#0001").
		Arg("user", "User to check the size of", false, kdgr.RouteArgUser)

	r.On("compatibility", funShip).
		Alias("ship", "compat").
		Desc("Check the compatibility between two users").
		Example("@l3af#0001").
		Arg("user", "User to check compatibility with", true, kdgr.RouteArgUser).
		Arg("user", "User to check compatibility with", false, kdgr.RouteArgUser)
}

func LoadFun(r *kdgr.Route) {
	r.Cat("Fun")

	loadFunCommands(r)
}
