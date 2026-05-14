package auth

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
)

const (
	ContentType             = "Content-Type"
	ApplicationJson         = "application/json"
	ExchangeRefreshTokenURL = "https://securetoken.googleapis.com/v1/token" //nolint: gosec
)

func (a *auth) exchangeRefreshToken(ctx context.Context, payLoad RefreshTokenRequest) (RefreshTokenResponse, error) {
	var result RefreshTokenResponse
	bodyPayload, err := a.json.Marshal(payLoad)
	if err != nil {
		return result, errors.NewWithCode(codes.CodeHttpMarshal, "%s", err.Error())
	}

	var param = url.Values{}
	param.Set("key", a.conf.Firebase.ApiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ExchangeRefreshTokenURL+"?"+param.Encode(), bytes.NewBuffer(bodyPayload))
	if err != nil {
		return result, errors.NewWithCode(codes.CodeErrorHttpNewRequest, "%s", err.Error())
	}
	req.Header.Set(ContentType, ApplicationJson)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return result, errors.NewWithCode(codes.CodeErrorHttpDo, "%s", err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, errors.NewWithCode(codes.CodeErrorIoutilReadAll, "%s", err.Error())
	}

	if resp.StatusCode != 200 {
		a.log.Error(ctx, fmt.Sprintf("error exchange refresh token, req: %v, resp: %s", payLoad, string(body)))
		return result, errors.NewWithCode(codes.CodeErrorHttpDo, "error exchangeRefreshToken")
	}

	err = a.json.Unmarshal(body, &result)
	if err != nil {
		return result, errors.NewWithCode(codes.CodeUnmarshal, "error exchangeRefreshToken")
	}

	return result, nil
}
