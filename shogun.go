package libshogun

import (
	"crypto/tls"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	"strconv"
)

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

func (c *ShogunClient) DoRequest(endpoint string) (response []byte, err error) {
	req, err := http.NewRequest("GET", "https://bugyo.hac.lp1.eshop.nintendo.net/shogun/v1"+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-DeviceAuthorization", c.dauth_token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return bytes, nil
}

func (c *ShogunClient) GetNsId(tid string) (ns int64, err error) {
	resp, err := c.DoRequest("/contents/ids?shop_id=4&lang=en&country=US&type=title&title_ids=" + tid)
	if err != nil {
		return 0, err
	}

	id, err := jsonparser.GetInt(resp, "id_pairs", "[0]", "id")
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *ShogunClient) GetTitleData(ns_id int64) (title Title) {
	resp, _ := c.DoRequest("/titles/" + strconv.FormatInt(ns_id, 10) + "?shop_id=4&lang=en&country=US")
	var (
		id           int
		name         string
		banner_url   string
		release_date string
		is_new       bool
		is_dlc       bool
		description  bool
		genre        string
		size         int64
		screenshots  []string
		movies       []string
		publisher    Publisher
	)

	values := [][]string{
		[]string{"id"},
		[]string{"formal_name"},
		[]string{"hero_banner_url"},
		[]string{"release_date"},
		[]string{"is_new"},
		[]string{"is_dlc"},
		[]string{"description"},
		[]string{"total_rom_size"},
		[]string{"publisher"},
	}

	jsonparser.EachKey(data, func(index int, value []byte) {
		switch index {

		}
	}, values...)
}
