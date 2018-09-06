package libshogun

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
					Certificates: []tls.Certificate{shopn},
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

func (c *ShogunClient) GetNsId(tid string) (ns int, err error) {
	resp, err := c.DoRequest("/contents/ids?shop_id=4&lang=en&country=US&type=title&title_ids=" + string(tid))
	if err != nil {
		return 0, err
	}

	fmt.Printf(string(resp) + "\n")
	data := &NsRequestResponse{}
	json.Unmarshal(resp, data)

	return data.id_pair.ns_id, nil
}
