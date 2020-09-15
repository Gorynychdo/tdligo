package model

type Message struct {
	ID               int64  `json:"id"`
	ChatID           int64  `json:"chat_id"`
	ReplyToMessageID int64  `json:"reply_to_message_id"`
	Date             int32  `json:"date"`
	EditDate         int32  `json:"edit_date"`
	UserID           int32  `json:"user_id"`
	Text             string `json:"text"`
	Username         string `json:"username"`
	UserPhone        string `json:"user_phone"`
	UserFirst        string `json:"user_first"`
	UserLast         string `json:"user_last"`
}
