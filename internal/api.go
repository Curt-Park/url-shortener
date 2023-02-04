package internal

import (
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/snowflake"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

var (
	configPath   = os.Getenv("CONFIG")
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

	// Create a unique id generator.
	node, err = snowflake.NewNode(1)
	if err != nil {
		log.Panic(err)
	}

	// Base62 encoder.
	enc = NewEncoding(characterSet)

	if len(configPath) == 0 {
		configPath = "config.yaml"
	}

	// Read the config file.
	yfile, err := os.ReadFile(configPath)
	if err != nil {
		log.Panic(err)
	}
	config := Config{}
	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Database - server: %s, db: %d", config.RedisServer, config.RedisDB)

	// Init the db client.
	db = NewDatabase(config.RedisServer, config.RedisPW, config.RedisDB)
}

type ShortenURLReq struct {
	URL string `json:"url" form:"url"`
}

type ShortenURLResp struct {
	Key string `json:"key" form:"key"`
}

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
