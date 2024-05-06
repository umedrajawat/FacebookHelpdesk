package configs

import (
	"net/http"
	"os"
	"time"
)

// var configDev = map[string]string{

const (
	SECRET = "*lkajs!dsandlk#12-21jlii!#/qweoqiw12aksd_l13129ljdlakj1912=93dlakd#klaasdds)1_-1l+1ksjai"

	LOG_PATH = "logger.txt"

	SERVER_PORT = ":8081"

	//REDIS_PATH                       = "localhost" //"localhost"
	REDIS_DB = "0"

	SESSION_CACHE_ALIAS             = "default"
	SESSION_COOKIE_AGE              = 3600 * 12 // expire in 12 hours
	SESSION_EXPIRE_AT_BROWSER_CLOSE = true

	// }

)

var Loc, _ = time.LoadLocation("Asia/Kolkata")

var BlankMap map[string]interface{}

var HttpClient *http.Client

var IOSUserAgents = []string{"webOS", "iPhone", "iPad", "iPod", "Darwin"}
var MobileUserAgents = []string{"Android", "webOS", "iPhone", "iPad", "iPod", "BlackBerry", "Windows Phone", "IEMobile", "Opera Mini", "Mobile", "mobile", "CriOS", "FxiOS"}

var ALLOWED_HOSTS = []string{"127.0.0.1", "0.0.0.0"}

var DummyTTL = time.Duration(84600) * time.Second

var MONGO_PATH = "localhost:27017"
var MONGO_USER = "" // umedrajawat
var MONGO_PASS = "" // umed271998
var MONGO_DB = "test"

var REDIS_PATH = os.Getenv("REDIS_ADDR")
var REDIS_PORT = os.Getenv("REDIS_PORT")
var REDIS_PASS = os.Getenv("REDIS_PASS")
