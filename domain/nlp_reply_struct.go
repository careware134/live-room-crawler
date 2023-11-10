package domain

import "live-room-crawler/constant"

type QueryRequest struct {
	ProjectID int    `json:"projectId"`
	Query     string `json:"query"`
	SessionID string `json:"sessionId"`
	TraceID   string `json:"traceId"`
}

type TaskQA struct {
	TriggeredQuery string   `json:"triggered_query"`
	Category       string   `json:"category"`
	Label          []string `json:"label"`
}

type Skill struct {
	Name string `json:"name"`
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Meta struct {
	UserName       string  `json:"user_name"`
	Catchall       bool    `json:"catchall"`
	Intent         string  `json:"intent"`
	Confidence     float64 `json:"confidence"`
	NoImage        bool    `json:"no_image"`
	NoSpeech       bool    `json:"no_speech"`
	NoText         bool    `json:"no_text"`
	TriggeredQuery string  `json:"triggered_query"`
	TaskQA         *TaskQA `json:"taskqa"`
	Skill          *Skill  `json:"skill"`
}

type QueryResponse struct {
	Trace          string                  `json:"trace"`
	Answer         string                  `json:"answer"`
	Text           string                  `json:"text"`
	Query          string                  `json:"query"`
	Meta           Meta                    `json:"meta"`
	ResponseStatus constant.ResponseStatus `json:"responseStatus"`
}

func (queryResponse *QueryResponse) ToPlayMessage() *CommandResponse {
	drivenType := TEXT
	if !queryResponse.Meta.NoSpeech {
		drivenType = AUDIO
	}
	if !queryResponse.Meta.NoText {
		drivenType = TEXT
	}

	// judge play mode from label list
	playMode := HOST_MODE
	if queryResponse.Meta.TaskQA != nil && queryResponse.Meta.TaskQA.Label != nil {
		labelList := queryResponse.Meta.TaskQA.Label
		for _, label := range labelList {
			if label == string(ASSIST_MODE) {
				playMode = ASSIST_MODE
				break
			}
		}
	}

	response := &CommandResponse{
		CommandType: PLAY,
		TraceId:     queryResponse.Trace,
		Content: PlayContent{
			DrivenType: drivenType,
			Text:       queryResponse.Text,
			Audio:      queryResponse.Answer,
			PlayMode:   playMode,
		},
		RuleMeta: RuleMeta{
			Name:     queryResponse.Meta.TriggeredQuery,
			Query:    queryResponse.Query,
			UserName: queryResponse.Meta.UserName,
			Type:     CHAT,
		},
	}
	return response
}
