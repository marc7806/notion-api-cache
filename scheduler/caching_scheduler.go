package scheduler

import (
	"log"
	"time"

	"github.com/marc7806/notion-cache/cache"
	"github.com/marc7806/notion-cache/config"
)

type Scheduler struct {
	timer *time.Timer
}

func Init() {
	log.Println("Initialize Caching Scheduler")
	sc := New()
	for {
		<-sc.timer.C
		sc.updateScheduler()
	}
}

func getDuration() time.Duration {
	now := time.Now()
	minuteTime := now.Minute() + config.CacheSchedulerMinutes
	hourTime := now.Hour() + config.CacheSchedulerHours
	dayTime := now.Day() + config.CacheSchedulerDays
	next := time.Date(now.Year(), now.Month(), dayTime, hourTime, minuteTime, now.Second(), 0, time.Local)
	return time.Until(next)
}

func New() Scheduler {
	return Scheduler{time.NewTimer(getDuration())}
}

func (sc Scheduler) updateScheduler() {
	log.Println("Caching scheduler triggered...")
	log.Println("Start notion databases cache refresh...")
	cache.HandleCacheRefresh()
	log.Println("Reset caching scheduler time interval...")
	sc.timer.Reset(getDuration())
}
