package commands

import (
	"L3afMe/Krul/kdgr"
	"encoding/json"
	"math/rand"
	"strings"

	"github.com/sirkon/go-format"
	"github.com/valyala/fasthttp"
)

type interactionSite int

const (
	siteNekosLife interactionSite = iota
)

type Interaction struct {
	Noun      string
	URL       string
	Responses []string
	Site      interactionSite
}

var (
	interactionsMap = map[string]Interaction{
		"kiss": {"Kiss", "kiss", []string{
			"${0} kissed ${1}",
		}, siteNekosLife},
		"slap": {"Slap", "slap", []string{
			"${0} slapped ${1}",
		}, siteNekosLife},
		"cuddle": {"Cuddle", "cuddle", []string{
			"${0} cuddled ${1}",
			"${0} snuggled up to ${1}",
		}, siteNekosLife},
		"hug": {"Hug", "hug", []string{
			"${0} hugged ${1}",
			"${0} wrapped their arms around ${1} and hugged them",
		}, siteNekosLife},
		"spank": {"Spank", "spank", []string{
			"${0} spanked ${1}",
			"${1} got their ass spanked by ${0}",
			"${0} beat ${1}'s ass",
		}, siteNekosLife},
		"bj": {"Give some head to", "bj", []string{
			"${0} gave ${1} same top",
			"${0} sucked off ${1}",
		}, siteNekosLife},
		"anal": {"Give some head to", "anal", []string{
			"${0} tore apart ${1}'s ass",
			"${0} put their dick in ${1}'s ass",
		}, siteNekosLife},
		"tickle": {"Tickle", "tickle", []string{
			"${0} tickled ${1}",
		}, siteNekosLife},
		"feed": {"Feed", "feed", []string{
			"${0} fed ${1} some food",
			"${0} fed ${1}",
			"${0} shoved some food down ${1}'s throat",
		}, siteNekosLife},
	}
)

func interaction(c *kdgr.Context) {
	user, err := c.Args.Get(0).AsUser(c.Ses)
	if err != nil {
		c.ReplyInvalidArg(0, "Invalid user specified.")
		return
	}

	inter := interactionsMap[c.Route.Name]
	var url string
	switch inter.Site {
	case siteNekosLife:
		{
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()

			req.SetRequestURIBytes([]byte(format.Formatp("https://nekos.life/api/v2/img/${}", inter.URL)))
			req.Header.SetMethodBytes([]byte("GET"))

			err := fasthttp.Do(req, resp)
			if err != nil {
				c.ReplyAutoHandle(kdgr.NewError(format.Formatp("Unable to get image: ${}", err)))
			}
			bodyStr := resp.Body()

			fasthttp.ReleaseResponse(resp)
			var resMap map[string]json.RawMessage
			err = json.Unmarshal(bodyStr, &resMap)
			if err != nil {
				c.ReplyAutoHandle(kdgr.NewError(format.Formatp("Unable to unmarshal response: ${}", err)))
			}

			err = json.Unmarshal(resMap["url"], &url)
			if err != nil {
				c.ReplyAutoHandle(kdgr.NewError(format.Formatp("Unable to unmarshal response: ${}", err)))
			}
		}
	}

	msg := kdgr.NewMessage(inter.Noun).
		Desc(format.Formatp(inter.Responses[rand.Intn(len(inter.Responses))], c.Msg.Author.Mention(), user.Mention())).
		Image(url)

	c.ReplyAutoHandle(msg)
}

func LoadInteractions(r *kdgr.Route) {
	r.Cat("Interactions")

	for name, inter := range interactionsMap {
		r.On(name, interaction).
			Desc(format.Formatp("${} a user", inter.Noun)).
			Arg("user", format.Formatp("The user to ${}", strings.ToLower(inter.Noun)), true, kdgr.RouteArgUser)
	}
}
