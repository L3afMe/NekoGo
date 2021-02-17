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

func utlPing(c *kdgr.Context) {
	delay := strconv.Itoa(int(c.Ses.LastHeartbeatAck.Sub(c.Ses.LastHeartbeatSent).Milliseconds()))
	content := format.Formatp("WS Delay: ${}ms", delay)

	start := time.Now()
	_, err := utils.GetDiscord(c.Ses.Token, "users/@me", nil)
	end := time.Now()
	delay = strconv.Itoa(int(end.Sub(start).Milliseconds()))
	if err != nil {
		c.ReplyAutoHandle(kdgr.NewError(format.
			Formatp("Failed to execute GET: ${}\n\n${}", err, content)))
		return
	}
	content += format.Formatp("\nGET Delay: ${}ms", delay)

	start = time.Now()
	body, err := utils.PostDiscord(c.Ses.Token,
		format.Formatp("channels/${}/messages", c.Msg.ChannelID),
		[]byte("{ \"content\": \"Loading ping...\" }"), nil)
	end = time.Now()
	delay = strconv.Itoa(int(end.Sub(start).Milliseconds()))
	if err != nil {
		c.ReplyAutoHandle(kdgr.NewError(format.
			Formatp("Failed to execute POST: ${}\n\n${}", err, content)))
		return
	}
	content += format.Formatp("\nPOST Delay: ${}ms", delay)

	var loadingMsg discordgo.Message
	err = json.Unmarshal([]byte(body), &loadingMsg)

	if err != nil {
		c.ReplyAutoHandle(kdgr.NewError(format.
			Formatp("Unable to unmarshal POST message response.\n\n${}", content)))
		return
	}

	c.EditAutoHandle(&loadingMsg, kdgr.NewMessage("Ping").Desc(content))
}

func utlBase64(c *kdgr.Context, enc bool) {
	inp := c.Args.All(" ")

	if enc {
		enc := base64.StdEncoding.EncodeToString([]byte(inp))
		c.ReplyAutoHandle(kdgr.NewMessage("Base64").AddField("Encode", "Input: `"+inp+"`\nOutput: `"+enc+"`", false))
	} else {
		dec, e := base64.StdEncoding.DecodeString(inp)
		if e != nil {
			c.ReplyInvalidArg(0, "Invalid Base64 string specified")
			return
		}

		c.ReplyAutoHandle(kdgr.
			NewMessage("Base64").
			AddField("Decode", "Input: `"+inp+"`\nOutput: `"+string(dec)+"`", false))
	}
}

