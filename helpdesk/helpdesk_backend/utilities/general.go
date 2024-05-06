package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"helpdesk_backend/logger"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CheckAuthUser(c *gin.Context) string {
	// session := GetApplicationSession(c, false)
	// userUuidI := session.Get("user_uuid")
	// fmt.Println("CHECK AUTH USER UUID: ", userUuidI)
	// userUuid := ""
	// switch o := userUuidI.(type) {
	// case string:
	// 	userUuid = userUuidI.(string)
	// default:
	// 	logger.Logger.Println("[verbose] | Nil user ", o, userUuid)
	// }

	return "abcddd"
}

func Sum(array []float64) float64 {
	result := 0.0
	for _, v := range array {
		result += v
	}
	return result
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func Max(i, j float64) float64 {
	if i > j {
		return i
	}
	return j
}

func GetUUIDv4() string {
	return uuid.NewV4().String()
}

func IsInterfaceString(i interface{}) bool {
	if i == nil {
		return false
	}
	if reflect.TypeOf(i).String() == "string" {
		return true
	}
	return false
}
func IsInterfaceBool(i interface{}) bool {
	if i == nil {
		return false
	}
	if reflect.TypeOf(i).String() == "bool" {
		return true
	}
	return false
}
func IsInterfaceTrue(i interface{}) bool {
	if i == nil {
		return false
	}
	if reflect.TypeOf(i).String() == "string" {
		if i.(string) == "true" {
			return true
		}
		return false
	}
	if reflect.TypeOf(i).String() == "bool" {
		if i.(bool) == true {
			return true
		}
		return false
	}
	return false
}

func StringFromFloat64(f float64) string {
	return fmt.Sprintf("%0.2f", f)
}
func InterfaceToString(i interface{}) (string, error) {
	switch i.(type) {
	case string:
		return i.(string), nil
	case float64, float32:
		return fmt.Sprintf("%0.2f", i), nil
	case int64, int32:
		return fmt.Sprintf("%d", i), nil
	}

	return "", errors.New("cannot interface of type " + reflect.TypeOf(i).String() + " to string")
}

func InterfaceSliceToString(arr []interface{}, joiningSymbol string) string {
	var istring string
	for i := range arr {
		v := arr[i]
		is, err := InterfaceToString(v)
		if err != nil {
			logger.ZapLogger.Error(err)
		}
		istring = istring + joiningSymbol + is
	}
	return istring
}

func MaxOf(vars ...float64) float64 {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}

	return max
}

func MinOf(vars ...float64) float64 {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

func PrintPrettyJSON(i interface{}) {
	data, _ := json.MarshalIndent(i, "", " ")
	fmt.Println(string(data))
}

func Includes(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(v, str) {
			return true
		}
	}

	return false
}

func Zip(lists ...[]interface{}) func() []interface{} {
	zip := make([]interface{}, len(lists))
	i := 0
	return func() []interface{} {
		for j := range lists {
			if i >= len(lists[j]) {
				return nil
			}
			zip[j] = lists[j][i]
		}
		i++
		return zip
	}
}

func DynamicSymParamsGenerator(dynamicStr string) []string {
	// dList := strings.Split(dynamicStr, "_")
	regex := `\(.*?\)`
	matches := RegexFindAll(regex, dynamicStr)
	return matches
}

func RegexFindAll(regex string, str string) []string {
	r := regexp.MustCompile(regex)
	return r.FindAllString(str, -1)
}

