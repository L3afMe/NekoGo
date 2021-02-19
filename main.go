package main

import (
	"L3afMe/NekoGo/commands"
	"L3afMe/NekoGo/config"
	"L3afMe/NekoGo/router"
	"encoding/binary"
	"flag"
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

func main() {
	stdlog.SetOutput(ioutil.Discard)
	flag.Parse()

	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logFormatter := logging.NewBackendFormatter(logBackend, logFormat)
	logging.SetBackend(logFormatter)

	db, err := bbolt.Open("NekoGo.db", 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(format.Formatp("Unable to initialize database: ${}", err))
	}
	log.Info("Initialized Database")

	conf = config.LoadConfig(db, *fToken)
	conf.Save()

	ses, err = discordgo.New(conf.Token)
	if err != nil {
		log.Panic(format.Formatp("Unable to connect to Discord: ${}", err))
	}
	log.Info("Initialized DiscordGo")

	root := router.NewRoot(conf)
	root.Before(before)

	commands.LoadConfig(root)
	commands.LoadFun(root)
	commands.LoadInteractions(root)
	commands.LoadMisc(root)
	commands.LoadUtility(root)

	log.Info("Loaded commands")

	ses.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		root.FindAndExecute(ses, m.Message)
	})

	ses.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) {
		log.Info(format.Formatp("NekoGo connected to ${}", ses.State.User.String()))
	})

	err = ses.Open()
	if err != nil {
		log.Fatal(format.Formatp("Unable to connect to Discord: &{}", err))
	}

	defer func() {
		ses.Close()
		log.Info("Disconnected from Discord")
		conf.Save()
		log.Info("Database saved")
	}()

	log.Notice("Successfully initialized NekoGo")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Info("Stopping NekoGo gracefully")
}

func before(ctx *router.Context) (cont bool) {
	if ctx.Msg.Author.ID == ctx.Ses.State.User.ID {
		if err := ctx.Ses.ChannelMessageDelete(ctx.Msg.ChannelID, ctx.Msg.ID); err != nil {
			ctx.Log.Warning(format.Formatp("Unable to delete message: ${}", err))
		}

		ctx.Log.Info(format.Formatp("Running '${}'", ctx.Route.GetFullName()))

		err := ctx.Route.Config.DB.Update(func(tx *bbolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists(config.ToBytes("usages"))
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
			ctx.Log.Error(format.Formatp("Unable to update usages in database. Error: ${}", err))
		}

		cont = true
	}

	return cont
}
