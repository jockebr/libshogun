// Package libshogun provides various utilities for working with the Nintendo Switch's title metadata server, Shogun
package libshogun

import (
	"crypto/tls"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	"strconv"
)

// NewShogunClient creates a new ShogunClient
func NewShogunClient(shopn_cert, shopn_key, dauth_token string) (client *ShogunClient, err error) {
	shopn, err := tls.LoadX509KeyPair(shopn_cert, shopn_key)
	if err != nil {
		return nil, err
	}

	return &ShogunClient{
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates:       []tls.Certificate{shopn},
					InsecureSkipVerify: true,
				},
			},
		},
		dauth_token,
	}, nil
}

// DoRequest sends a request to the specified URL with the proper authentication
func (c *ShogunClient) DoRequest(url string) (response []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-DeviceAuthorization", c.DauthToken)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return bytes, nil
}

// DoShogunRequest sends a request to the specified Shogun endpoint with the proper authentication
func (c *ShogunClient) DoShogunRequest(endpoint string) (response []byte, err error) {
	req, err := http.NewRequest("GET", "https://bugyo.hac.lp1.eshop.nintendo.net/shogun/v1"+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-DeviceAuthorization", c.DauthToken)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return bytes, nil
}

// GetNsId returns the NS ID for the given Title ID
func (c *ShogunClient) GetNsId(tid string) (ns int64, err error) {
	resp, err := c.DoShogunRequest("/contents/ids?shop_id=4&lang=en&country=US&type=title&title_ids=" + tid)
	if err != nil {
		return 0, err
	}

	id, err := jsonparser.GetInt(resp, "id_pairs", "[0]", "id")
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetTitleData returns a Title instance for the given NS ID
func (c *ShogunClient) GetTitleData(ns_id int64) (title *Title, err error) {
	resp, err := c.DoShogunRequest("/titles/" + strconv.FormatInt(ns_id, 10) + "?shop_id=4&lang=en&country=US")
	if err != nil {
		return &Title{}, err
	}

	id, err := jsonparser.GetInt(resp, "id")
	if err != nil {
		return &Title{}, err
	}

	name, err := jsonparser.GetString(resp, "formal_name")
	if err != nil {
		return &Title{}, err
	}

	banner_url, err := jsonparser.GetString(resp, "hero_banner_url")
	if err != nil {
		return &Title{}, err
	}
	banner_url = "https://bugyo.hac.lp1.eshop.nintendo.net" + banner_url

	release_date, err := jsonparser.GetString(resp, "release_date_on_eshop")
	if err != nil {
		return &Title{}, err
	}

	is_new, err := jsonparser.GetBoolean(resp, "is_new")
	if err != nil {
		return &Title{}, err
	}

	is_dlc, err := jsonparser.GetBoolean(resp, "in_app_purchase")
	if err != nil {
		return &Title{}, err
	}

	description, err := jsonparser.GetString(resp, "description")
	if err != nil {
		return &Title{}, err
	}

	genre, err := jsonparser.GetString(resp, "genre")
	if err != nil {
		return &Title{}, err
	}

	size, err := jsonparser.GetInt(resp, "total_rom_size")
	if err != nil {
		return &Title{}, err
	}

	screenshots := []string{}
	jsonparser.ArrayEach(resp, func(value []byte, value_type jsonparser.ValueType, offset int, err error) {
		// todo: add error checking
		screenshots = append(screenshots, "https://bugyo.hac.lp1.eshop.nintendo.net"+string(value))
	}, "images", "url")

	movies := []*Movie{}
	jsonparser.ArrayEach(resp, func(value []byte, value_type jsonparser.ValueType, offset int, err error) {
		// todo: add error checking
		url, _ := jsonparser.GetString(value, "movie_url")
		thumbnail_url, _ := jsonparser.GetString(value, "thumbnail_url")

		movies = append(movies, &Movie{
			URL:       "https://bugyo.hac.lp1.eshop.nintendo.net" + url,
			Thumbnail: "https://bugyo.hac.lp1.eshop.nintendo.net" + thumbnail_url,
		})
	}, "movies")

	pub_id, err := jsonparser.GetInt(resp, "publisher", "id")
	if err != nil {
		return &Title{}, err
	}

	pub_name, err := jsonparser.GetString(resp, "publisher", "name")
	if err != nil {
		return &Title{}, err
	}

	title_id, err := jsonparser.GetString(resp, "applications", "[0]", "id")
	if err != nil {
		return &Title{}, err
	}

	icon_url, err := jsonparser.GetString(resp, "applications", "[0]", "image_url")
	if err != nil {
		return &Title{}, err
	}
	icon_url = "https://bugyo.hac.lp1.eshop.nintendo.net" + icon_url

	return &Title{
		id,
		name,
		banner_url,
		release_date,
		is_new,
		is_dlc,
		description,
		genre,
		size,
		screenshots,
		movies,
		&Publisher{
			pub_id,
			pub_name,
		},
		title_id,
		icon_url,
	}, nil
}
