package controllers

import (
	"encoding/json"
	"fmt"
	"helpdesk_backend/db"
	"helpdesk_backend/logger"
	"helpdesk_backend/models"
	"helpdesk_backend/utilities"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// var (
// 	connections = make(map[string]*websocket.Conn)
// 	// connections    = *websocket.Conn
// 	connectionsMux sync.Mutex
// )

// type Connection struct {
// 	ConnectionID string `json:"connection_id"`
// 	PageID       string `json:"page_id"`
// }

type Connection struct {
	Conn   *websocket.Conn
	PageID string
}

var (
	connections = make(map[string]*Connection)
	mutex       sync.Mutex
)

type User struct {
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
}

func CreateUser(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	email := ctx.PostForm("email")

	user := User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	err := ctx.ShouldBindJSON(&user)

	err = db.InsertMongoDocument("user_profile", user)
	if err != nil {
		logger.ZapLogger.Error("error inserting user", err)
		utilities.Render(ctx, utilities.ERRORS["INTERNAL_ERROR"], utilities.ERRORS["INTERNAL_ERROR"].StatusCode)
		return
	}
	authToken := utilities.GenerateUuid()
	// db.SetRedisKey(db.RedisClient, name, authToken, -1)
	db.Cache.SetCustomCache(name, authToken)

	resp := map[string]interface{}{
		"status": "success",
		"token":  authToken,
		"user":   user,
	}

	// resp := map[string]interface{}{
	// 	"status": "success",
	// }

	utilities.Render(ctx, resp, http.StatusOK)

}

func LoginUser(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")

	var user User

	qb := map[string]interface{}{
		"name":     name,
		"password": password,
	}

	err := db.FindOneMongo(models.Collections["UserProfile"], qb, &user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			resp := map[string]interface{}{
				"status":    "error",
				"error_msg": "user not found",
			}
			utilities.Render(ctx, resp, http.StatusOK)
			return
		}
		logger.ZapLogger.Error("error inserting user", err)
		utilities.Render(ctx, utilities.ERRORS["INTERNAL_ERROR"], utilities.ERRORS["INTERNAL_ERROR"].StatusCode)
		return
	}

	authToken := utilities.GenerateUuid()
	// db.SetRedisKey(db.RedisClient, name, authToken, -1)
	db.Cache.SetCustomCache(name, authToken)

	resp := map[string]interface{}{
		"status": "success",
		"token":  authToken,
	}

	utilities.Render(ctx, resp, http.StatusOK)
}

func GetUser(ctx *gin.Context) {
	user := ctx.GetHeader("x-user")

	authToken := strings.Split(ctx.GetHeader("authorization"), " ")[1]

	// authToken := ctx.GetHeader("x-authToken")
	fetchedToken, err := db.Cache.GetCustomCache(user) // (db.RedisClient, user)
	var resp map[string]interface{}

	if err != nil {
		resp = map[string]interface{}{
			"status": "error",
			// "user":   user,
		}

	}

	if fetchedToken.(string) == authToken {
		resp = map[string]interface{}{
			"status": "success",
			"user":   user,
		}
	}

	utilities.Render(ctx, resp, http.StatusOK)

}

type Message struct {
	Mid  string `json:"mid"`
	Text string `json:"text"`
}

type Messaging struct {
	Message   Message `json:"message"`
	Recipient struct {
		Id string `json:"id"`
	} `json:"recipient"`
	Sender struct {
		Id string `json:"id"`
	} `json:"sender"`
	Timestamp int64 `json:"timestamp"`
}

type Entry struct {
	Id        string      `json:"id"`
	Messaging []Messaging `json:"messaging"`
	Time      int64       `json:"time"`
}

type WebhookData struct {
	Entry  []Entry `json:"entry"`
	Object string  `json:"object"`
}

