package bapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	APP_KEY    = "go-blog"
	APP_SECRET = "go-blog"
)

var accessToken AccessToken

type AccessToken struct {
	Token  string `json:"token"`
	Expire time.Time
}

type API struct {
	URL string
}

func NewAPI(url string) *API {
	return &API{URL: url}
}

func (api *API) httpGet(ctx context.Context, path string) ([]byte, error) {
	log.Printf("http get url: %s", path)
	resp, err := http.Get(fmt.Sprintf("%s/%s", api.URL, path))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (api *API) httpPost(ctx context.Context, path string, param url.Values) ([]byte, error) {

	log.Printf("http post url: %s, param: %v", path, param)
	resp, err := http.PostForm(
		fmt.Sprintf("%s/%s", api.URL, path),
		param,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (api *API) getAccessToken(ctx context.Context) (string, error) {
	if accessToken.Token != "" && accessToken.Expire.Sub(time.Now()) > 0 {
		return accessToken.Token, nil
	}
	params := url.Values{
		"appKey":    {APP_KEY},
		"appSecret": {APP_SECRET},
	}
	body, err := api.httpPost(ctx, "auth", params)
	if err != nil {
		return "", err
	}

	_ = json.Unmarshal(body, &accessToken)
	accessToken.Expire = time.Now().Add(time.Duration(7000 * time.Second))
	return accessToken.Token, nil
}

func (api *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
	token, err := api.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	body, err := api.httpGet(ctx, fmt.Sprintf(
		"%s?token=%s&name=%s",
		"api/v1/tags",
		token, name,
	))
	if err != nil {
		return nil, nil
	}

	return body, err
}
