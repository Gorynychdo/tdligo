package tdclient

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Arman92/go-tdlib"
	"github.com/Gorynychdo/tdligo.git/internal/model"
	"github.com/pkg/errors"
)

type TDClient struct {
	client   *tdlib.Client
	config   *model.Config
	sender   chan model.OutgoingMessage
	receiver chan tdlib.UpdateMsg
}

func NewTDClient(config *model.Config) *TDClient {
	tdlib.SetLogVerbosityLevel(1)
	client := &TDClient{
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
	}
	client.receiver = client.client.GetRawUpdatesChannel(100)
	return client
}

func (c *TDClient) Start() error {
	if err := c.authorize(); err != nil {
		return err
	}
	log.Println("Telegram client started")
	c.listenAndServe()
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

func (c *TDClient) listenAndServe() {
	for {
		select {
		case update := <-c.receiver:
			if err := c.handleIncoming(update); err != nil {
				log.Println(errors.Wrap(err, "handle incoming message"))
			}
		case message := <-c.sender:
			if err := c.sendMessage(message); err != nil {
				log.Println(errors.Wrap(err, "send message"))
			}
		}
	}
}

func (c *TDClient) handleIncoming(update tdlib.UpdateMsg) error {
	if update.Data["@type"] != "updateNewMessage" {
		return nil
	}

	var data tdlib.UpdateNewMessage
	if err := json.Unmarshal(update.Raw, &data); err != nil {
		return errors.Wrap(err, "unmarshall data")
	}

	message := data.Message
	if message == nil || message.Content == nil || message.IsOutgoing {
		return nil
	}
	content, ok := message.Content.(*tdlib.MessageText)
	if !ok {
		return nil
	}

	mes := &model.IncomingMessage{
		ID:               message.ID,
		UserID:           message.SenderUserID,
		ChatID:           message.ChatID,
		Date:             message.Date,
		EditDate:         message.EditDate,
		ReplyToMessageID: message.ReplyToMessageID,
		Text:             content.Text.Text,
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

	mesRow, err := json.Marshal(mes)
	if err != nil {
		return errors.Wrap(err, "marshall message")
	}
	fmt.Println(string(mesRow))

	return nil
}

func (c *TDClient) sendMessage(message model.OutgoingMessage) error {
	content := &tdlib.InputMessageText{
		Text: tdlib.NewFormattedText(message.Text, nil),
	}
	_, err := c.client.SendMessage(message.ChatID, 0, false, false, nil, content)
	return err
}
