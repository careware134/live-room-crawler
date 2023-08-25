package domain

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
	Catchall       bool    `json:"catchall"`
	Intent         string  `json:"intent"`
	Confidence     float64 `json:"confidence"`
	NoImage        bool    `json:"no_image"`
	NoSpeech       bool    `json:"no_speech"`
	NoText         bool    `json:"no_text"`
	TriggeredQuery string  `json:"triggered_query"`
	TaskQA         TaskQA  `json:"taskqa"`
	Skill          Skill   `json:"skill"`
}

type ResponseStatus struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type QueryResponse struct {
	Trace          string         `json:"trace"`
	Answer         string         `json:"answer"`
	Text           string         `json:"text"`
	Query          string         `json:"query"`
	Meta           Meta           `json:"meta"`
	ResponseStatus ResponseStatus `json:"responseStatus"`
}

func (queryResponse *QueryResponse) ToPlayMessage() *CommandResponse {
	drivenType := TEXT
	if !queryResponse.Meta.NoSpeech {
		drivenType = AUDIO
	}
	if !queryResponse.Meta.NoText {
		drivenType = TEXT
	}
	response := &CommandResponse{
		CommandType: PLAY,
		TraceId:     queryResponse.Trace,
		Content: PlayContent{
			DrivenType: drivenType,
			Text:       queryResponse.Text,
			Audio:      queryResponse.Answer,
		},
	}
	return response
}
