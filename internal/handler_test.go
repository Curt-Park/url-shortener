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
		assert.Equal(t, err, nil)
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
		assert.Equal(t, err, nil)
		assert.Equal(t, resp, firstResp)
	}
}
