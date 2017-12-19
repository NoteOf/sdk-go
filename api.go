package sdk

const NoteOfItEndpoint = "http://api.noteof.it"

type api struct {
	Endpoint string
}

func (a *api) GetEndpoint() string {
	if a.Endpoint != "" {
		return a.Endpoint
	}

	return NoteOfItEndpoint
}

type UnauthenticatedAPI struct {
	api
}

type AuthenticatedAPI struct {
	token string
	UnauthenticatedAPI
}

func NewAuthenticatedApi(token string) *AuthenticatedAPI {
	return &AuthenticatedAPI{
		token: token,
	}
}
