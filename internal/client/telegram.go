package client

import (
	"fmt"

	"github.com/Gorynychdo/tdligo.git/internal/model"

	"github.com/Arman92/go-tdlib"
	"github.com/pkg/errors"
)

// TelegramClient обертка над tdlib
type TelegramClient struct {
	client   *tdlib.Client
	settings *model.TelegramSettings
}

// NewTelegramClient конструктор
func NewTelegramClient() TelegramClient {
	tdlib.SetLogVerbosityLevel(1)
	return TelegramClient{
		client: tdlib.NewClient(tdlib.Config{
			APIID:              "api_id",
			APIHash:            "api_hash",
			SystemLanguageCode: "en",
			DeviceModel:        "Server",
			SystemVersion:      "1.0.0",
			ApplicationVersion: "1.0.0",
			DatabaseDirectory:  "./tdlib-db",
			IgnoreFileNames:    false,
		}),
		settings: &model.TelegramSettings{
			Number:   "phone_number",
			Password: "password",
		},
	}
}

// Start запуск клиента на прослушивание событий
func (c TelegramClient) Start() error {
	if err := c.authorize(); err != nil {
		return err
	}
	return c.listenForMessages()
}

// nolint:gocyclo // норм
func (c TelegramClient) authorize() error {
	for {
		state, err := c.client.Authorize()
		if err != nil {
			return err
		}

		switch state.GetAuthorizationStateEnum() {
		case tdlib.AuthorizationStateWaitPhoneNumberType:
			if _, err = c.client.SendPhoneNumber(c.settings.Number); err != nil {
				return errors.Wrap(err, "sending phone number")
			}
		case tdlib.AuthorizationStateWaitCodeType:
			if err = c.setAuthCode(); err != nil {
				return errors.Wrap(err, "setting auth code")
			}
			if _, err = c.client.SendAuthCode(c.settings.Code); err != nil {
				return errors.Wrap(err, "sending auth code")
			}
		case tdlib.AuthorizationStateWaitPasswordType:
			if _, err = c.client.SendAuthPassword(c.settings.Password); err != nil {
				return errors.Wrap(err, "sending auth password")
			}
		case tdlib.AuthorizationStateReadyType:
			fmt.Println("Authorization Ready! Let's rock")
			return nil
		}
	}
}

func (c TelegramClient) setAuthCode() error {
	if c.settings.Code != "" {
		return nil
	}
	fmt.Print("Enter code: ")
	_, err := fmt.Scanln(&c.settings.Code)
	return err
}

func (c TelegramClient) listenForMessages() error {
	for update := range c.client.GetRawUpdatesChannel(100) {
		fmt.Println(update.Data)
		fmt.Print("\n\n")
	}
	return nil
}