func utlTokenInfo(c *kdgr.Context, withBilling bool) {
	tkn := c.Args.Get(0).AsString()

	body, err := utils.GetDiscord(tkn, "users/@me", nil)
	if err != nil {
		c.ReplyAutoHandle(kdgr.
			NewError("Unable to execute GET request").
			AddField("Error", format.Formatp("${}", err), false))
	}

	if strings.Contains(body, "401: Unauthorized") {
		c.ReplyInvalidArg(1, "Unauthorized token specified")

		return
	}

	var user discordgo.User
	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		c.ReplyAutoHandle(kdgr.
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
			c.ReplyAutoHandle(msg.
				AddField("Billing Information",
					format.Formatp("Error occurred while fetching billing information.\nError: ${}", err),
					false))

			return
		}

		if strings.Contains(body, "401: Unauthorized") {
			msg.AddField("Billing Information",
				"Error occurred while fetching billing information.\nError: 401 Unauthorized",
				false)
			c.ReplyAutoHandle(msg)

			return
		}

		var billing []BillingInformation
		err = json.Unmarshal([]byte(body), &billing)
		if err != nil {
			msg.AddField("Billing Information",
				"Unable to unmarshal response.\n JSON: ```\n"+body+"\n```",
				false)
			c.ReplyAutoHandle(msg)

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

	c.ReplyAutoHandle(msg)
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

func utlUserInfo(c *kdgr.Context) {
	user := c.Msg.Author
	if len(c.Args) == 0 {
		var err error
		user, err = c.Args.Get(0).AsUser(c.Ses)
		if err != nil {
			c.ReplyInvalidArg(0, "Invalid user specified.")
			return
		}
	}

	utlUserInfoShow(c, user)
}

func utlUserInfoShow(c *kdgr.Context, user *discordgo.User) {
	c.Log.Info(user.String())
}

func utlChannelInfo(c *kdgr.Context) {
	var chnl *discordgo.Channel
	var err error
	if len(c.Args) == 1 {
		if chnl, err = c.Args.Get(0).AsChannel(c.Ses); err != nil {
			c.ReplyInvalidArg(0, "Invalid channel specified.")
			return
		}
	} else {
		if chnl, err = c.Ses.State.Channel(c.Msg.ChannelID); err != nil {
			c.ReplyAutoHandle(kdgr.NewError("Unable to get current channel"))
			return
		}
	}

	utlChannelInfoShow(c, chnl)
}

func utlChannelInfoShow(c *kdgr.Context, chnl *discordgo.Channel) {
	c.Log.Info(chnl.Name)
}

func utlGuildInfo(c *kdgr.Context) {
	g, err := c.Ses.State.Guild(c.Msg.GuildID)
	if err != nil {
		c.ReplyAutoHandle(kdgr.NewError("Unable to get current guild."))
		return
	}

	utlGuildInfoShow(c, g)
}

func utlGuildInfoShow(c *kdgr.Context, g *discordgo.Guild) {
	c.Log.Info(g.Name)
}

func utlRoleInfo(c *kdgr.Context) {
	g, err := c.Ses.State.Guild(c.Msg.GuildID)
	if err != nil {
		c.ReplyAutoHandle(kdgr.NewError("Unable to get current guild."))
		return
	}

	if arg := c.Args.Get(0); arg.IsRole() {
		if role, err := arg.AsRole(g); err == nil {
			utlRoleInfoShow(c, role)
			return
		}
	}

	roles := utlRoleInfoFindRole(c, g, c.Args.All(" "))
	if len(roles) == 0 {
		c.ReplyAutoHandle(kdgr.NewError(format.Formatp("No roles match `${}`", c.Args.All(" "))))
	} else if len(roles) == 1 {
		utlRoleInfoShow(c, roles[0])
	} else {
		roleList := make([]string, len(roles))
		for _, role := range roles {
			roleList = append(roleList, format.Formatp("${} - ${}", role.Mention(), role.ID))
		}
		c.ReplyAutoHandle(kdgr.
			NewMessage("Role Info").
			Desc("Too many roles match.").
			AddField("IDs", strings.Join(roleList, "\n"), false))
	}
}

func utlRoleInfoShow(c *kdgr.Context, role *discordgo.Role) {
	c.Log.Info(role.Name)
}

func utlRoleInfoFindRole(c *kdgr.Context, g *discordgo.Guild, roleName string) []*discordgo.Role {
	matchRoles := []*discordgo.Role{}

	for _, role := range g.Roles {
		if role.Name == roleName {
			matchRoles = append(matchRoles, role)
		}
	}

	if len(matchRoles) != 0 {
		return matchRoles
	}
	matchRoles = []*discordgo.Role{}

	for _, role := range g.Roles {
		if strings.EqualFold(role.Name, roleName) {
			matchRoles = append(matchRoles, role)
		}
	}

	if len(matchRoles) != 0 {
		return matchRoles
	}
	matchRoles = []*discordgo.Role{}

	for _, role := range g.Roles {
		if utils.SContainsCI(role.Name, roleName) {
			matchRoles = append(matchRoles, role)
		}
	}

	return matchRoles
}

func utlSmartInfo(c *kdgr.Context) {
	if len(c.Args) == 1 {
		arg := c.Args.Get(0)

		if arg.IsUser() {
			user, err := arg.AsUser(c.Ses)
			if err == nil {
				utlUserInfoShow(c, user)
				return
			}
		}
		if arg.IsChannel() {
			user, err := arg.AsChannel(c.Ses)
			if err == nil {
				utlChannelInfoShow(c, user)
				return
			}
		}
		if c.Msg.GuildID != "" {
			g, err := c.Ses.State.Guild(c.Msg.GuildID)
			if err == nil {
				roles := utlRoleInfoFindRole(c, g, arg.AsString())
				switch len(roles) {
				case 0:
					break
				case 1:
					utlRoleInfoShow(c, roles[0])
					return
				default:
					c.ReplyAutoHandle(kdgr.NewError(
						format.Formatp("Too many roles match `${}`, use `roleinfo` instead", arg.AsString()),
					))
					return
				}
			}
		}
	}

	if c.Msg.GuildID != "" {
		g, err := c.Ses.State.Guild(c.Msg.GuildID)
		if err == nil {
			utlGuildInfoShow(c, g)
			return
		}
	}
	chnl, err := c.Ses.Channel(c.Msg.ChannelID)
	if err == nil {
		if len(chnl.Recipients) == 1 {
			utlUserInfoShow(c, chnl.Recipients[0])
		} else {
			utlChannelInfoShow(c, chnl)
		}
		return
	}

	c.ReplyAutoHandle(kdgr.NewError("Unable to detect info type.\nIf this happens please an an issue on GitHub"))
}

func loadUtilityCommands(r *kdgr.Route) {
	r.On("avatar", utlAvatar).
		Desc("Get the avatar of the mentioned user").
		Alias("av", "pfp").
		Arg("user", "The user to get their avatar", false, kdgr.RouteArgUser)

	r.On("ping", utlPing).
		Desc("Get the client WebSocket, GET, and POST latency").
		Alias("latency")

	rBase64 := r.On("base64", func(c *kdgr.Context) { c.ReplyInvalidArg(0, "Expected 'encode' or 'decode'") }).
		Alias("b64").
		Desc("Encode and decode Base64 values").
		Example("decode SW5vcmkgaXMgdGhlIGJlc3Qgd2FpZnU=", "encode I agree").
		Arg("encode/decode", "Whether to encode or decode value", true, kdgr.RouteArgString).
		Arg("value...", "The value to encode/decode", true, kdgr.RouteArgString)

	rBase64.On("decode", func(c *kdgr.Context) { utlBase64(c, false) }).
		Alias("d", "dec").
		Desc("Decode a Base64 string to a Unicode string").
		Example("SW5vcmkgaXMgdGhlIGJlc3Qgd2FpZnU=").
		Arg("value...", "The value to encode", true, kdgr.RouteArgString)

	rBase64.On("encode", func(c *kdgr.Context) { utlBase64(c, true) }).
		Alias("e", "enc").
		Desc("Encode a Unicode string to a Base64 sring").
		Example("I agree").
		Arg("value...", "The value to encode", true, kdgr.RouteArgString)

	r.On("tokeninfo", func(c *kdgr.Context) { utlTokenInfo(c, false) }).
		Alias("token").
		Desc("Check a tokens user account. Optionally including billing details.").
		Example("ODAyMTc5MzM1OTg0NTEzMDY0.YArduw.30nmw_xqSuUX6hzRAC_li05Jw3Q",
			"billing ODAyMTc5MzM1OTg0NTEzMDY0.YArduw.30nmw_xqSuUX6hzRAC_li05Jw3Q").
		Arg("billing", "Whether te display billing information", false, kdgr.RouteArgString).
		Arg("token", "The token to check", true, kdgr.RouteArgString).
		On("billing", func(c *kdgr.Context) { utlTokenInfo(c, true) }).
		Alias("b", "bill").
		Desc("Check a tokens user account including billing details.").
		Example("ODAyMTc5MzM1OTg0NTEzMDY0.YArduw.30nmw_xqSuUX6hzRAC_li05Jw3Q").
		Arg("token", "The token to check", true, kdgr.RouteArgString)

	r.On("info", utlSmartInfo).
		Desc("Attempts to detect what to view info about.").
		Example("l3af@0001", "").
		Arg("user/channel/role", "The object to view info about", false, kdgr.RouteArgString)

	r.On("guildinfo", utlGuildInfo).
		Desc("Display information about a specific guild.").
		In(kdgr.RouteInGuild).
		Alias("gi", "guild", "si", "server", "serverinfo")

	r.On("roleinfo", utlRoleInfo).
		Desc("Display information about a specific role.").
		Example("Mod", "").
		Alias("ri", "role").
		Arg("role...", "The role to view info about", true, kdgr.RouteArgString)

	r.On("channelinfo", utlChannelInfo).
		Desc("Display information about a specific channel.").
		In(kdgr.RouteInGuild).
		Example("#General", "").
		Alias("ci", "channel").
		Arg("channel", "The channel to view info about", false, kdgr.RouteArgChannel)

	r.On("userinfo", utlUserInfo).
		Desc("Display information about a specic user.").
		Example("l3af@0001", "").
		Alias("ui", "user").
		Arg("user", "The user to view info about", false, kdgr.RouteArgUser)
}

func LoadUtility(r *kdgr.Route) {
	r.Group(func(r *kdgr.Route) {
		r.Cat("Utility")

		loadUtilityCommands(r)
	})
}
