package kdgr

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirkon/go-format"
)

type messageField struct {
	name   string
	value  string
	inline bool
}

type messageFooter struct {
	text     string
	imageURL string
}

func (m *messageFooter) toFooter() *discordgo.MessageEmbedFooter {
	return &discordgo.MessageEmbedFooter{
		IconURL: m.imageURL,
		Text:    m.text,
	}
}

type messageHeader struct {
	text     string
	imageURL string
	URL      string
}

func (m *messageHeader) toAuthor() *discordgo.MessageEmbedAuthor {
	return &discordgo.MessageEmbedAuthor{
		IconURL: m.imageURL,
		URL:     m.URL,
		Name:    m.text,
	}
}

type Message struct {
	title        string
	description  string
	color        int
	thumbnailURL string
	imageURL     string
	fields       []*messageField
	header       *messageHeader
	footer       *messageFooter
}

func NewMessage(title string) *Message {
	return &Message{title: title, color: 16429549, header: &messageHeader{}, footer: &messageFooter{}}
}

func (msg *Message) ToText() (content string) {
	content = format.Formatp("```\n  [${}]\n${}", msg.title, msg.description)

	for _, field := range msg.fields {
		content += format.Formatp("\n\n  ${}\n${}", field.name, field.value)
	}

	content += format.Formatp("\n```\n${}", msg.imageURL)

	return
}

func (msg *Message) ToEmbed() *discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{}
	for _, field := range msg.fields {
		fields = append(fields, &discordgo.MessageEmbedField{Name: field.name, Value: field.value, Inline: field.inline})
	}

	return &discordgo.MessageEmbed{
		Title:       msg.title,
		Description: msg.description,
		Color:       msg.color,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: msg.thumbnailURL},
		Image:       &discordgo.MessageEmbedImage{URL: msg.imageURL},
		Fields:      fields,
		Author:      msg.header.toAuthor(),
		Footer:      msg.footer.toFooter(),
	}
}

func NewError(description string) *Message {
	return &Message{title: "Error", description: description, color: 16429549, header: &messageHeader{}, footer: &messageFooter{}}
}

func (m *Message) Desc(description string) *Message {
	m.description = description
	return m
}

func (m *Message) Thumbnail(thumbnail string) *Message {
	m.thumbnailURL = thumbnail
	return m
}

func (m *Message) Image(image string) *Message {
	m.imageURL = image
	return m
}

func (m *Message) AddField(name, value string, inline bool) *Message {
	m.fields = append(m.fields, &messageField{name, value, inline})
	return m
}

func (m *Message) Header(name, url, imageURL string) *Message {
	m.header = &messageHeader{name, url, imageURL}
	return m
}

func (m *Message) Footer(name, imageURL string) *Message {
	m.footer = &messageFooter{name, imageURL}
	return m
}

func (c *Context) EditAuto(oldMsg *discordgo.Message, msg *Message) (*discordgo.Message, error) {
	perms, err := c.Ses.State.UserChannelPermissions(c.Ses.State.User.ID, oldMsg.ChannelID)

	embedPerms := int64(0x00004000)
	if err != nil || perms&embedPerms != embedPerms {
		return c.Ses.ChannelMessageEdit(oldMsg.ChannelID, oldMsg.ID, msg.ToText())
	}

	return c.Ses.ChannelMessageEditEmbed(oldMsg.ChannelID, oldMsg.ID, msg.ToEmbed())
}

func (c *Context) EditAutoHandle(oldMsg *discordgo.Message, msg *Message) {
	if _, err := c.EditAuto(oldMsg, msg); err != nil {
		log.Error("Unable to update message. Error:", err)
	}
}

func (c *Context) ReplyAuto(msg *Message) (*discordgo.Message, error) {
	return c.SendAuto(c.Msg.ChannelID, msg)
}

func (c *Context) ReplyAutoHandle(msg *Message) {
	if _, err := c.ReplyAuto(msg); err != nil {
		log.Error("Unable to send message. Error:", err)
	}
}

func (c *Context) SendAuto(chnlID string, msg *Message) (*discordgo.Message, error) {
	perms, err := c.Ses.State.UserChannelPermissions(c.Ses.State.User.ID, chnlID)

	embedPerms := int64(0x00004000)
	if err != nil || perms&embedPerms != embedPerms {
		return c.Ses.ChannelMessageSend(chnlID, msg.ToText())
	}

	return c.Ses.ChannelMessageSendEmbed(chnlID, msg.ToEmbed())
}
