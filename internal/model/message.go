package model

type Message struct {
	ID               int64  `json:"id"`
	SenderUserID     int32  `json:"sender_user_id"`
	ChatID           int64  `json:"chat_id"`
	Date             int32  `json:"date"`
	EditDate         int32  `json:"edit_date"`
	ReplyToMessageID int64  `json:"reply_to_message_id"`
	Text             string `json:"text"`
}
