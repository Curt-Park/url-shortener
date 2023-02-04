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
	nodeID       = "0" + os.Getenv("NODE_ID")
	shortURLDB   *Database
	longURLDB    *Database
	node         *snowflake.Node
	enc          *Encoding
	characterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

type Config struct {
	RedisServer     string `yaml:"RedisServer"`
	RedisPW         string `yaml:"RedisPW"`
	RedisShortURLDB int    `yaml:"RedisShortURLDB"`
	RedisLongURLDB  int    `yaml:"RedisLongURLDB"`
}

func init() {
	var err error
	var nid int

	// Get the unique node id.
	if n, err := strconv.Atoi(nodeID); err != nil {
		nid = rand.Intn(1024)
	} else {
		nid = n % 1024
	}

	// Create a unique id generator.
	log.Printf("nodeID: %d", nid)
	node, err = snowflake.NewNode(int64(nid))
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
	log.Printf(
		"Database - server: %s, short url db: %d, long url db: %d",
		config.RedisServer,
		config.RedisShortURLDB,
		config.RedisLongURLDB,
	)
	shortURLDB = NewDatabase(config.RedisServer, config.RedisPW, config.RedisShortURLDB)
	longURLDB = NewDatabase(config.RedisServer, config.RedisPW, config.RedisLongURLDB)
}

type ShortenURLReq struct {
	URL string `json:"url" form:"url"`
}

type ShortenURLResp struct {
	Key string `json:"key" form:"key"`
}

// @Summary     Shorten the URL.
// @Description Shorten the URL as 11-length Base62 string.
// @Accept      json
// @Produce     json
// @Param       request body     ShortenURLReq true "url"
// @Success     200     {object} ShortenURLResp
// @Failure     400     {object} error
// @Router      /shorten [post].
func ShortenURL(c echo.Context) error {
	shortenURLReq := ShortenURLReq{}
	if err := c.Bind(&shortenURLReq); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	longURL := shortenURLReq.URL
	key, exist := shortURLDB.Get(longURL)
	if !exist {
		// Generate a Base62 unique key.
		key = generateShortURL()

		// Store the value in the database.
		shortURLDB.Set(shortenURLReq.URL, key)
		longURLDB.Set(key, shortenURLReq.URL)
	}

	return c.JSON(http.StatusOK, ShortenURLResp{Key: key})
}

// @Summary     Redirect to the original URL.
// @Description Redirect to the original URL.
// @Param       key path     string true "key"
// @Success     302 {object} string
// @Failure     404 {object} error
// @Router      /{key} [get].
func OriginalURL(c echo.Context) error {
	key := c.Param("key")
	longURL, exist := longURLDB.Get(key)
	if !exist {
		return c.String(http.StatusNotFound, key)
	}
	return c.Redirect(http.StatusFound, longURL)
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
