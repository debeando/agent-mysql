package main

import (
	"slices"
	"strings"
	"time"

	"github.com/debeando/agent-mysql/metrics"
	"github.com/debeando/go-common/env"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

func main() {
	log.Info("Start DeBeAndo Agent for MySQL")

	if getDebug() {
		log.SetLevel(log.DebugLevel)
	}

	log.DebugWithFields("Environment Variables", log.Fields{
		"DEBUG":           getDebug(),
		"DISABLE":         getDisableList(),
		"INFLUXDB_BUCKET": influxDB.Bucket,
		"INFLUXDB_HOST":   influxDB.Host,
		"INFLUXDB_PORT":   influxDB.Port,
		"INFLUXDB_TOKEN":  influxDB.Token,
		"INTERVAL":        getInterval(),
		"MYSQL_HOST":      MySQL.Host,
		"MYSQL_PASSWORD":  MySQL.Password,
		"MYSQL_PORT":      MySQL.Port,
		"MYSQL_TIMEOUT":   MySQL.Timeout,
		"MYSQL_USER":      MySQL.Username,
		"SERVER":          getServer(),
	})

	influxDB.New()
	defer influxDB.Close()
	MySQL.Connection = mysql.New(MySQL.Host, MySQL.DSN())
	defer MySQL.Connection.Close()

	for {
		items := metrics.Metrics{}

		MySQL.Connection.Connect()

		for _, query := range Queries {
			metric := metrics.Metric{}

			if !query.IsTime(query.Interval) {
				continue
			}

			if slices.Contains(getDisableList(), query.Name) {
				log.DebugWithFields("Ignore metric collect", log.Fields{
					"name": query.Name,
				})
				continue
			}

			MySQL.Connection.FetchAll(query.Beautifier(), func(row map[string]string) {
				metric.Measurement = query.Name
				metric.AddTag(metrics.Tag{
					Name:  "server",
					Value: getServer(),
				})

				if query.UnPivot {
					for column, value := range row {
						if valueParsed, ok := mysql.ParseNumberValue(value); ok {
							metric.AddField(metrics.Field{
								Name:  column,
								Value: valueParsed,
							})
						} else {
							metric.AddTag(metrics.Tag{
								Name:  column,
								Value: value,
							})
						}
					}
				} else if valueParsed, ok := mysql.ParseNumberValue(row[query.Value]); ok {
					metric.AddField(metrics.Field{
						Name:  row[query.Key],
						Value: valueParsed,
					})
				}
			})

			items.Add(metric)
		}

		if items.Count() > 0 {
			influxDB.Write(items)
		}

		items.Reset()
		log.Debug("Wait until next collect metrics.")
		time.Sleep(getInterval())
	}
}

func getDebug() bool {
	return env.GetBool("DEBUG", true)
}

func getInterval() time.Duration {
	return time.Duration(env.GetInt("INTERVAL", 10)) * time.Second
}

func getServer() string {
	return env.Get("SERVER", MySQL.Host)
}

func getDisableList() []string {
	list := strings.Split(env.Get("DISABLE", ""), ",")

	for item := range list {
		list[item] = strings.TrimSpace(list[item])
	}

	return list
}
