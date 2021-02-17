package main

import (
	"L3afMe/Krul/commands"
	"L3afMe/Krul/config"
	"L3afMe/Krul/kdgr"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	stdlog "log" //nolint:depguard //Needed to disable 3rd party library logging
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/op/go-logging"
	"github.com/sirkon/go-format"
	"go.etcd.io/bbolt"
)

var (
	log       = logging.MustGetLogger("Main")
	logFormat = logging.MustStringFormatter(`[%{time:15:04:05.000}][%{color}%{level:.4s}%{color:reset}][%{module}] %{message}`)
	fToken    = flag.String("t", "", "Set a new token to login with")

	conf *config.Config
	ses  *discordgo.Session
)

func Test() struct{} {
	return struct{}{}
}

func main() {
	x := 42
	fmt.Println(x)

	flag.Parse()
	stdlog.SetOutput(ioutil.Discard)

	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logFormatter := logging.NewBackendFormatter(logBackend, logFormat)
	logging.SetBackend(logFormatter)

	db, err := bbolt.Open("NekoGo.db", 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("Unable to initialize database:", err)
	}
	log.Info("Initialized Database")

	conf = config.LoadConfig(db, *fToken)
	conf.Save()

	ses, err = discordgo.New(conf.Token)
	if err != nil {
		log.Panic("Unable to initialize DiscordGo:", err)
	}
	log.Info("Initialized DiscordGo")

	router := kdgr.New(conf).
		Before(func(ctx *kdgr.Context) bool {
			if ctx.Msg.Author.ID == ctx.Ses.State.User.ID {
				err = ctx.Ses.ChannelMessageDelete(ctx.Msg.ChannelID, ctx.Msg.ID)
				if err != nil {
					log.Warning("Unable to delete message:", err)
				}

				log.Info("Running '" + ctx.Route.GetFullName() + "'")

				err = ctx.Route.Config.DB.Update(func(tx *bbolt.Tx) error {
					var b *bbolt.Bucket
					b, err = tx.CreateBucketIfNotExists(config.ToBytes("usages"))
					if err != nil {
						return err
					}

					keyStr := ctx.Route.GetRootParent().Name
					key := config.ToBytes(keyStr)

					bVal := make([]byte, 8)

					val := b.Get(key)
					if val == nil {
						binary.LittleEndian.PutUint32(bVal, 1)
					} else {
						curVal := binary.LittleEndian.Uint32(val)
						binary.LittleEndian.PutUint32(bVal, curVal+1)
					}

					config.SafePut(b, keyStr, bVal)

					return nil
				})
				if err != nil {
					log.Error("Unable to update usages in database.")
				}

				return true
			}

			return false
		}).
		After(func(ctx *kdgr.Context) {
			log.Info("Finished running '" + ctx.Route.GetFullName() + "'")
		})

	commands.LoadConfig(router)
	commands.LoadFun(router)
	commands.LoadInteractions(router)
	commands.LoadMisc(router)
	commands.LoadUtility(router)

	log.Info("Loaded commands")

	ses.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(ses, conf.Prefix, ses.State.User.ID, m.Message)
	})

	ses.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) {
		log.Info(format.Formatp("NekoGo connected to ${}", ses.State.User.String()))
	})

	err = ses.Open()
	if err != nil {
		log.Fatal(format.Formatp("Unable to connect to Discord: ", err))
	}

	defer func() {
		ses.Close()
		log.Info("Disconnected from Discord")
		conf.Save()
		log.Info("Database saved")
	}()

	log.Notice("Successfully initialized NekoGo")
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Info("Stopping NekoGo gracefully")
}
