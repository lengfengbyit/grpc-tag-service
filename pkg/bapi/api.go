package bapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
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

func (api *API) setTracing(ctx context.Context, url string, req *http.Request) *http.Request {
	// 设置链路追踪
	span, newCtx := opentracing.StartSpanFromContext(
		ctx, "HTTP GET: "+api.URL,
		opentracing.Tag{
			Key:   string(ext.Component),
			Value: "HTTP",
		},
	)
	span.SetTag("url", url)
	_ = opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	return req.WithContext(newCtx)
}

func (api *API) httpDo(ctx context.Context, method, path, token string, param url.Values) ([]byte, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	uri := fmt.Sprintf("%s/%s", api.URL, path)
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		uri,
		strings.NewReader(param.Encode()),
	)
	if err != nil {
		return nil, err
	}

	// 设置链路追踪
	req = api.setTracing(ctx, uri, req)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if token != "" {
		req.Header.Set("token", token)
	}

	log.Printf("http %s path: %s, param: %v", method, path, param)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (api *API) httpGet(ctx context.Context, path string) ([]byte, error) {
	log.Printf("http get path: %s", path)
	resp, err := http.Get(fmt.Sprintf("%s/%s", api.URL, path))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (api *API) httpPost(ctx context.Context, path string, param url.Values) ([]byte, error) {
	body, nil := api.httpDo(ctx, "POST", path, "", param)
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
	token, _ := api.getAccessToken(ctx)
	path := fmt.Sprintf("api/v1/tags?name=%s", name)
	body, err := api.httpDo(ctx, "GET", path, token, nil)
	if err != nil {
		return nil, err
	}

	return body, err
}
