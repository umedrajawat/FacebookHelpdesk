package utilities

import (
	"bytes"
	"fmt"
	"helpdesk_backend/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func Render(c *gin.Context, response interface{}, responseCode int) {
	switch c.Request.Header.Get("Accept") {
	// case "application/json":
	default:
		// Respond with JSON
		fmt.Println("responseCode", responseCode)
		c.JSON(responseCode, response)
	}
	// c.JSON(responseCode, response)
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func GetEODTime(t time.Time, loc *time.Location, h int, m int, s int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, t.Nanosecond(), loc)
}

func GenerateSessionKey() string {

	return RandStringBytesMaskImprSrc(32, LetterBytes+NumberBytes)
}

func GetPayloadData(c *gin.Context) interface{} {
	var f map[string]interface{}
	if err := c.ShouldBindQuery(&f); err != nil {
		logger.ZapLogger.Errorf("[error] | Error binding get playlod with error -> %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	fmt.Println(f, c.Request.URL.Query())
	return f
}

func PostPayloadData(c *gin.Context) interface{} {
	var f map[string]interface{}
	if c.Request.Header.Get("Content-type") == "application/json" {
		if err := c.ShouldBindJSON(&f); err != nil {
			logger.ZapLogger.Errorf("[error] | Error binding json post playlod with error -> %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else if c.Request.Header.Get("Content-type") == "application/x-www-form-urlencoded" {
		if err := c.Bind(&f); err != nil {
			logger.ZapLogger.Errorf("[error] | Error binding form post playlod with error -> %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		if err := c.ShouldBind(&f); err != nil {
			logger.ZapLogger.Errorf("[error] | Error binding post playlod with error -> %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
	return f
}

func EscapeHtmlEncoding(b []byte) []byte {
	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	return b
}

func GetCurrentLocalTime() time.Time {
	ct := time.Now()
	loc := time.FixedZone("UTC", 0)
	uct := time.Date(ct.Year(), ct.Month(), ct.Day(), ct.Hour(), ct.Minute(), ct.Second(), ct.Nanosecond(), loc) //.Add(time.Hour*time.Duration(5) + time.Minute*time.Duration(30))
	return uct
}

func DownSample() {

}