func GetMapInterfaceKeys(myMap interface{}) []string {
	var convertedMap map[string]interface{}
	jsonBytes, _ := json.Marshal(myMap)
	json.Unmarshal(jsonBytes, &convertedMap)
	keys := make([]string, len(convertedMap))
	var count int
	for key := range convertedMap {
		keys[count] = key
		count++
		// keys = append(keys, key)
	}

	return keys
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func MapToUrlParams(m map[string]string) string {
	var params string
	for key, value := range m {
		params += key + "=" + value + "&"
	}
	return params
}

func ValidatePhoneNumber(phoneNumber string) bool {
	regex := `^\d{10}$`
	matched, _ := regexp.MatchString(regex, phoneNumber)
	return matched
}

// def last_date_of_month(year,month,weekday=3):
//  lastDayOfMonth = datetime.datetime(year,month,calendar.monthrange(year,month)[1])
//  while lastDayOfMonth.weekday() != weekday:
//      lastDayOfMonth-=datetime.timedelta(days=1)
//  return lastDayOfMonth

func LastDateOfMonth(year int, month int, weekday int) time.Time {
	if weekday == 0 {
		weekday = 3
	}
	lastDayOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
	for lastDayOfMonth.Weekday() != time.Weekday(weekday) {
		lastDayOfMonth = lastDayOfMonth.AddDate(0, 0, -1)
	}
	return lastDayOfMonth
}

func CLTime(cusTime time.Time) time.Time { // CURRENT LOCAL TIME
	ct := cusTime
	loc := time.FixedZone("UTC", 0)
	// loc, _ := time.LoadLocation("Asia/Kolkata")
	// loc = Loc
	uct := time.Date(ct.Year(), ct.Month(), ct.Day(), ct.Hour(), ct.Minute(), ct.Second(), ct.Nanosecond(), loc)
	// .Add(time.Hour*time.Duration(5) + time.Minute*time.Duration(30))
	return uct
}

func MakeDefaultNilValue(object interface{}, objectStruct interface{}) interface{} {
	m := make(map[string]interface{})
	jsonBytes, _ := json.Marshal(object)
	json.Unmarshal(jsonBytes, &m)
	for i := 0; i < reflect.Indirect(reflect.ValueOf(objectStruct)).NumField(); i++ {
		if m[reflect.Indirect(reflect.ValueOf(objectStruct)).Type().Field(i).Tag.Get("json")] == nil {
			switch reflect.Indirect(reflect.ValueOf(objectStruct)).Type().Field(i).Type.Kind() {
			case reflect.Slice:
				m[reflect.Indirect(reflect.ValueOf(objectStruct)).Type().Field(i).Tag.Get("json")] = []string{}
			case reflect.Map:
				m[reflect.Indirect(reflect.ValueOf(objectStruct)).Type().Field(i).Tag.Get("json")] = map[string]interface{}{}
			default:
				m[reflect.Indirect(reflect.ValueOf(objectStruct)).Type().Field(i).Tag.Get("json")] = nil
			}
		}
	}
	return m
}

func AssignDefaultValue(object interface{}, objectStruct interface{}) interface{} {
	m := make(map[string]interface{})
	jsonBytes, _ := json.Marshal(object)
	json.Unmarshal(jsonBytes, &m)
	objectStructValue := reflect.Indirect(reflect.ValueOf(objectStruct))
	for i := 0; i < objectStructValue.NumField(); i++ {
		field := objectStructValue.Type().Field(i)
		jsonTag := field.Tag.Get("json")
		defaultValue := field.Tag.Get("default")

		// Check if the key is not present in the map
		if _, exists := m[jsonTag]; !exists && defaultValue != "" {
			m[jsonTag] = defaultValue
		}
	}
	return m
}

func IsEmptyStruct(s interface{}) bool {
	val := reflect.ValueOf(s)

	// Check if the value is a struct
	if val.Kind() != reflect.Struct {
		return false
	}

	// Iterate through struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Check if the field is zero value
		if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			continue
		}

		// If any field is not zero value, the struct is not considered empty
		return false
	}

	// All fields are zero value, the struct is considered empty
	return true
}

func UnescapeHtmlEncoding(s string) string {
	s = strings.Replace(s, "&lt;", "<", -1)
	s = strings.Replace(s, "&gt;", ">", -1)
	s = strings.Replace(s, "&amp;", "&", -1)
	s = strings.Replace(s, "&quot;", "\"", -1)
	s = strings.Replace(s, "&apos;", "'", -1)
	return s
}

func GenerateCurlString(c *gin.Context) string {
	// Get request method and URL
	method := c.Request.Method
	url := c.Request.URL.String()

	// Initialize the cURL command
	curlCmd := fmt.Sprintf("curl -X %s '%s'", method, url)

	// Add request headers to the cURL command
	for key, values := range c.Request.Header {
		for _, value := range values {
			curlCmd += fmt.Sprintf(" -H '%s: %s'", key, value)
		}
	}

	// Add request body, if present
	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		c.Request.ParseForm()
		body := c.Request.PostForm.Encode()
		curlCmd += fmt.Sprintf(" -d '%s'", body)
	}

	return curlCmd
}
