package ent

import (
	"os"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/DataDog/datadog-go/statsd"
)

type DBStatsD struct {
	poolName    string
	statsd      *statsd.Client
	driver      *entsql.Driver
	interval    time.Duration
	stop        chan bool
	defaultTags []string
}

func NewDBStatsD(pollName string, statsd *statsd.Client, driver *entsql.Driver) *DBStatsD {
	result := &DBStatsD{
		poolName: pollName,
		statsd:   statsd,
		driver:   driver,
		stop:     make(chan bool),
		interval: 10 * time.Second,
	}
	result.defaultTags = result.buildStaticTags()
	return result
}

func (d *DBStatsD) buildStaticTags() []string {
	tags := []string{"pool:" + d.poolName}

	if os.Getenv("HOSTNAME") != "" {
		tags = append(tags, "host:"+os.Getenv("HOSTNAME"))
	}

	if os.Getenv("POST_NAME") != "" {
		tags = append(tags, "post_name:"+os.Getenv("POST_NAME"))
	}
	return tags
}

func (d *DBStatsD) buildTags() []string {
	tags := d.defaultTags
	return tags
}
func (d *DBStatsD) sendStats() {
	currentStats := d.driver.DB().Stats()
	tags := d.buildTags()
	d.statsd.Gauge("db.max_open_connections", float64(currentStats.MaxOpenConnections), tags, 1)
	d.statsd.Gauge("db.open_connections", float64(currentStats.OpenConnections), tags, 1)
	d.statsd.Gauge("db.in_use", float64(currentStats.InUse), tags, 1)
	d.statsd.Gauge("db.idle", float64(currentStats.Idle), tags, 1)
	d.statsd.Gauge("db.wait_count", float64(currentStats.WaitCount), tags, 1)
	d.statsd.Gauge("db.wait_duration", float64(currentStats.WaitDuration.Milliseconds()), tags, 1)
	d.statsd.Gauge("db.max_idle_closed", float64(currentStats.MaxIdleClosed), tags, 1)
	d.statsd.Gauge("db.max_idle_time_closed", float64(currentStats.MaxIdleTimeClosed), tags, 1)
	d.statsd.Gauge("db.max_lifetime_closed", float64(currentStats.MaxLifetimeClosed), tags, 1)
	d.statsd.Gauge("db.connections", float64(currentStats.OpenConnections), tags, 1)
}

func (d *DBStatsD) Start() {
	go d.start()
}

func (d *DBStatsD) start() {
	ticker := time.NewTicker(d.interval)
	for {
		select {
		case <-ticker.C:
			d.sendStats()
		case <-d.stop:
			ticker.Stop()
			return
		}
	}
}

func (d *DBStatsD) Stop() {
	d.stop <- true
}
