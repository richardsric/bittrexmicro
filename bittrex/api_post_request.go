package bittrex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/richardsric/bittrexmicro/helper"
)

// New returns an instanciated bittrex struct
func New(apiKey, apiSecret string) *Bittrex {
	clientInstance := NewClient(apiKey, apiSecret)
	return &Bittrex{clientInstance}
}

// NewClient return a new Bittrex HTTP client
func NewClient(apiKey, apiSecret string) (c *Client) {
	return &Client{apiKey, apiSecret, &http.Client{}}
}

// NewClientWithCustomHttpConfig returns a new Bittrex HTTP client using the predefined http client
func NewClientWithCustomHttpConfig(apiKey, apiSecret string, httpClient *http.Client) (c *Client) {
	return &Client{apiKey, apiSecret, httpClient}
}

// doTimeoutRequest do a HTTP request with timeout
func (c *Client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.HttpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from Bittrex API")
	}
}

// do prepare and process HTTP request to Bittrex API
func (c *Client) do(method string, ressource string, payload string, authNeeded bool) (response []byte, err error) {
	//connectTimer := time.NewTimer(ReqTimeOut * time.Second)
	connectTimer := time.NewTimer(helper.TimeOut * time.Second)
	var rawurl string
	if strings.HasPrefix(ressource, "http") {
		rawurl = ressource
	} else {
		rawurl = fmt.Sprintf("%s/%s", helper.ApiUrl, ressource)
	}

	req, err := http.NewRequest(helper.ApiMethod, rawurl, strings.NewReader(payload))
	if err != nil {
		return
	}
	if helper.ApiMethod == "POST" || helper.ApiMethod == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Add("Accept", "application/json")

	// Auth
	if authNeeded {
		if len(c.ApiKey) == 0 || len(c.ApiSecret) == 0 {
			err = errors.New("You need to set API Key and API Secret to call this method")
			return
		}
		nonce := time.Now().UnixNano()
		q := req.URL.Query()
		q.Set("apikey", c.ApiKey)
		q.Set("nonce", fmt.Sprintf("%d", nonce))
		req.URL.RawQuery = q.Encode()
		mac := hmac.New(sha512.New, []byte(c.ApiSecret))
		_, err = mac.Write([]byte(req.URL.String()))
		sig := hex.EncodeToString(mac.Sum(nil))
		req.Header.Add("apisign", sig)
	}
	//fmt.Println("the request is: ", req)
	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	//fmt.Println(fmt.Sprintf("reponse %s", response), err)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}
	return response, err
}
