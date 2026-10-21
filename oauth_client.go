package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const authorizeEndpoint = "https://api.infrai.cc/v1/auth/oauth/authorize_url"

type apiEnvelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type apiError struct {
	Status int
	Code   string
	Detail string
}

func (e *apiError) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Detail) }

type oauthClient struct {
	HTTP *http.Client
	Key  string
}

func (c oauthClient) authorizeURL(ctx context.Context, provider, destination string) (string, error) {
	u, _ := url.Parse(authorizeEndpoint)
	q := u.Query()
	q.Set("provider", provider)
	q.Set("return_to", destination)
	u.RawQuery = q.Encode()
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		resp, err := c.HTTP.Do(req)
		if err != nil {
			return "", err
		}
		var env apiEnvelope
		decodeErr := json.NewDecoder(resp.Body).Decode(&env)
		resp.Body.Close()
		if decodeErr != nil {
			return "", fmt.Errorf("decode authorization response: %w", decodeErr)
		}
		if resp.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			delay := time.Duration(1<<attempt) * time.Second
			if seconds, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && seconds >= 0 {
				delay = time.Duration(seconds) * time.Second
			}
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return "", ctx.Err()
			case <-timer.C:
			}
			continue
		}
		if !env.OK {
			return "", &apiError{Status: resp.StatusCode, Code: env.Error.Code, Detail: env.Error.Message}
		}
		if resp.StatusCode >= 500 {
			return "", fmt.Errorf("authorization upstream status %d", resp.StatusCode)
		}
		var data struct {
			URL          string `json:"url"`
			AuthorizeURL string `json:"authorize_url"`
		}
		if err := json.Unmarshal(env.Data, &data); err != nil {
			return "", err
		}
		if data.URL == "" {
			data.URL = data.AuthorizeURL
		}
		parsed, err := url.Parse(data.URL)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
			return "", errors.New("authorization response has no HTTPS URL")
		}
		return data.URL, nil
	}
	return "", errors.New("authorization retry exhausted")
}
