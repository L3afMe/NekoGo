package commands

import (
	"L3afMe/Krul/kdgr"
	"L3afMe/Krul/utils"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirkon/go-format"
)

type PaymentType int

func (b PaymentType) Format() string {
	switch b {
	case 2:
		return "PayPal"
	default:
		return "Unknown"
	}
}

type BillingInformation struct {
	ID             string      `json:"id"`
	Type           PaymentType `json:"type"`
	Invalid        bool        `json:"invalid"`
	Email          string      `json:"email"`
	BillingAddress struct {
		Name       string `json:"name"`
		Line1      string `json:"line_1"`
		Line2      string `json:"line_2"`
		City       string `json:"city"`
		State      string `json:"state"`
		Country    string `json:"country"`
		PostalCode string `json:"postal_code"`
	} `json:"billing_address"`
	Country string `json:"country"`
	Default bool   `json:"default"`
}

func (b *BillingInformation) Format() (output string) {
	output = "Type: `" + b.Type.Format() + "`\n" +
		"Email: `" + b.Email + "`\n" +
		"Country: `" + b.Country + "`\n" +
		"ID: `" + b.ID + "`\n" +
		"Valid: `" + strconv.FormatBool(!b.Invalid) + "`\n" +
		"Billing Address\n" +
		"> " + b.BillingAddress.Name + "\n" +
		"> " + b.BillingAddress.Line1 + "\n"

	if b.BillingAddress.Line2 != "" {
		output += "> " + b.BillingAddress.Line2 + "\n"
	}

	output += "> " + b.BillingAddress.City + "\n"

	if b.BillingAddress.State != "" {
		output += "> " + b.BillingAddress.State + "\n"
	}

	output += "> " + b.BillingAddress.Country + "\n" +
		"> " + b.BillingAddress.PostalCode + "\n"

	return
}

func utlPing(ctx *kdgr.Context) {
	delay := strconv.Itoa(int(ctx.Ses.LastHeartbeatAck.Sub(ctx.Ses.LastHeartbeatSent).Milliseconds()))
	content := format.Formatp("WS Delay: ${}ms", delay)

	start := time.Now()
	body, err := utils.GetDiscord(ctx.Ses.Token, "users/@me", nil)
	end := time.Now()
	delay = strconv.Itoa(int(end.Sub(start).Milliseconds()))
	if err != nil {
		ctx.ReplyAutoHandle(kdgr.NewError(format.
			Formatp("Failed to execute GET: ${}\n\n${}", err, content)))
		return
	}
	content += format.Formatp("\nGET Delay: ${}ms", delay)

	start = time.Now()
	body, err = utils.PostDiscord(ctx.Ses.Token,
		format.Formatp("channels/${}/messages", ctx.Msg.ChannelID),
		[]byte("{ \"content\": \"Loading ping...\" }"), nil)
	end = time.Now()
	delay = strconv.Itoa(int(end.Sub(start).Milliseconds()))
	if err != nil {
		ctx.ReplyAutoHandle(kdgr.NewError(format.
			Formatp("Failed to execute POST: ${}\n\n${}", err, content)))
		return
	}
	content += format.Formatp("\nPOST Delay: ${}ms", delay)

	var loadingMsg discordgo.Message
	err = json.Unmarshal([]byte(body), &loadingMsg)

	if err != nil {
		ctx.ReplyAutoHandle(kdgr.NewError(format.
			Formatp("Unable to unmarshal POST message response.\n\n${}", content)))
		return
	}

	ctx.EditAutoHandle(&loadingMsg, kdgr.NewMessage("Ping").Desc(content))
}

func utlBase64(ctx *kdgr.Context, enc bool) {
	inp := ctx.Args.All(" ")

	if enc {
		enc := base64.StdEncoding.EncodeToString([]byte(inp))
		ctx.ReplyAutoHandle(kdgr.NewMessage("Base64").AddField("Encode", "Input: `"+inp+"`\nOutput: `"+enc+"`", false))
	} else {
		dec, e := base64.StdEncoding.DecodeString(inp)
		if e != nil {
			ctx.ReplyInvalidArg(0, "Invalid Base64 string specified")
			return
		}

		ctx.ReplyAutoHandle(kdgr.
			NewMessage("Base64").
			AddField("Decode", "Input: `"+inp+"`\nOutput: `"+string(dec)+"`", false))
	}
}

