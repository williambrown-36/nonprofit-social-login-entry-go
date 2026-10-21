package main

import (
	"errors"
	"net/url"
	"strings"
)

type loginRequest struct {
	Provider string `json:"provider"`
	Workflow string `json:"workflow"`
	Campaign string `json:"campaign"`
}

var destinations = map[string]string{
	"receipt":  "/donor/receipts",
	"reminder": "/volunteer/reminders",
	"report":   "/campaign/reports",
}

func returnTo(base string, input loginRequest) (string, error) {
	if input.Provider != "google" && input.Provider != "github" {
		return "", errors.New("provider must be google or github")
	}
	path, ok := destinations[input.Workflow]
	if !ok || strings.TrimSpace(input.Campaign) == "" {
		return "", errors.New("workflow and campaign are required")
	}
	u, err := url.Parse(base)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", errors.New("APP_BASE_URL must be an HTTPS origin")
	}
	u.Path = strings.TrimRight(u.Path, "/") + path
	q := u.Query()
	q.Set("campaign", input.Campaign)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
