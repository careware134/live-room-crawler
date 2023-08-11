package common

import (
	"github.com/google/uuid"
	"time"
)

type CounterType string

const (
	ONLINE  CounterType = "online"
	LIKE    CounterType = "enter"
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
	Count    uint64 `json:"count,omitempty"`
	Incr     bool   `json:"incr,omitempty"`
	LastPush int    `json:"lastPush,omitempty"`
}

func BuildStatisticsCounter(count uint64, incr bool) StatisticCounter {
	return StatisticCounter{
		Count:    count,
		LastPush: 0,
		Incr:     incr,
	}
}

func AddStatisticsCounter(base *StatisticCounter, count uint64) StatisticCounter {
	if base != nil {
		base.Count += count
		return *base
	}

	return BuildStatisticsCounter(count, true)
}

func (c *StatisticCounter) AddCounter(count uint64) {
	c.Count += count
}

func (c *StatisticCounter) Add(other StatisticCounter) {
	if other.Incr {
		c.Count += other.Count
	} else {
		c.Count = other.Count
	}

}

func InitStatisticStruct() LiveStatisticsStruct {
	return LiveStatisticsStruct{
		Online:  BuildStatisticsCounter(0, false),
		Like:    BuildStatisticsCounter(0, false),
		Gift:    BuildStatisticsCounter(0, false),
		Follow:  BuildStatisticsCounter(0, false),
		View:    BuildStatisticsCounter(0, false),
		Comment: BuildStatisticsCounter(0, false),
	}
}

func (s *LiveStatisticsStruct) Add(other LiveStatisticsStruct) {
	s.Online.Add(other.Online)
	s.Like.Add(other.Like)
	s.Gift.Add(other.Gift)
	s.Follow.Add(other.Follow)
	s.View.Add(other.View)
	s.Comment.Add(other.Comment)
}

func (s *LiveStatisticsStruct) AddCounter(counterType CounterType, other StatisticCounter) {
	if counterType == ONLINE {
		s.Online.Add(other)
	}
	if counterType == LIKE {
		s.Like.Add(other)
	}
	if counterType == GIFT {
		s.Gift.Add(other)
	}
	if counterType == FOLLOW {
		s.Follow.Add(other)
	}
	if counterType == VIEW {
		s.View.Add(other)
	}
	if counterType == COMMENT {
		s.Comment.Add(other)
	}
}
