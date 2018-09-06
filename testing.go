package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	shogun_base = "https://bugyo.hac.lp1.eshop.nintendo.net/shogun/v1"
)

func main() {
	shopn, err := tls.LoadX509KeyPair("shopn.cert", "shopn.key")
	if err != nil {
		panic(err)
	}

	console, err := tls.LoadX509KeyPair("console.cert", "console.key")
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{shopn, console},
				InsecureSkipVerify: true,
			},
		},
	}

	fmt.Printf("[info] testing a connection to shogun...\n")
	resp, err := client.Get(shogun_base + "/contents/ids?shop_id=4&lang=en&country=US&type=title&title_ids=0100000000010000")
	if err != nil {
		fmt.Printf("[err] there was an error while contacting shogun:\n")
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(body))
}
