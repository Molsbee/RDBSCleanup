package clc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type API struct {
	Config      Config
	BearerToken string
}

type Config struct {
	CLCUsername    string
	CLCPassword    string
	RDBSAppfogUser string
	RDBSAppfogPass string
}

func NewAPI(config Config) (*API, error) {
	login, _ := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: config.CLCUsername,
		Password: config.CLCPassword,
	})
	resp, err := http.Post("https://api.ctl.io/v2/authentication/login", "application/json", bytes.NewReader(login))
	if err != nil {
		return nil, errors.New("unable to login to clc api")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non 200 status code received from api (%d)", resp.StatusCode)
	}

	var res map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("unable to parse response from clc api - (%s)", err)
	}

	return &API{BearerToken: res["bearerToken"].(string)}, nil
}

func (a API) GetAccount(accountID string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", "https://api.ctl.io/v2-experimental/accounts/"+accountID, nil)
	req.Header.Add("Authorization", "Bearer "+a.BearerToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var account map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a API) DeleteRDBSSubscription(accountAlias string, subscriptionID int) error {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("https://api.rdbs.ctl.io/v1/%s/subscriptions/%d", accountAlias, subscriptionID), nil)
	req.Header.Add("Authorization", "Bearer "+a.BearerToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("non successful status code received from api (%d)", resp.StatusCode)
	}
	return nil
}

func (a API) DeleteAppfogSubscriptions(externalID string) error {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("https://api.rdbs.ctl.io/partners/appfog/%s", externalID), nil)
	req.SetBasicAuth(a.Config.RDBSAppfogUser, a.Config.RDBSAppfogPass)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("non successful status code received from api (%d)", resp.StatusCode)
	}
	return nil
}
