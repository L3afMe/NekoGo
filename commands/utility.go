package commands

import (
	"L3afMe/NekoGo/router"
	"L3afMe/NekoGo/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
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
	// Use repeat formats rather than straight
	// chaining concatenation to improve readability
	// TODO: Find a better way to do this while keeping readability
	output += format.Formatp("Type: `${}`\n", b.Type.Format())
	output += format.Formatp("Email: `${}`\n", b.Email)
	output += format.Formatp("Country: `${}`\n", b.Country)
	output += format.Formatp("ID: `${}`\n", b.ID)
	output += format.Formatp("Valid: `${}`\n", strconv.FormatBool(!b.Invalid))
	output += "Billing Address\n"
	output += format.Formatp("> ${}\n", b.BillingAddress.Name)
	output += format.Formatp("> ${}\n", b.BillingAddress.Line1)

	if b.BillingAddress.Line2 != "" {
		output += format.Formatp("> ${}\n", b.BillingAddress.Line2)
	}

	output += format.Formatp("> ${}\n", b.BillingAddress.City)

	if b.BillingAddress.State != "" {
		output += format.Formatp("> ${}\n", b.BillingAddress.State)
	}

	output += format.Formatp("> ${}\n", b.BillingAddress.Country)
	output += format.Formatp("> ${}\n", b.BillingAddress.PostalCode)

	return
}

func utlPing(c *router.Context) {
	lastPingRec := c.Ses.LastHeartbeatAck
	lastPingSent := c.Ses.LastHeartbeatSent
	delay := lastPingRec.Sub(lastPingSent).Milliseconds()

	content := format.Formatp("WS Delay: ${}ms", delay)

	start := time.Now()
	_, err := utils.GetDiscord(c.Ses.Token, "users/@me", nil)
	end := time.Now()
	delay = end.Sub(start).Milliseconds()

	if err != nil {
		msgErr := router.NewError(format.Formatp(
			"Failed to execute GET: ${}\n\n${}",
			err,
			content,
		))
		c.ReplyAutoHandle(msgErr)

		return
	}

	content += format.Formatp("\nGET Delay: ${}ms", delay)

	url := format.Formatp("channels/${}/messages", c.Msg.ChannelID)
	postJSON := []byte("{ \"content\": \"Loading ping...\" }")

	start = time.Now()
	body, err := utils.PostDiscord(c.Ses.Token, url, postJSON, nil)
	end = time.Now()
	delay = end.Sub(start).Milliseconds()

	if err != nil {
		msgErr := router.NewError(format.Formatp(
			"Failed to execute POST: ${}\n\n${}",
			err,
			content,
		))
		c.ReplyAutoHandle(msgErr)

		return
	}
	content += format.Formatp("\nPOST Delay: ${}ms", delay)

	var loadingMsg discordgo.Message
	if json.Unmarshal([]byte(body), &loadingMsg) != nil {
		msgErr := router.NewError(format.Formatp(
			"Unable to unmarshal POST message response.\n\n${}",
			content,
		))
		c.ReplyAutoHandle(msgErr)

		return
	}

	c.EditAutoHandle(&loadingMsg, router.NewMessage("Ping").Desc(content))
}

func utlBase64(c *router.Context, enc bool) {
	msg := router.NewMessage("Base64")
	inp := c.Args.All(" ")

	if enc {
		enc := base64.StdEncoding.EncodeToString([]byte(inp))

		out := format.Formatp(
			"Input: `${}`\nOutput: `${}`",
			inp,
			enc,
		)
		msg.AddField("Encode", out, false)
	} else {
		dec, err := base64.StdEncoding.DecodeString(inp)
		if err != nil {
			c.ReplyInvalidArg(0, "Invalid Base64 string specified.")

			return
		}

		out := format.Formatp(
			"Input: `${}`\nOutput: `${}`",
			inp,
			string(dec),
		)
		msg.AddField("Decode", out, false)
	}

	c.ReplyAutoHandle(msg)
}

