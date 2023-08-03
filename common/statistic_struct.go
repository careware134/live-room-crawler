package common

import "time"

type UpdateRegistryEvent struct {
	Statistics LiveStatisticsStruct `json:"statistics,omitempty"`
	ActionList []UserActionStruct   `json:"actionList,omitempty"`
}

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
	COMMENT ActionType = "comment"
	ENTER   ActionType = "enter"
	GIFT    ActionType = "gift"
	FOLLOW  ActionType = "follow"
)

type UserActionStruct struct {
	Username  string     `json:"username"`
	Action    ActionType `json:"action"`
	Content   string     `json:"content"`
	EventTime time.Time  ``
}

type StatisticCounter struct {
	count uint64 `json:"count,omitempty"`
	incr  bool   `json:"incr,omitempty"`
}

func BuildStatisticsCounter(count uint64, incr bool) StatisticCounter {
	return StatisticCounter{
		count: count,
		incr:  incr,
	}
}

func AddStatisticsCounter(base *StatisticCounter, count uint64) StatisticCounter {
	if base != nil {
		base.count += count
		return *base
	}

	return BuildStatisticsCounter(count, true)
}

func (c *StatisticCounter) AddCounter(count uint64) {
	c.count += count
}
