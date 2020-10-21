package bapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/grpc/metadata"

	"golang.org/x/net/context/ctxhttp"
)

const (
	APPKey    = "knight"
	APPSecret = "go_tour"
)

type AccessToken struct {
	Token string `json:"token"`
}

type API struct {
	URL string
}

func NewAPI(url string) *API {
	return &API{URL: url}
}

func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
	resp, err := ctxhttp.Get(ctx, http.DefaultClient, fmt.Sprintf("%s/%s", a.URL, path))
	//resp, err := http.Get(fmt.Sprintf("%s/%s", a.URL, path))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (a *API) getAccessToken(ctx context.Context) (string, error) {
	// 调用第三方需要设置超时时间
	body, err := a.httpGet(ctx, fmt.Sprintf("%s?app_key=%s&app_secret=%s", "auth", APPKey, APPSecret))
	if err != nil {
		return "", err
	}

	var accessToken AccessToken
	_ = json.Unmarshal(body, accessToken)
	return accessToken.Token, nil
}

func (a *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
	// 调用第三方需要设置超时时间
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	md, _ := metadata.FromIncomingContext(ctx)
	log.Printf("md: %+v", md)
	body, err := a.httpGet(ctx, fmt.Sprintf("%s?token=%s&name=%s", "api/v1/tags", token, name))
	if err != nil {
		return nil, err
	}

	return body, nil
}
