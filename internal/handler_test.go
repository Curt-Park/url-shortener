package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type DatabaseMock struct {
	MockDB map[string]interface{}
}

func (db *DatabaseMock) Set(key string, value interface{}) {
	db.MockDB[key] = value
}

func (db *DatabaseMock) Get(key string) (string, bool) {
	var value interface{}
	var found bool
	if value, found = db.MockDB[key]; !found {
		return "", false
	}
	return value.(string), true
}

func (db *DatabaseMock) Delete(key string) {
	delete(db.MockDB, key)
}

func TestShortenURL(t *testing.T) {
	// Setup.
	e := echo.New()
	r, _ := regexp.Compile("[a-zA-Z0-9]{11}")
	userJSON := `{"url":"https://www.longlonglong-url.com/"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &DBHandler{
		shortURLDB: &DatabaseMock{make(map[string]interface{})},
		longURLDB:  &DatabaseMock{make(map[string]interface{})},
	}

	// Assertions.
	firstResp := ShortenURLResp{}
	if assert.NoError(t, h.ShortenURL(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		err := json.NewDecoder(rec.Body).Decode(&firstResp)
		assert.Equal(t, nil, err)
		assert.True(t, r.MatchString(firstResp.Key))
	}

	// 2nd w/ the existing value in DB.
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, h.ShortenURL(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		resp := ShortenURLResp{}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		assert.Equal(t, nil, err)
		assert.Equal(t, firstResp, resp)
	}
}

func TestOriginalURL(t *testing.T) {
	// Init DB.
	h := &DBHandler{
		shortURLDB: &DatabaseMock{make(map[string]interface{})},
		longURLDB:  &DatabaseMock{make(map[string]interface{})},
	}

	// Init URL info.
	long := "https://www.longlonglong-url.com/"
	short := "M8urCp1G000"
	h.longURLDB.Set(short, long)

	// Setup.
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:key")
	c.SetParamNames("key")
	c.SetParamValues(short)

	// Assertions.
	if assert.NoError(t, h.OriginalURL(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, long, rec.HeaderMap.Get("Location"))
	}

	// 2nd w/ a value not existing.
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/:key")
	c.SetParamNames("key")
	c.SetParamValues(short + "1")

	// Assertions.
	if assert.NoError(t, h.OriginalURL(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}
