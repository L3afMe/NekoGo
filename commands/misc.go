package commands

import (
	"L3afMe/NekoGo/config"
	"L3afMe/NekoGo/router"
	"L3afMe/NekoGo/utils"
	"encoding/binary"
	"sort"
	"strconv"

	"github.com/sirkon/go-format"
	"go.etcd.io/bbolt"
)

type pair struct {
	key   string
	value int
}

type pairList []pair

// Len is needed for sorting
func (p pairList) Len() int { return len(p) }

// Swap is needed for sorting
func (p pairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Less is needed for sorting
func (p pairList) Less(i, j int) bool { return p[i].value < p[j].value }

func mscUsages(c *router.Context) {
	uses := pairList{}
	total := 0

	err := c.Route.Config.DB.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(config.ToBytes("usages"))
		if err != nil {
			return err
		}

		for _, route := range c.Route.Parent.Routes {
			val := b.Get(config.ToBytes(route.Name))

			if val != nil {
				count := int(binary.LittleEndian.Uint32(val))

				uses = append(uses, pair{
					route.Name,
					count,
				})

				total += count
			}
		}

		return nil
	})

	if err != nil {
		msgErr := router.NewError(format.Formatp(
			"Error getting usages from database.\n${}",
			err,
		))
		c.ReplyAutoHandle(msgErr)

		return
	}

	sort.Sort(uses)

	msg := router.NewMessage("Usages")
	msg.Desc(format.Formatp(
		"Total: ${}",
		total,
	))

	for _, use := range uses {
		msg.AddField(use.key, strconv.Itoa(use.value), true)
	}

	c.ReplyAutoHandle(msg)
}

func mscAbout(c *router.Context) {
	msg := router.NewMessage("About")
	msg.Desc("A simple to use Discord selfbot written in Golang with a focus on speed and plethora of commands")
	msg.AddField("Version", utils.Version, true)
	msg.AddField("Developer", "[GitHub](https://github.com/L3afMe) | L3af#0001", true)
	msg.AddField("GitHub", "[L3afMe/NekoGo](https://github.com/L3afMe/NekoGo)", true)
	msg.AddField("Framework", "[DiscordGo](https://github.com/bwmarrin/discordgo)", true)
	msg.Thumbnail("https://user-images.githubusercontent.com/72546287/" +
		"108258608-16cd9580-71c5-11eb-9544-6b6c25951c55.png")

	c.ReplyAutoHandle(msg)
}

func LoadMisc(r *router.Route) {
	r.Group(func(r *router.Route) {
		r.Cat("Miscellaneous")

		cUsages := r.On("usages", mscUsages)
		cUsages.Desc("Display how often commands are used.")

		cAbout := r.On("about", mscAbout)
		cAbout.Desc("Show some info about NekoGo.")

		cHelp := r.On("help", router.SendHelp)
		cHelp.Alias("?", "man")
		cHelp.Desc(
			"Displays the help menu.\n" +
				"**Keys**\n" +
				"`<>` - Required\n" +
				"`[]` - Optional\n" +
				"`...` - More than one allowed",
		)
		cHelp.Arg("categery/command...", "Category or command to diplay help about", false, router.ArgString)
		cHelp.Example("tokeninfo billing", "utility")
	})
}
