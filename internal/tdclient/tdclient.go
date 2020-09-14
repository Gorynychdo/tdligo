package tdclient

import (
	"fmt"

	"github.com/Arman92/go-tdlib"
	"github.com/Gorynychdo/tdligo.git/internal/model"
	"github.com/pkg/errors"
)

// TDClient tdlib wrap
type TDClient struct {
	client *tdlib.Client
	config *model.TDConfig
}

// NewTDClient constructor
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

// Start TDClient
func (c TDClient) Start() error {
	if err := c.authorize(); err != nil {
		return err
	}
	return c.listenForMessages()
}

// nolint:gocyclo // ...
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
		fmt.Println(update.Data)
		fmt.Print("\n\n")
	}
	return nil
}