func utlTokenInfo(c *router.Context, withBilling bool) {
	tkn := c.Args.Get(0).AsString()

	body, err := utils.GetDiscord(tkn, "users/@me", nil)
	if err != nil {
		msgErr := router.NewError("Unable to execute GET request.")
		msgErr.AddField("Error", err.Error(), false)

		c.ReplyAutoHandle(msgErr)
	}

	if strings.Contains(body, "401: Unauthorized") {
		c.ReplyInvalidArg(1, "Unauthorized token specified.")

		return
	}

	var user discordgo.User
	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		errStr := format.Formatp(
			"```\n${}\n```",
			body,
		)
		msgErr := router.NewError("Unable to unmarshal response.")
		msgErr.AddField("Response JSON", errStr, false)

		c.ReplyAutoHandle(msgErr)
	}

	msg := router.NewMessage("Token Info")
	msg.Desc("Information about `" + tkn + "`")

	nitro := ""
	switch user.PremiumType {
	case 2:
		nitro = "Boost"
	case 1:
		nitro = "Classic"
	default:
		nitro = "None"
	}

	content := format.Formatp("Tag: `${}`\n", user.String())
	content += format.Formatp("ID: `${}`\n", user.ID)
	content += format.Formatp("Email: `${}`\n", user.Email)
	content += format.Formatp("Verified: `${}`\n", strconv.FormatBool(user.Verified))
	content += format.Formatp("MFA Enabled: `${}`\n", strconv.FormatBool(user.MFAEnabled))
	content += format.Formatp("Locale: `${}`\n", user.Locale)
	content += format.Formatp("Nitro: `${}`\n", nitro)

	msg.AddField("User Information", content, true)

	if withBilling {
		if err = utlTokenInfoBilling(tkn, msg); err != nil {
			msg.AddField("Billing Information", err.Error(), false)
		}
	}

	c.ReplyAutoHandle(msg)
}

func utlTokenInfoBilling(tkn string, msg *router.Message) error {
	body, err := utils.GetDiscord(tkn, "users/@me/billing/payment-sources", nil)
	if err != nil {
		errStr := format.Formatp(
			"Error occurred while fetching billing information.\nError: ${}",
			err,
		)

		return errors.New(errStr)
	}

	if strings.Contains(body, "401: Unauthorized") {
		errStr := "Error occurred while fetching billing information.\n" +
			"Error: 401 Unauthorized"

		//nolint:stylecheck // Used in message so allow capitalization
		return errors.New(errStr)
	}

	var billing []BillingInformation
	if json.Unmarshal([]byte(body), &billing) != nil {
		errStr := format.Formatp(
			"Unable to unmarshal person.\nJSON: ```\n${}\n```",
			body,
		)

		return errors.New(errStr)
	}

	for i := range billing[0:utils.Min(9, len(billing))] {
		content := billing[i].Format()

		defaultStr := ""
		if billing[i].Default {
			defaultStr = " (Default)"
		}

		title := format.Formatp(
			"Billing Information${} (${}/${})",
			defaultStr,
			i+1,
			len(billing),
		)

		msg.AddField(title, content, true)
	}

	return nil
}

func utlAvatar(c *router.Context) {
	user := c.Msg.Author
	if len(c.Args) == 1 {
		var err error

		if user, err = c.Args.Get(0).AsUser(c.Ses); err != nil {
			c.ReplyInvalidArg(0, "Invalid user specified.")

			return
		}
	}

	msg := router.NewMessage("Avatar")
	msg.Desc(format.Formatp(
		"${}'s Avatar\n[Download Link](${})",
		user.Mention(),
		user.AvatarURL("2048"),
	))
	msg.Footer(format.Formatp("ID: ${}", user.ID), "")
	msg.Image(user.AvatarURL("2048"))

	c.ReplyAutoHandle(msg)
}

func utlUserInfo(c *router.Context) {
	user := c.Msg.Author
	if len(c.Args) == 1 {
		var err error

		if user, err = c.Args.Get(0).AsUser(c.Ses); err != nil {
			c.ReplyInvalidArg(0, "Invalid user specified.")

			return
		}
	}

	utlUserInfoShow(c, user)
}

func utlUserInfoShow(c *router.Context, user *discordgo.User) {
	c.Log.Info(user.String())
}

func utlChannelInfo(c *router.Context) {
	var chnl *discordgo.Channel
	var err error

	if len(c.Args) == 1 {
		if chnl, err = c.Args.Get(0).AsChannel(c.Ses); err != nil {
			c.ReplyInvalidArg(0, "Invalid channel specified.")
			return
		}
	} else {
		if chnl, err = c.Ses.State.Channel(c.Msg.ChannelID); err != nil {
			msgErr := router.NewError("Unable to get current channel.")
			c.ReplyAutoHandle(msgErr)

			return
		}
	}

	utlChannelInfoShow(c, chnl)
}

func utlChannelInfoShow(c *router.Context, chnl *discordgo.Channel) {
	c.Log.Info(chnl.Name)
}

