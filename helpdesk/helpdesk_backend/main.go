// COREAPP GOLANG APIs
//
// .
//
//	BasePath: /
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"helpdesk_backend/logger"
	"helpdesk_backend/server"
	"time"
)

// func init() {
// 	configs.REDIS_PATH = os.Getenv("REDIS_ADDR")
// }

func main() {
	go logger.LogFileChanger()
	time.Sleep(2 * time.Second)

	server.Init()
	// tests.Init()
}
