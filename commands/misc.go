package commands

import (
	"L3afMe/Krul/config"
	"L3afMe/Krul/kdgr"
	"encoding/binary"
	"sort"
	"strconv"

	"github.com/sirkon/go-format"
	"go.etcd.io/bbolt"
)

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func mscUsages(c *kdgr.Context) {
	uses := PairList{}
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
				uses = append(uses, Pair{
					route.Name,
					count,
				})
				total += count
			}
		}

		return nil
	})

	if err != nil {
		c.ReplyAutoHandle(kdgr.NewError(format.Formatp("Error getting usages from database.\n${}", err)))
		return
	}

	sort.Sort(uses)

	msg := kdgr.NewMessage("Usages").Desc(format.Formatp("Total: ${}", total))
	for _, use := range uses {
		msg.AddField(use.Key, strconv.Itoa(use.Value), true)
	}

	c.ReplyAutoHandle(msg)
}

func loadMiscCommands(r *kdgr.Route) {
	r.On("usages", mscUsages).
		Desc("Display how often commands are used")

	r.On("help", kdgr.SendHelp).
		Alias("?").
		Desc("Displays the help menu.\n**Keys**\n`<>` - Required\n`[]` - Optional\n`...` - More than one allowed").
		Arg("categery/command...", "Category or command to diplay help about", false, kdgr.RouteArgString).
		Example("tokeninfo billing", "utility")
}

func LoadMisc(r *kdgr.Route) {
	r.Group(func(r *kdgr.Route) {
		r.Cat("Miscellaneous")

		loadMiscCommands(r)
	})
}