func utlGuildInfo(c *router.Context) {
	g, err := c.Ses.State.Guild(c.Msg.GuildID)
	if err != nil {
		msgErr := router.NewError("Unable to get current guild.")
		c.ReplyAutoHandle(msgErr)

		return
	}

	utlGuildInfoShow(c, g)
}

func utlGuildInfoShow(c *router.Context, g *discordgo.Guild) {
	c.Log.Info(g.Name)
}

func utlRoleInfo(c *router.Context) {
	g, err := c.Ses.State.Guild(c.Msg.GuildID)
	if err != nil {
		msgErr := router.NewError("Unable to get current guild.")
		c.ReplyAutoHandle(msgErr)

		return
	}

	if arg := c.Args.Get(0); arg.IsRole() {
		if role, err := arg.AsRole(g); err == nil {
			utlRoleInfoShow(c, role)
			return
		}
	}

	roles := utlRoleInfoFindRole(g, c.Args.All(" "))
	switch len(roles) {
	case 0:
		msgErr := router.NewError(format.Formatp(
			"No roles match `${}`",
			c.Args.All(" "),
		))
		c.ReplyAutoHandle(msgErr)

	case 1:
		utlRoleInfoShow(c, roles[0])

	default:
		{
			roleList := make([]string, len(roles))
			for _, role := range roles {
				roleList = append(roleList, format.Formatp(
					"${} - ${}",
					role.Mention(),
					role.ID,
				))
			}

			msg := router.NewMessage("Role Info")
			msg.Desc(format.Formatp(
				"Too many roles match `${}`.",
				c.Args.All(" "),
			))
			msg.AddField("IDs", strings.Join(roleList, "\n"), false)

			c.ReplyAutoHandle(msg)
		}
	}
}

func utlRoleInfoShow(c *router.Context, role *discordgo.Role) {
	c.Log.Info(role.Name)
}

func utlRoleInfoFindRole(g *discordgo.Guild, roleName string) []*discordgo.Role {
	for x := 0; x < 3; x++ {
		matchRoles := []*discordgo.Role{}

		for _, role := range g.Roles {
			if (x == 0 && role.Name == roleName) ||
				(x == 1 && strings.EqualFold(role.Name, roleName)) ||
				(x == 2 && utils.SContainsCI(role.Name, roleName)) {
				matchRoles = append(matchRoles, role)
			}
		}

		if len(matchRoles) != 0 {
			return matchRoles
		}
	}

	return make([]*discordgo.Role, 0)
}

func utlSmartInfo(c *router.Context) {
	if len(c.Args) == 1 {
		arg := c.Args.Get(0)

		if arg.IsUser() {
			if user, err := arg.AsUser(c.Ses); err == nil {
				utlUserInfoShow(c, user)

				return
			}
		}
		if arg.IsChannel() {
			if user, err := arg.AsChannel(c.Ses); err == nil {
				utlChannelInfoShow(c, user)

				return
			}
		}
		if c.Msg.GuildID != "" {
			if g, err := c.Ses.State.Guild(c.Msg.GuildID); err == nil {
				roles := utlRoleInfoFindRole(g, arg.AsString())

				switch len(roles) {
				case 0:
					break
				case 1:
					utlRoleInfoShow(c, roles[0])

					return
				default:
					msgErr := router.NewError(format.Formatp(
						"Too many roles match `${}`, use `roleinfo` instead.",
						arg.AsString(),
					))
					c.ReplyAutoHandle(msgErr)

					return
				}
			}
		}
	}

	if c.Msg.GuildID != "" {
		if g, err := c.Ses.State.Guild(c.Msg.GuildID); err == nil {
			utlGuildInfoShow(c, g)

			return
		}
	}

	if chnl, err := c.Ses.Channel(c.Msg.ChannelID); err == nil {
		if len(chnl.Recipients) == 1 {
			utlUserInfoShow(c, chnl.Recipients[0])
		} else {
			utlChannelInfoShow(c, chnl)
		}

		return
	}

	msgErr := router.NewError(
		"Unable to detect info type.\nIf this happens please " +
			"[open an an issue](https://github.com/L3afMe/NekoGo/issues/new) on GitHub",
	)
	c.ReplyAutoHandle(msgErr)
}

