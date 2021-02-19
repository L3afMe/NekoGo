package commands

import (
	"L3afMe/NekoGo/router"
	"encoding/json"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirkon/go-format"
	"github.com/valyala/fasthttp"
)

type interactionSite int

const (
	siteNekosLife interactionSite = iota
)

type interaction struct {
	Noun      string
	URL       string
	Responses []string
	Site      interactionSite
}

var (
	interactionsMap = map[string]interaction{
		"kiss": {
			Noun: "Kiss",
			URL:  "kiss",
			Responses: []string{
				"${0} kissed ${1}",
			},
			Site: siteNekosLife,
		},
		"slap": {
			Noun: "Slap",
			URL:  "slap",
			Responses: []string{
				"${0} slapped ${1}",
			},
			Site: siteNekosLife,
		},
		"cuddle": {
			Noun: "Cuddle",
			URL:  "cuddle",
			Responses: []string{
				"${0} cuddled ${1}",
				"${0} snuggled up to ${1}",
			},
			Site: siteNekosLife,
		},
		"hug": {
			Noun: "Hug",
			URL:  "hug",
			Responses: []string{
				"${0} hugged ${1}",
				"${0} wrapped their arms around ${1} and hugged them",
			},
			Site: siteNekosLife,
		},
		"spank": {
			Noun: "Spank",
			URL:  "spank",
			Responses: []string{
				"${0} spanked ${1}",
				"${1} got their ass spanked by ${0}",
				"${0} beat ${1}'s ass",
			},
			Site: siteNekosLife,
		},
		"bj": {
			Noun: "Give some head to",
			URL:  "bj",
			Responses: []string{
				"${0} gave ${1} same top",
				"${0} sucked off ${1}",
			},
			Site: siteNekosLife,
		},
		"anal": {
			Noun: "Give some head to",
			URL:  "anal",
			Responses: []string{
				"${0} tore apart ${1}'s ass",
				"${0} put their dick in ${1}'s ass",
			},
			Site: siteNekosLife,
		},
		"tickle": {
			Noun: "Tickle",
			URL:  "tickle",
			Responses: []string{
				"${0} tickled ${1}",
			},
			Site: siteNekosLife,
		},
		"pat": {
			Noun: "Pat",
			URL:  "pat",
			Responses: []string{
				"${0} patted ${1}",
				"${0} gave ${1} some head pats",
			},
			Site: siteNekosLife,
		},
		"feed": {
			Noun: "Feed",
			URL:  "feed",
			Responses: []string{
				"${0} fed ${1} some food",
				"${0} fed ${1}",
				"${0} shoved some food down ${1}'s throat",
			},
			Site: siteNekosLife,
		},
		"poke": {
			Noun: "Poke",
			URL:  "poke",
			Responses: []string{
				"${0} poked ${1}",
			},
			Site: siteNekosLife,
		},
	}
)

func interactionFunc(c *router.Context) {
	var user *discordgo.User
	var err error

	if len(c.Args) == 0 {
		var chnl *discordgo.Channel
		if chnl, err = c.Ses.State.Channel(c.Msg.ChannelID); err == nil {
			if len(chnl.Recipients) == 1 {
				user = chnl.Recipients[0]
			}
		}

		if user == nil {
			c.ReplyInvalidArg(0, "Invalid user specified.")
			return
		}
	} else if user, err = c.Args.Get(0).AsUser(c.Ses); err != nil {
		c.ReplyInvalidArg(0, "Invalid user specified.")
		return
	}

	inter := interactionsMap[c.Route.Name]

	var imgURL string
	switch inter.Site { //nolint:gocritic // More will come in future
	case siteNekosLife:
		{
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()

			url := format.Formatp(
				"https://nekos.life/api/v2/img/${}",
				inter.URL,
			)

			req.SetRequestURIBytes([]byte(url))
			req.Header.SetMethodBytes([]byte("GET"))

			if fasthttp.Do(req, resp) != nil {
				msgErr := router.NewError(format.Formatp(
					"Unable to get image: ${}",
					err,
				))
				c.ReplyAutoHandle(msgErr)

				return
			}
			bodyStr := resp.Body()
			fasthttp.ReleaseResponse(resp)

			var resMap map[string]json.RawMessage
			if json.Unmarshal(bodyStr, &resMap) != nil {
				msgErr := router.NewError(format.Formatp(
					"Unable to unmarshal response: ${}",
					err,
				))
				c.ReplyAutoHandle(msgErr)

				return
			}

			if json.Unmarshal(resMap["url"], &imgURL) != nil {
				msgErr := router.NewError(format.Formatp(
					"Unable to unmarshal URL from response: ${}",
					err,
				))
				c.ReplyAutoHandle(msgErr)

				return
			}
		}
	}

	response := inter.Responses[rand.Intn(len(inter.Responses))]

	msg := router.NewMessage(inter.Noun)
	msg.Desc(format.Formatp(
		response,
		c.Msg.Author.Mention(),
		user.Mention(),
	))
	msg.Image(imgURL)

	c.ReplyAutoHandle(msg)
}

func LoadInteractions(r *router.Route) {
	r.Group(func(r *router.Route) {
		r.Cat("Interactions")

		for name, inter := range interactionsMap {
			argDesc := format.Formatp(
				"The user to ${}",
				strings.ToLower(inter.Noun),
			)

			cmd := r.On(name, interactionFunc)
			cmd.Desc(format.Formatp("${} a user.", inter.Noun))
			cmd.Arg("user", argDesc, false, router.ArgUser)
		}
	})
}
