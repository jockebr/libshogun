package libshogun

import (
	"net/http"
)

type ShogunClient struct {
	http        *http.Client
	dauth_token string
}

type Title struct {
	id int
	name
	banner_url
	release_date string
	is_new
	is_dlc bool
	description
	genre       string
	size        int64
	screenshots []string
	movies      []string
	publisher   Publisher
}

type Publisher struct {
	id   int
	name string
}
