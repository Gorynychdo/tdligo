package tdclient

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Arman92/go-tdlib"
	"github.com/Gorynychdo/tdligo/internal/model"
	"github.com/pkg/errors"
)

type TDClient struct {
	client *tdlib.Client
	config *model.Config
	files  map[int32]*tdlib.File
}

func NewTDClient(config *model.Config) *TDClient {
	tdlib.SetLogVerbosityLevel(1)
	return &TDClient{
		client: tdlib.NewClient(tdlib.Config{
			APIID:              config.TelegramAPIID,
			APIHash:            config.TelegramAPIHash,
			SystemLanguageCode: "en",
			DeviceModel:        "Server",
			SystemVersion:      "1.0.0",
			ApplicationVersion: "1.0.0",
			DatabaseDirectory:  "./tdlib-db",
			IgnoreFileNames:    false,
		}),
		config: config,
		files:  make(map[int32]*tdlib.File),
	}
}

func (c *TDClient) Start() error {
	if err := c.authorize(); err != nil {
		return err
	}
	log.Println("Telegram client started")
	c.listenIncoming()
	return nil
}

func (c *TDClient) authorize() error {
	for {
		state, err := c.client.Authorize()
		if err != nil {
			return err
		}

		switch state.GetAuthorizationStateEnum() {
		case tdlib.AuthorizationStateWaitPhoneNumberType:
			if _, err = c.client.SendPhoneNumber(c.config.PhoneNumber); err != nil {
				return errors.Wrap(err, "sending phone number")
			}
		case tdlib.AuthorizationStateWaitCodeType:
			if err = c.setAuthCode(); err != nil {
				return errors.Wrap(err, "setting auth code")
			}
			if _, err = c.client.SendAuthCode(c.config.AuthCode); err != nil {
				return errors.Wrap(err, "sending auth code")
			}
		case tdlib.AuthorizationStateWaitPasswordType:
			if _, err = c.client.SendAuthPassword(c.config.AuthPassword); err != nil {
				return errors.Wrap(err, "sending auth password")
			}
		case tdlib.AuthorizationStateReadyType:
			return nil
		}
	}
}

func (c *TDClient) setAuthCode() error {
	if c.config.AuthCode != "" {
		return nil
	}
	fmt.Print("Enter code: ")
	_, err := fmt.Scanln(&c.config.AuthCode)
	return err
}

func (c *TDClient) listenIncoming() {
	for update := range c.client.GetRawUpdatesChannel(100) {
		switch update.Data["@type"] {
		case string(tdlib.UpdateAuthorizationStateType):
			log.Println(string(update.Raw))
		case string(tdlib.UpdateNewMessageType):
			var data tdlib.UpdateNewMessage
			if err := json.Unmarshal(update.Raw, &data); err != nil {
				log.Println(errors.Wrap(err, "unmarshall data"))
				continue
			}
			if err := c.handleIncoming(data.Message); err != nil {
				log.Println(errors.Wrap(err, "handle incoming message"))
			}
		case string(tdlib.UpdateFileType):
			log.Println(string(update.Raw))
		}
	}
}

func (c *TDClient) handleIncoming(message *tdlib.Message) error {
	if message == nil || message.IsOutgoing {
		return nil
	}

	switch message.Content.(type) {
	case *tdlib.MessageText:
		return c.handleNewMessage(message)
	case *tdlib.MessageDocument:
		return c.handleNewFile(message)
	}
	return nil
}

func (c *TDClient) handleNewMessage(message *tdlib.Message) error {
	mes := c.getMessage(message)
	mes.Text = message.Content.(*tdlib.MessageText).Text.Text
	mesRow, err := json.Marshal(mes)
	if err == nil {
		fmt.Println(string(mesRow))
	}
	return errors.Wrap(err, "marshall message")
}

func (c *TDClient) handleNewFile(message *tdlib.Message) error {
	mes := c.getMessage(message)
	content := message.Content.(*tdlib.MessageDocument)
	document := content.Document
	mes.Text = content.Caption.Text
	file, err := c.client.DownloadFile(document.Document.ID, 1)
	if err == nil {
		mes.File = &model.File{
			FileName: document.FileName,
			MimeType: document.MimeType,
			FilePath: file.Local.Path,
			Size:     file.Size,
		}
	} else {
		log.Println(errors.Wrap(err, "downloading file with id"))
	}

	mesRow, err := json.Marshal(mes)
	if err == nil {
		fmt.Println(string(mesRow))
	}
	return errors.Wrap(err, "marshall message")
}

func (c *TDClient) getMessage(message *tdlib.Message) *model.IncomingMessage {
	mes := &model.IncomingMessage{
		ID:               message.ID,
		UserID:           message.SenderUserID,
		ChatID:           message.ChatID,
		Date:             message.Date,
		EditDate:         message.EditDate,
		ReplyToMessageID: message.ReplyToMessageID,
	}

	user, err := c.client.GetUser(message.SenderUserID)
	if err == nil {
		mes.Username = user.Username
		mes.UserPhone = user.PhoneNumber
		mes.UserFirst = user.FirstName
		mes.UserLast = user.LastName
	} else {
		log.Println(errors.Wrap(err, "get user from id"))
	}
	return mes
}

func (c *TDClient) SendMessage(message model.OutgoingMessage) error {
	content := tdlib.NewInputMessageText(tdlib.NewFormattedText(message.Text, nil), false, false)
	_, err := c.client.SendMessage(message.ChatID, 0, false, false, nil, content)
	return err
}
