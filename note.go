package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type NoteText struct {
	// NoteTextID int64     `json:"-"`
	// NoteID     int64     `json:"-"`
	NoteTextValue string    `json:"note_text"`
	Created       time.Time `json:"created"`
}

type NoteTag string

type Note struct {
	NoteID int64 `json:"note_id"`
	// UserID     int       `json:"-"`
	// NoteTextID int64     `json:"-"` // maybe get rid of?
	Archived bool      `json:"archived"`
	Created  time.Time `json:"created"`

	Tags []NoteTag `json:"tags"`

	CurrentText *NoteText   `json:"current"`
	TextHistory []*NoteText `json:"history,omitempty"`
}

func (a *AuthenticatedAPI) authRequest(r *http.Request) {
	r.Header.Add("Authorization", "Token "+a.token)
}

func (a *AuthenticatedAPI) GetNotes() ([]*Note, error) {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", a.GetEndpoint()+"/notes", nil)
	a.authRequest(req)

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, ErrInvalidCredentials
	} else if resp.StatusCode != http.StatusOK {
		text, err := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d: %s%s", resp.StatusCode, text, err)
	}

	tresp := []*Note{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tresp)
	if err != nil {
		return nil, err
	}

	return tresp, nil
}

func (a *AuthenticatedAPI) PostNewNote(n *Note) (*Note, error) {
	client := &http.Client{}

	j, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest("POST", a.GetEndpoint()+"/notes", bytes.NewBuffer(j))
	a.authRequest(req)

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

	tresp := &Note{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tresp)
	if err != nil {
		return nil, err
	}

	return tresp, nil
}
