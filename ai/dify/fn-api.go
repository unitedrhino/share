package dify

import (
	"context"
	"fmt"
	"net/http"
)

func setAPIAuthorization(dc *DifyClient, req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dc.Key))
	req.Header.Set("Content-Type", "application/json")
}

func SendGetRequestToAPI(dc *DifyClient, api string) (httpCode int, bodyText []byte, err error) {
	return SendGetRequest(false, dc, api)
}

func SendPostRequestToAPI(dc *DifyClient, api string, postBody interface{}) (httpCode int, bodyText []byte, err error) {
	return SendPostRequest(false, dc, api, postBody)
}
func SendPostSseRequestToAPI[rsp any](ctx context.Context, dc *DifyClient, api string, postBody interface{}) (httpCode int, bodyText chan rsp, err error) {
	return SendPostRequestSse[rsp](ctx, false, dc, api, postBody)
}