func LoadUtility(r *router.Route) {
	r.Group(func(r *router.Route) {
		r.Cat("Utility")

		cAvatar := r.On("avatar", utlAvatar)
		cAvatar.Desc("Get the avatar of the mentioned user.")
		cAvatar.Alias("av", "pfp")
		cAvatar.Arg("user", "The user to get their avatar", false, router.ArgUser)

		cPing := r.On("ping", utlPing)
		cPing.Desc("Get the client WebSocket, GET, and POST latency.")
		cPing.Alias("latency")

		cBase64 := r.On("base64", func(c *router.Context) {
			c.ReplyInvalidArg(0, "Expected 'encode' or 'decode'")
		})
		cBase64.Alias("b64")
		cBase64.Desc("Encode and decode Base64 values.")
		cBase64.Example("decode SW5vcmkgaXMgdGhlIGJlc3Qgd2FpZnU=", "encode I agree")
		cBase64.Arg("encode/decode", "Whether to encode or decode value", true, router.ArgString)
		cBase64.Arg("value...", "The value to encode/decode", true, router.ArgString)

		cBase64Enc := cBase64.On("decode", func(c *router.Context) {
			utlBase64(c, false)
		})
		cBase64Enc.Alias("d", "dec")
		cBase64Enc.Desc("Decode a Base64 string to a Unicode string.")
		cBase64Enc.Example("SW5vcmkgaXMgdGhlIGJlc3Qgd2FpZnU=")
		cBase64Enc.Arg("value...", "The value to encode", true, router.ArgString)

		cBase64Dec := cBase64.On("encode", func(c *router.Context) {
			utlBase64(c, true)
		})
		cBase64Dec.Alias("e", "enc")
		cBase64Dec.Desc("Encode a Unicode string to a Base64 sring.")
		cBase64Dec.Example("I agree")
		cBase64Dec.Arg("value...", "The value to encode", true, router.ArgString)

		cToken := r.On("tokeninfo", func(c *router.Context) {
			utlTokenInfo(c, false)
		})
		cToken.Alias("token")
		cToken.Desc("Check a tokens user account. Optionally including billing details.")
		cToken.Example(
			"ODAyMTc5MzM1OTg0NTEzMDY0.YArduw.30nmw_xqSuUX6hzRAC_li05Jw3Q",
			"billing ODAyMTc5MzM1OTg0NTEzMDY0.YArduw.30nmw_xqSuUX6hzRAC_li05Jw3Q",
		)
		cToken.Arg("billing", "Whether te display billing information", false, router.ArgString)
		cToken.Arg("token", "The token to check", true, router.ArgString)

		cTokenBilling := cToken.On("billing", func(c *router.Context) {
			utlTokenInfo(c, true)
		})
		cTokenBilling.Alias("b", "bill")
		cTokenBilling.Desc("Check a tokens user account including billing details.")
		cTokenBilling.Example("ODAyMTc5MzM1OTg0NTEzMDY0.YArduw.30nmw_xqSuUX6hzRAC_li05Jw3Q")
		cTokenBilling.Arg("token", "The token to check", true, router.ArgString)

		cInfo := r.On("info", utlSmartInfo)
		cInfo.Desc("Attempts to detect what to view info about.")
		cInfo.Example("l3af@0001", "#General", "Mod", "")
		cInfo.Arg("user/channel/role", "The object to view info about", false, router.ArgString)

		cGuildInfo := r.On("guildinfo", utlGuildInfo)
		cGuildInfo.Desc("Display information about a specific guild.")
		cGuildInfo.In(router.InGuild)
		cGuildInfo.Alias("gi", "guild", "si", "server", "serverinfo")

		cRoleInfo := r.On("roleinfo", utlRoleInfo)
		cRoleInfo.Desc("Display information about a specific role.")
		cRoleInfo.Example("Mod", "")
		cRoleInfo.Alias("ri", "role")
		cRoleInfo.Arg("role...", "The role to view info about", true, router.ArgString)

		cChannelInfo := r.On("channelinfo", utlChannelInfo)
		cChannelInfo.Desc("Display information about a specific channel.")
		cChannelInfo.In(router.InGuild)
		cChannelInfo.Example("#General", "")
		cChannelInfo.Alias("ci", "channel")
		cChannelInfo.Arg("channel", "The channel to view info about", false, router.ArgChannel)

		cUserInfo := r.On("userinfo", utlUserInfo)
		cUserInfo.Desc("Display information about a specic user.")
		cUserInfo.Example("l3af@0001", "")
		cUserInfo.Alias("ui", "user")
		cUserInfo.Arg("user", "The user to view info about", false, router.ArgUser)
	})
}
