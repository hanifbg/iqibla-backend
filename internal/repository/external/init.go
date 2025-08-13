package external

import (
	"net/http"

	"github.com/hanifbg/landing_backend/config"
)

type ExternalApiWrapper struct {
	WAApi *WAApi
}

type WAApi struct {
	cfg    *config.AppConfig
	client *http.Client
}

func New(cfg *config.AppConfig, client *http.Client) *ExternalApiWrapper {

	ExternalApiWrapper := &ExternalApiWrapper{
		WAApi: &WAApi{
			cfg:    cfg,
			client: client,
		},
	}

	return ExternalApiWrapper
}
