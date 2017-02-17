package githubstat

import (
	"time"

	"github.com/BurntSushi/toml"
)

const (
	DimensionOverall = "Overall"
	DimensionWeek    = "Week"

	NoSort              = 0
	SortByMergedPRs     = 1
	SortByMergedCommits = 2
)

type Load interface {
}

var Config Configuration

type User struct {
	Name     string
	RealName string
}

type Configuration struct {
	StatBeginTime    time.Time
	StatEndTime      time.Time
	AccessToken      string
	Users            []User
	Repos            []string
	Metrics          string
	Dimension        string
	WeekFirstDay     time.Weekday
	ThisWeekFirstDay time.Time
	Sort             int
}

func getWeekFirstDay(t time.Time) time.Time {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	for {
		weekday := t.Weekday()
		if weekday != Config.WeekFirstDay {
			t = t.Add(-24 * time.Hour)
		} else {
			break
		}
	}
	return t

}

// read config file
func init() {
	if _, err := toml.DecodeFile("./config.toml", &Config); err != nil {
		panic(err)
	}
	if !Config.StatEndTime.IsZero() && !Config.StatEndTime.After(Config.StatBeginTime) {
		panic("stat end time must be after stat begin time")
	}
	Config.ThisWeekFirstDay = getWeekFirstDay(time.Now())
}
