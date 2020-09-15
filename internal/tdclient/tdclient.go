package tdclient

import (
	"encoding/json"
	"fmt"

	"github.com/Arman92/go-tdlib"
	"github.com/Gorynychdo/tdligo.git/internal/model"
	"github.com/pkg/errors"
)

type TDClient struct {
	client *tdlib.Client
	config *model.TDConfig
}

func NewTDClient(config *model.TDConfig) TDClient {
	tdlib.SetLogVerbosityLevel(1)
	return TDClient{
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
}

func (c TDClient) Start() error {
	if err := c.authorize(); err != nil {
		return err
	}
	return c.listenForMessages()
}

func (c TDClient) authorize() error {
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
			fmt.Println("Authorization Ready! Let's rock")
			return nil
		}
	}
}

func (c TDClient) setAuthCode() error {
	if c.config.AuthCode != "" {
		return nil
	}
	fmt.Print("Enter code: ")
	_, err := fmt.Scanln(&c.config.AuthCode)
	return err
}

func (c TDClient) listenForMessages() error {
	for update := range c.client.GetRawUpdatesChannel(100) {
		if update.Data["@type"] != "updateNewMessage" {
			continue
		}

		var data tdlib.UpdateNewMessage
		if err := json.Unmarshal(update.Raw, &data); err != nil {
			return errors.Wrap(err, "unmarshall data")
		}
		if err := c.handleMessage(data.Message); err != nil {
			return err
		}
	}
	return nil
}

func (c TDClient) handleMessage(message *tdlib.Message) error {
	if message == nil || message.Content == nil || message.IsOutgoing {
		return nil
	}
	content, ok := message.Content.(*tdlib.MessageText)
	if !ok {
		return nil
	}

	user, err := c.client.GetUser(message.SenderUserID)
	if err != nil {
		return errors.Wrap(err, "get user from id")
	}

	mes := &model.Message{
		ID:               message.ID,
		UserID:           message.SenderUserID,
		ChatID:           message.ChatID,
		Date:             message.Date,
		EditDate:         message.EditDate,
		ReplyToMessageID: message.ReplyToMessageID,
		Text:             content.Text.Text,
		Username:         user.Username,
		UserPhone:        user.PhoneNumber,
		UserFirst:        user.FirstName,
		UserLast:         user.LastName,
	}

	mesRow, err := json.Marshal(mes)
	if err != nil {
		return errors.Wrap(err, "marshall message")
	}
	fmt.Println(string(mesRow))

	return nil
}
