package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"helpdesk_backend/db"
	"helpdesk_backend/logger"
	"helpdesk_backend/models"
	"helpdesk_backend/utilities"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ClientMessages struct {
	ClientID string            `json:"client_id"`
	SenderID string            `json:"sender_id"`
	PageID   string            `json:"page_id"`
	Messages []models.Mesaages `json:"messages"`
}

func groupMessagesBySenders(allMessages []models.Mesaages) map[string]ClientMessages {
	messagesGroupedBySenders := make(map[string]ClientMessages)

	for _, item := range allMessages {
		// utilities.PrintPrettyJSON(item)
		if _, ok := messagesGroupedBySenders[item.Clientid]; !ok {
			messagesGroupedBySenders[item.Clientid] = ClientMessages{
				SenderID: item.SenderId,
				ClientID: item.Clientid,
				PageID:   item.PageId,
				Messages: []models.Mesaages{},
			}
		}

		// Append the message to the Messages slice inside ClientMessages
		client := messagesGroupedBySenders[item.Clientid]
		// fmt.Println("item:", item.SenderId, item.ClientId, item.Mesaage)

		// fmt.Println("going to append message")

		// utilities.PrintPrettyJSON(client)

		client.Messages = append(messagesGroupedBySenders[item.Clientid].Messages, models.Mesaages{
			Mesaage:   item.Mesaage,
			SenderId:  item.SenderId,
			Clientid:  item.Clientid,
			PageId:    item.Clientid,
			Timestamp: item.Timestamp,
		})

		messagesGroupedBySenders[item.Clientid] = client

	}

	return messagesGroupedBySenders
}

func GetAllMessagesController(c *gin.Context) {

	pageId := c.Query("pageId")
	fmt.Println("pageId: ", pageId)

	var msgs []models.Mesaages
	qb := map[string]interface{}{
		// "client_id": pageId,
		"page_id": pageId,
	}
	utilities.PrintPrettyJSON(qb)
	msgCursor, _ := db.FindManyMongo(models.Collections["Messages"], qb, nil)

	for msgCursor.Next(context.TODO()) {
		var b models.Mesaages
		err := msgCursor.Decode(&b)
		if err != nil {
			logger.ZapLogger.Error(err.Error())
			continue
		}
		msgs = append(msgs, b)
		// utilities.PrintPrettyJSON(msgs)
	}
	logger.ZapLogger.Info("length", len(msgs))
	msgCursor.Close(context.TODO())
	// utilities.PrintPrettyJSON(msgs)

	messagesGroupedBySenders := groupMessagesBySenders(msgs)

	// Convert grouped messages to payload
	payload := make([]ClientMessages, 0, len(messagesGroupedBySenders))
	for _, v := range messagesGroupedBySenders {
		payload = append(payload, v)
	}

	resp := map[string]interface{}{
		"messages": payload,
		"status":   "success",
		"message":  "Message reveived successfully",
	}

	utilities.Render(c, resp, http.StatusOK)
}

type Body struct {
	PageID      string `json:"pageId"`
	ClientID    string `json:"clientId"`
	Message     string `json:"message"`
	AccessToken string `json:"accessToken"`
}

type Response struct {
	Message string `json:"message"`
}

func SendMesaage(c *gin.Context) {
	var body Body
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/messages?access_token=%s", body.PageID, body.AccessToken)
	token := fmt.Sprintf("Bearer %s", body.AccessToken)
	// jsonValue, _ := json.Marshal(payload)
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": token,
	}

	payload := map[string]interface{}{
		"recipient":      map[string]string{"id": body.ClientID},
		"messaging_type": "RESPONSE",
		"message":        map[string]string{"text": body.Message},
	}

	response, statusCode, err := utilities.CallApi(url, "POST", header, payload)
	if err != nil {
		logger.ZapLogger.Error(err, "[error->Elastic] api calling error!")
	}
	resp := map[string]interface{}{}
	err = json.Unmarshal(response, &resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if statusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Facebook API error: %v", statusCode)})
		return
	}

	msg := models.Mesaages{
		SenderId:  body.PageID,
		PageId:    body.PageID,
		Mesaage:   body.Message,
		Clientid:  body.PageID,
		Timestamp: time.Now().Unix(),
	}

	fmt.Println("sending message: ", msg)

	utilities.PrintPrettyJSON(msg)

	err = db.InsertMongoDocument(models.Collections["Messages"], msg)

	if err != nil {
		logger.ZapLogger.Error(err)
	}

	c.JSON(http.StatusOK, Response{Message: "Message sent successfully"})
}
