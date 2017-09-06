package mode

import (
	"github.com/griesbacher/check_prometheus/helper"
	"net/url"
	"path"
	"fmt"
	"time"
	"encoding/json"
	"strings"
	"github.com/griesbacher/check_x"
)

type targets struct {
	Status string `json:"status"`
	Data   struct {
		       ActiveTargets []struct {
			       DiscoveredLabels struct {
							Address     string `json:"__address__"`
							MetricsPath string `json:"__metrics_path__"`
							Scheme      string `json:"__scheme__"`
							Job         string `json:"job"`
						} `json:"discoveredLabels"`
			       Labels           struct {
							Instance string `json:"instance"`
							Job      string `json:"job"`
						} `json:"labels"`
			       ScrapeURL        string    `json:"scrapeUrl"`
			       LastError        string    `json:"lastError"`
			       LastScrape       time.Time `json:"lastScrape"`
			       Health           string    `json:"health"`
		       } `json:"activeTargets"`
	       } `json:"data"`
}

func getTargets(address string) (*targets, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/api/v1/targets")
	jsonBytes, err := helper.DoAPIRequest(u.String())
	if err != nil {
		return nil, err
	}
	var dat targets
	if err = json.Unmarshal(jsonBytes, &dat); err != nil {
		return nil, err
	}
	return &dat, nil
}

func TargetsHealth(address, warning, critical string) (err error) {
	warn, err := check_x.NewThreshold(warning)
	if err != nil {
		return
	}

	crit, err := check_x.NewThreshold(critical)
	if err != nil {
		return
	}

	targets, err := getTargets(address)
	if err != nil {
		return
	}
	if (*targets).Status != "success" {
		return fmt.Errorf("The API target returnstatus was %s", (*targets).Status)
	}
	msg := ""
	healthy := 0
	unhealthy := 0
	for _, target := range (*targets).Data.ActiveTargets {
		msg += fmt.Sprintf("Job: %s, Instance: %s, Health: %s, Last Error: %s\n", target.Labels.Job, target.Labels.Instance, target.Health, target.LastError)
		health := 0.0
		if target.Health != "up" {
			health = 1
			unhealthy += 1
		} else {
			healthy += 1
		}
		check_x.NewPerformanceData(target.Labels.Job, health)
	}
	var health_rate float64
	if unhealthy == 0{
		health_rate = 1
	}else {
		health_rate = float64(healthy) / float64(len((*targets).Data.ActiveTargets))
	}
	check_x.NewPerformanceData("health_rate", health_rate).Warn(warn).Crit(crit).Min(0).Max(1)
	check_x.NewPerformanceData("targets", float64(len((*targets).Data.ActiveTargets))).Min(0)
	state := check_x.Evaluator{Warning: warn, Critical: warn}.Evaluate(health_rate)
	check_x.LongExit(state, fmt.Sprintf("There are %d healthy and %d unhealthy targets", healthy, unhealthy), strings.TrimRight(msg, "\n"))
	return
}
