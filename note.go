package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/unicode/norm"
)

type NoteText struct {
	// NoteTextID int64     `json:"-"`
	// NoteID     int64     `json:"-"`
	NoteTextValue string    `json:"note_text"`
	Created       time.Time `json:"created"`
}

// NoteTag is a string that includes methods for canonicalizing the tag as used
// by the server itself
type NoteTag string

// c14nReg matches Combining Diacritical Marks
// see: https://en.wikipedia.org/wiki/Combining_Diacritical_Marks
var c14nReg = regexp.MustCompile("[\u0300-\u036f]")

// C14n provides the standard method by which tags are canonicalized.
//
// The process involves converting the string to NFD normalization and
// removing all Combining Diacritical Marks
//
// There should be little need to call this directly other than perhaps
// comparing strings that have not yet been canonicalized by the server,
// it is provided here as a helper for such comparison but shall not
// be needed before a roundtrip to the server.
func (nt NoteTag) C14n() string {
	s := strings.ToLower(string(nt))
	s = norm.NFD.String(s)
	s = c14nReg.ReplaceAllString(s, "")

	return s
}

type Note struct {
	PublicID string `json:"public_id"`
	// UserID     int       `json:"-"`
	// NoteTextID int64     `json:"-"` // maybe get rid of?
	Archived bool      `json:"archived"`
	Starred  bool      `json:"starred"`
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
	if err != nil {
		return nil, err
	}

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
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%d: %s", resp.StatusCode, text)
	}

	tresp := []*Note{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tresp)
	if err != nil {
		return nil, err
	}

	return tresp, nil
}

func (a *AuthenticatedAPI) GetNote(publicID string) (*Note, error) {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/notes/%s", a.GetEndpoint(), publicID), nil)
	if err != nil {
		return nil, err
	}

	a.authRequest(req)

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, ErrInvalidCredentials
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	} else if resp.StatusCode != http.StatusOK {
		text, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%d: %s", resp.StatusCode, text)
	}

	tresp := &Note{}

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
	if err != nil {
		return nil, err
	}
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
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%d: %s", resp.StatusCode, text)
	}

	tresp := &Note{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tresp)
	if err != nil {
		return nil, err
	}

	return tresp, nil
}

func (a *AuthenticatedAPI) PutUpdateNote(n *Note) (*Note, error) {
	client := &http.Client{}

	j, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/notes/%s", a.GetEndpoint(), n.PublicID), bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
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
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%d: %s", resp.StatusCode, text)
	}

	tresp := &Note{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tresp)
	if err != nil {
		return nil, err
	}

	return tresp, nil
}

func (a *AuthenticatedAPI) DeleteNote(publicID string) (bool, error) {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/notes/%s", a.GetEndpoint(), publicID), nil)
	if err != nil {
		return false, err
	}
	a.authRequest(req)

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return false, ErrInvalidCredentials
	} else if resp.StatusCode == http.StatusNotFound {
		// already deleted
		return false, ErrNotFound
	} else if resp.StatusCode != http.StatusNoContent {
		text, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}

		return false, fmt.Errorf("%d: %s", resp.StatusCode, text)
	}

	return true, nil
}
