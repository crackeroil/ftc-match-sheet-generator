package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type TeamsResponse struct {
	Teams []struct {
		TeamNumber int    `json:"teamNumber"`
		NameShort  string `json:"nameShort"`
	} `json:"teams"`
	TeamCountPage int `json:"teamCountPage"`
	PageCurrent   int `json:"pageCurrent"`
	PageTotal     int `json:"pageTotal"`
}

type ScheduleResponse struct {
	Schedule []struct {
		MatchNumber int `json:"matchNumber"`
		Teams       []struct {
			TeamNumber int    `json:"teamNumber"`
			Station    string `json:"station"`
		} `json:"teams"`
	} `json:"schedule"`
}

func ParseTeamsResponse(resp *http.Response) (TeamsResponse, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return TeamsResponse{}, err
	}

	var out TeamsResponse
	if err = json.Unmarshal(body, &out); err != nil {
		panic(err)
	}

	return out, nil
}

func ParseScheduleResponse(resp *http.Response) (ScheduleResponse, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return ScheduleResponse{}, err
	}

	var out ScheduleResponse
	if err = json.Unmarshal(body, &out); err != nil {
		panic(err)
	}

	return out, nil
}
