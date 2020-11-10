package model

type IncomingMessage struct {
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
	File             *File  `json:"file"`
}

type File struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	MimeType string `json:"mime_type"`
	Size     int32  `json:"size"`
}

type OutgoingMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}
