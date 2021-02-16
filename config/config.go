package config

import (
	"bufio"
	"encoding/hex"
	"fmt" //nolint:depguard // Used for logging token input
	"os"
	"strings"

	"github.com/op/go-logging"
	"github.com/valyala/fasthttp"
	"go.etcd.io/bbolt"
)

var log = logging.MustGetLogger("Config")

type Config struct {
	Token  string
	Prefix string

	DB *bbolt.DB
}

func ToBytes(src string) []byte {
	in := []byte(src)
	out := make([]byte, hex.EncodedLen(len(in)))
	hex.Encode(out, in)

	return out
}

func FromBytes(src []byte, fallback string) string {
	out := make([]byte, hex.DecodedLen(len(src)))
	if _, err := hex.Decode(out, src); err != nil {
		log.Error("Unable to decode '" + string(src) + "' to string")
		return fallback
	}

	if len(out) == 0 {
		return fallback
	}

	return string(out)
}

func SafePut(b *bbolt.Bucket, key string, value []byte) {
	if err := b.Put(ToBytes(key), value); err != nil {
		log.Error("Unable to put value with key '"+key+"' in database. Error:", err)
	}
}

func (c *Config) Save() {
	err := c.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(ToBytes("Config"))

		SafePut(b, "Token", ToBytes(c.Token))
		SafePut(b, "Prefix", ToBytes(c.Prefix))

		return nil
	})

	if err != nil {
		log.Error("Error occurred while updating database:", err)
	}
}

func LoadConfig(db *bbolt.DB, fToken string) *Config {
	var Token string
	var Prefix string

	err := db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(ToBytes("Config")); err != nil {
			log.Panic("Unable to create 'Config' bucket in database. Error:", err)
		}

		return nil
	})

	if err != nil {
		log.Panic("Error occurred while initializing database:", err)
	}

	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(ToBytes("Config"))

		bToken := b.Get(ToBytes("Token"))
		if fToken != "" {
			Token = fToken
		} else {
			Token = FromBytes(bToken, "")
		}

		bPrefix := b.Get(ToBytes("Prefix"))
		Prefix = FromBytes(bPrefix, "~")

		return nil
	})

	if err != nil {
		log.Error("Error occurred while loading config from database:", err)
	}

	log.Info("Checking token")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if checkToken(Token) {
			break
		}

		fmt.Print("\nInvalid token specificied.\nInput your Discord token.\n > ")
		scanner.Scan()
		Token = strings.Trim(scanner.Text(), " \n\r")
	}

	log.Info("Loaded config")

	return &Config{
		Token,
		Prefix, db,
	}
}

func checkToken(token string) bool {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.Header.Set("authorization", token)
	req.Header.SetMethodBytes([]byte("GET"))
	req.SetRequestURIBytes([]byte("https://discord.com/api/v8/users/@me"))
	if err := fasthttp.Do(req, resp); err != nil {
		fasthttp.ReleaseResponse(resp)
		return false
	}

	body := string(resp.Body())
	fasthttp.ReleaseResponse(resp)

	return !strings.Contains(body, "401: Unauthorized")
}
