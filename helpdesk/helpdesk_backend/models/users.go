package models

/*
UserProfile , models for user profile
*/
type Mesaages struct {
	SenderId  string `bson:"sender_id" json:"sender_id"`
	Clientid  string `bson:"client_id" json:"client_id"`
	Mesaage   string `bson:"mesaage" json:"mesaage"`
	PageId    string `bson:"page_id" json:"page_id"`
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
}
