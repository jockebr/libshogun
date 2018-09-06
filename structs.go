package libshogun

import (
	"net/http"
)

type ShogunClient struct {
	http        *http.Client
	dauth_token string
}

type NsRequestResponse struct {
	id_pair *IdPair `json:"id_pairs"`
}

type IdPair struct {
	ns_id    int `json:"id"`
	title_id int `json:"title_id"`
}
