package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func GetTeamSchedule(
	client *http.Client,
	season string,
	eventCode string,
	teamNumber string,
	apiKey string,
) (resp *http.Response, err error) {
	return FTCAPICall(
		client,
		"GET",
		fmt.Sprintf("https://ftc-api.firstinspires.org/v2.0/%s/schedule/%s?teamNumber=%s", season, eventCode, teamNumber),
		map[string]string{},
		apiKey,
	)
}

func GetTeams(
	client *http.Client,
	season string,
	eventCode string,
	page int,
	apiKey string,
) (resp *http.Response, err error) {
	return FTCAPICall(
		client,
		"GET",
		fmt.Sprintf("https://ftc-api.firstinspires.org/v2.0/%s/teams?eventCode=%s&page=%d", season, eventCode, page),
		map[string]string{},
		apiKey,
	)
}

func FTCAPICall(
	client *http.Client,
	method string,
	url string,
	headers map[string]string,
	apiKey string,
) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	authToken := base64.StdEncoding.EncodeToString([]byte(apiKey))

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authToken))

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	return client.Do(req)
}