func utlTokenInfo(ctx *kdgr.Context, withBilling bool) {
	tkn := ctx.Args.Get(0).AsString()

	body, err := utils.GetDiscord(tkn, "users/@me", nil)
	if err != nil {
		ctx.ReplyAutoHandle(kdgr.
			NewError("Unable to execute GET request").
			AddField("Error", format.Formatp("${}", err), false))
	}

	if strings.Contains(body, "401: Unauthorized") {
		ctx.ReplyInvalidArg(1, "Unauthorized token specified")

		return
	}

	var user discordgo.User
	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		ctx.ReplyAutoHandle(kdgr.
			NewError("Unable to unmarshal response").
			AddField("Response JSON", "```"+body+"```", false))
	}

	msg := kdgr.NewMessage("Token Info").Desc("Information about `" + tkn + "`")

	nitro := ""
	switch user.PremiumType {
	case 2:
		nitro = "Boost"
	case 1:
		nitro = "Classic"
	default:
		nitro = "None"
	}

	content := "" +
		"Tag: `" + user.String() + "`\n" +
		"ID: `" + user.ID + "`\n" +
		"Email: `" + user.Email + "`\n" +
		"Verified: `" + strconv.FormatBool(user.Verified) + "`\n" +
		"MFA Enabled: `" + strconv.FormatBool(user.MFAEnabled) + "`\n" +
		"Locale: `" + user.Locale + "`\n" +
		"Nitro: `" + nitro + "`\n"

	msg.AddField("User Information", content, true)

	if withBilling {
		body, err = utils.GetDiscord(tkn, "users/@me/billing/payment-sources", nil)
		if err != nil {
			ctx.ReplyAutoHandle(msg.
				AddField("Billing Information",
					format.Formatp("Error occurred while fetching billing information.\nError: ${}", err),
					false))

			return
		}

		if strings.Contains(body, "401: Unauthorized") {
			msg.AddField("Billing Information",
				"Error occurred while fetching billing information.\nError: 401 Unauthorized",
				false)
			ctx.ReplyAutoHandle(msg)

			return
		}

		var billing []BillingInformation
		err = json.Unmarshal([]byte(body), &billing)
		if err != nil {
			msg.AddField("Billing Information",
				"Unable to unmarshal response.\n JSON: ```\n"+body+"\n```",
				false)
			ctx.ReplyAutoHandle(msg)

			return
		}

		for i := range billing[0:utils.Min(9, len(billing))] {
			content = billing[i].Format()

			defaultStr := ""
			if billing[i].Default {
				defaultStr = " (Default)"
			}
			msg.AddField(
				format.Formatp("Billing Information${} (${}/${})",
					defaultStr, i+1, len(billing)),
				content, true,
			)
		}
	}

	ctx.ReplyAutoHandle(msg)
}

func utlAvatar(c *kdgr.Context) {
	var user *discordgo.User
	if len(c.Args) == 1 {
		var err error
		user, err = c.Args.Get(0).AsUser(c.Ses)
		if err != nil {
			c.ReplyInvalidArg(0, "Invalid user specified")
			return
		}
	} else {
		user = c.Msg.Author
	}

	c.ReplyAutoHandle(kdgr.NewMessage("Avatar").
		Desc(format.Formatp("${}'s Avatar\n[Download Link](${})", user.Mention(), user.AvatarURL("2048"))).
		Footer(format.Formatp("ID: ${}", user.ID), "").
		Image(user.AvatarURL("2048")))
}

func loadUtilityCommands(r *kdgr.Route) {
	r.On("avatar", utlAvatar).
		Desc("Get the avatar of the mentioned user").
		Alias("av", "pfp").
		Arg("user", "The user to get their avatar", false, kdgr.RouteArgUser)

	r.On("ping", utlPing).
		Desc("Get the client WebSocket, GET, and POST latency").
		Alias("latency")

	rBase64 := r.On("base64", func(ctx *kdgr.Context) { ctx.ReplyInvalidArg(0, "Expected 'encode' or 'decode'") }).
		Alias("b64").
		Desc("Encode and decode Base64 values").
		Example("decode SW5vcmkgaXMgdGhlIGJlc3Qgd2FpZnU=", "encode I agree").
		Arg("encode/decode", "Whether to encode or decode value", true, kdgr.RouteArgString).
		Arg("value...", "The value to encode/decode", true, kdgr.RouteArgString)

	rBase64.On("decode", func(ctx *kdgr.Context) { utlBase64(ctx, false) }).
		Alias("d", "dec").
		Desc("Decode a Base64 string to a Unicode string").
		Example("SW5vcmkgaXMgdGhlIGJlc3Qgd2FpZnU=").
		Arg("value...", "The value to encode", true, kdgr.RouteArgString)

	rBase64.On("encode", func(ctx *kdgr.Context) { utlBase64(ctx, true) }).
		Alias("e", "enc").
		Desc("Encode a Unicode string to a Base64 sring").
		Example("I agree").
		Arg("value...", "The value to encode", true, kdgr.RouteArgString)

	r.On("tokeninfo", func(ctx *kdgr.Context) { utlTokenInfo(ctx, false) }).
		Alias("token").
		Desc("Check a tokens user account. Optionally including billing details.").
		Example("ODAyMTc5MzM1OTg0NTEzMDY0.YArduw.30nmw_xqSuUX6hzRAC_li05Jw3Q",
			"billing ODAyMTc5MzM1OTg0NTEzMDY0.YArduw.30nmw_xqSuUX6hzRAC_li05Jw3Q").
		Arg("billing", "Whether te display billing information", false, kdgr.RouteArgString).
		Arg("token", "The token to check", true, kdgr.RouteArgString).
		On("billing", func(ctx *kdgr.Context) { utlTokenInfo(ctx, true) }).
		Alias("b", "bill").
		Desc("Check a tokens user account including billing details.").
		Example("ODAyMTc5MzM1OTg0NTEzMDY0.YArduw.30nmw_xqSuUX6hzRAC_li05Jw3Q").
		Arg("token", "The token to check", true, kdgr.RouteArgString)
}

func LoadUtility(r *kdgr.Route) {
	r.Cat("Utility")

	loadUtilityCommands(r)
}