func Webhook(ctx *gin.Context) {

	fmt.Println(ctx)
	fmt.Println("ctx.Request.Method", ctx.Request.Method)

	if ctx.Request.Method == "POST" {

		var data WebhookData
		// var data map[string]interface{}

		if err := ctx.BindJSON(&data); err != nil {
			fmt.Println("Error decoding webhook request:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error decoding webhook request"})
			return
		}

		// Process the incoming webhook event
		fmt.Println("Received webhook event:", data)
		// pageID := data["page_id"].(string)
		utilities.PrintPrettyJSON(data)

		senderId := data.Entry[0].Messaging[0].Sender.Id
		clientId := data.Entry[0].Messaging[0].Recipient.Id
		message := data.Entry[0].Messaging[0].Message.Text
		timeStamp := data.Entry[0].Messaging[0].Timestamp

		msg := models.Mesaages{
			SenderId:  senderId,
			Clientid:  clientId, // senderId
			Mesaage:   message,
			Timestamp: timeStamp,
			PageId:    clientId,
		}
		utilities.PrintPrettyJSON(msg)

		db.InsertMongoDocument(models.Collections["Messages"], msg)

		// fmt.Println(senderId, message, clientId)

		// Query MongoDB to get the first connection ID associated with the page ID
		conn, _ := GetConnectionByPageID(clientId)
		fmt.Println("connection", conn)

		if conn != nil {
			// conn, ok := connections[connectionID]
			// if !ok {
			// 	log.Println("Error: Connection not found for connection ID:", connectionID)
			// 	return
			// }
			msg := map[string]interface{}{
				"action":     "message",
				"senderId":   senderId,
				"clientId":   clientId,
				"pageId":     clientId,
				"created_at": time.Now(),
				"message":    message,
			}

			utilities.PrintPrettyJSON(msg)

			jsonData, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error marshaling JSON:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Println("Error writing message to WebSocket:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
		}

		ctx.Status(http.StatusOK)
		return
	}

	challenge := ctx.Query("hub.challenge")
	if challenge == "" {
		// If no challenge value is provided, respond with an error
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing challenge value"})
		return
	}

	// Respond with the challenge value to verify the webhook URL
	ctx.String(http.StatusOK, challenge)
}

func HandleWebsocket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.ZapLogger.Error(err)
		return
	}
	defer conn.Close()
	fmt.Println("conn", conn)
	for {
		// Read message from the WebSocket connection
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// Handle read error (connection closed, etc.)
			logger.ZapLogger.Error(err)
			break
		}

		// Process the received message based on its type
		switch messageType {
		case websocket.TextMessage:
			// Handle text message
			message := string(p)
			logger.ZapLogger.Info("Received text message:", message)
			utilities.PrintPrettyJSON(message)
			var msg map[string]interface{}
			if err := json.Unmarshal(p, &msg); err != nil {
				logger.ZapLogger.Error("Error decoding JSON:", err)
				continue // Ignore malformed messages
			}

			// Access the pageId field of the struct
			pageID := msg["pageId"].(string)
			logger.ZapLogger.Info("Received page ID:", pageID)
			// Create a new Connection instance
			connection := &Connection{
				Conn:   conn,
				PageID: pageID,
			}

			// Store the connection in the global map
			mutex.Lock()
			connections[pageID] = connection
			mutex.Unlock()

			// Remove the connection from the map when the function exits
			defer func() {
				mutex.Lock()
				delete(connections, pageID)
				mutex.Unlock()
			}()
			// Add your logic to handle text messages here

		case websocket.BinaryMessage:
			// Handle binary message
			// Note: You may need to implement custom logic for binary messages
			logger.ZapLogger.Info("Received binary message:", p)
			// Add your logic to handle binary messages here

		case websocket.CloseMessage:
			// Handle close message (connection closed)
			logger.ZapLogger.Info("Connection closed by client")
			return // Exit the loop and close the connection

		default:
			// Ignore other message types
			logger.ZapLogger.Info("Ignoring message of type:", messageType)
		}
	}
	for {
		conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
		time.Sleep(time.Second)
	}
}

func GetConnectionByPageID(pageID string) (*websocket.Conn, error) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, connection := range connections {
		if connection.PageID == pageID {
			return connection.Conn, nil
		}
	}

	return nil, fmt.Errorf("connection not found for pageID: %s", pageID)
}
