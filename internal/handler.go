package internal

import (
	"errors"
)

func shortenURL(
	longURL string,
	shortURLDB DatabaseOperations,
	longURLDB DatabaseOperations,
) string {
	key, exist := shortURLDB.Get(longURL)
	if !exist {
		// Generate a Base62 unique key.
		key = generateShortURL()

		// Store the value in the database.
		shortURLDB.Set(longURL, key)
		longURLDB.Set(key, longURL)
	}
	return key
}

func originalURL(key string, longURLDB DatabaseOperations) (string, error) {
	longURL, exist := longURLDB.Get(key)
	var err error
	if !exist {
		err = errors.New("key does not exist")
	}
	return longURL, err
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
