package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TokenRequest struct {
	Usage string `json:"usage"`
}

type TokenResponse struct {
	User *User `json:"user"`
	*Token
}

var (
	ErrInvalidCredentials       = errors.New("invalid username or password")
	ErrUnexpectedServerResponse = errors.New("invalid server response")
)

func (u *UnauthenticatedAPI) DoAuth(username, password, usage string) (*TokenResponse, error) {
	outjson, err := json.Marshal(TokenRequest{
		Usage: usage,
	})
	if err != nil {
		return nil, err
	}
	reqbody := bytes.NewBuffer(outjson)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", u.GetEndpoint()+"/auth", reqbody)

	// Headers
	req.SetBasicAuth(username, password)

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, ErrInvalidCredentials
	} else if resp.StatusCode != http.StatusCreated {
		text, err := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d: %s%s", resp.StatusCode, text, err)
	}

	tresp := TokenResponse{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tresp)
	if err != nil {
		return nil, err
	}

	if tresp.APIToken == "" {
		return nil, ErrUnexpectedServerResponse
	}

	return &tresp, nil
}
