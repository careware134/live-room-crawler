package common

import (
	"github.com/google/uuid"
	"time"
)

type CounterType string

const (
	ONLINE  CounterType = "online"
	LIKE    CounterType = "like"
	GIFT    CounterType = "gift"
	FOLLOW  CounterType = "follow"
	VIEW    CounterType = "view"
	COMMENT CounterType = "comment"
)

type LiveStatisticsStruct struct {
	Online  StatisticCounter `json:"online,omitempty"`
	Like    StatisticCounter `json:"like,omitempty"`
	Gift    StatisticCounter `json:"gift,omitempty"`
	Follow  StatisticCounter `json:"follow,omitempty"`
	View    StatisticCounter `json:"view,omitempty"`
	Comment StatisticCounter `json:"comment,omitempty"`
}

type ActionType string

const (
	ON_COMMENT ActionType = "comment"
	ON_ENTER   ActionType = "enter"
	ON_GIFT    ActionType = "gift"
	ON_FOLLOW  ActionType = "follow"
)

type UserActionEvent struct {
	Username  string     `json:"username"`
	Action    ActionType `json:"action"`
	Content   string     `json:"content"`
	EventTime time.Time  ``
}

func (event *UserActionEvent) ToPlayMessage() *CommandResponse {
	response := &CommandResponse{
		CommandType: PLAY,
		TraceId:     "play-" + uuid.NewString(),
		Content: PlayContent{
			DrivenType: TEXT,
			Text:       event.Content,
		},
	}
	return response
}

type StatisticCounter struct {
	Count     uint64 `json:"count,omitempty"`
	Incr      bool   `json:"incr,omitempty"`
	LastMatch int    `json:"lastPush,omitempty"`
}

func BuildStatisticsCounter(count uint64, incr bool) *StatisticCounter {
	return &StatisticCounter{
		Count:     count,
		LastMatch: 0,
		Incr:      incr,
	}
}

func (c *StatisticCounter) AddCounter(count uint64) {
	c.Count += count
}

func (c *StatisticCounter) Add(other *StatisticCounter) {
	if other.Incr {
		c.Count += other.Count
	} else {
		c.Count = other.Count
	}

}

func InitStatisticStruct() map[CounterType]*StatisticCounter {
	registry := make(map[CounterType]*StatisticCounter)
	registry[ONLINE] = BuildStatisticsCounter(0, false)
	registry[LIKE] = BuildStatisticsCounter(0, false)
	registry[GIFT] = BuildStatisticsCounter(0, false)
	registry[FOLLOW] = BuildStatisticsCounter(0, false)
	registry[VIEW] = BuildStatisticsCounter(0, false)
	registry[COMMENT] = BuildStatisticsCounter(0, false)
	return registry

}
