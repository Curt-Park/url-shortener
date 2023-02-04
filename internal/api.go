package internal

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

var (
	configPath   = os.Getenv("CONFIG")
	replicaID    = os.Getenv("REPLICA_ID")
	db           *Database
	node         *snowflake.Node
	enc          *Encoding
	characterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

type Config struct {
	RedisServer string `yaml:"RedisServer"`
	RedisPW     string `yaml:"RedisPW"`
	RedisDB     int    `yaml:"RedisDB"`
}

func init() {
	var err error
	var nodeID int64

	// Get the unique node id.
	if len(replicaID) == 0 {
		nodeID = rand.Int63()
	} else {
		if n, err := strconv.Atoi(replicaID); err != nil {
			log.Panic(err)
		} else {
			nodeID = int64(n)
		}
	}
	nodeID %= 1024

	// Create a unique id generator.
	log.Printf("nodeID: %d", nodeID)
	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Panic(err)
	}

	// Base62 encoder.
	log.Printf("Base62 characterset: %s", characterSet)
	enc = NewEncoding(characterSet)

	if len(configPath) == 0 {
		configPath = "config.yaml"
	}

	// Read the config file.
	log.Printf("config path: %s", configPath)
	yfile, err := os.ReadFile(configPath)
	if err != nil {
		log.Panic(err)
	}
	config := Config{}
	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		log.Panic(err)
	}

	// Init the db client.
	log.Printf("Database - server: %s, db: %d", config.RedisServer, config.RedisDB)
	db = NewDatabase(config.RedisServer, config.RedisPW, config.RedisDB)
}

type ShortenURLReq struct {
	URL string `json:"url" form:"url"`
}

type ShortenURLResp struct {
	Key string `json:"key" form:"key"`
}

//	@Summary		Shorten the URL.
//	@Description	Shorten the URL as 11-length Base62 string.
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ShortenURLReq	true	"url"
//	@Success		200		{object}	error
//	@Failure		400		{object}	error
//	@Router			/shorten [post].
func ShortenURL(c echo.Context) error {
	shortenURLReq := ShortenURLReq{}
	if err := c.Bind(&shortenURLReq); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	// Generate a Base62 unique key.
	key := generateShortURL()

	// Store the value in the database.
	db.Set(key, shortenURLReq.URL)

	return c.JSON(http.StatusOK, ShortenURLResp{Key: key})
}

//	@Summary		Redirect to the original URL.
//	@Description	Redirect to the original URL.
//	@Param			key	path		string	true	"key"
//	@Success		302	{object}	error
//	@Failure		404	{object}	error
//	@Router			/{key} [get].
func OriginalURL(c echo.Context) error {
	key := c.Param("key")
	originalURL, exist := db.Get(key)
	if !exist {
		return c.String(http.StatusNotFound, key)
	}
	return c.Redirect(http.StatusFound, originalURL)
}

// Helpers
func generateShortURL() string {
	// Generate a snowflake ID.
	id := node.Generate()
	id_bytes := []byte{
		byte(0xff & id),
		byte(0xff & (id >> 8)),
		byte(0xff & (id >> 16)),
		byte(0xff & (id >> 24)),
		byte(0xff & (id >> 32)),
		byte(0xff & (id >> 40)),
		byte(0xff & (id >> 48)),
		byte(0xff & (id >> 56)),
	}

	// Convert to Base62.
	key := enc.EncodeToString(id_bytes)

	return key
}
