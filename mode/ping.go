package mode

import (
	"time"
	"github.com/griesbacher/check_prometheus/helper"
	"context"
	"fmt"
	"github.com/prometheus/common/model"
	"encoding/json"
	"github.com/griesbacher/check_x"
)

type buildInfo struct {
	Metric struct {
		       Name      string `json:"__name__"`
		       Branch    string `json:"branch"`
		       Goversion string `json:"goversion"`
		       Instance  string `json:"instance"`
		       Job       string `json:"job"`
		       Revision  string `json:"revision"`
		       Version   string `json:"version"`
	       } `json:"metric"`
	Value []interface{} `json:"value"`
}

func Ping(address string) (err error) {
	apiClient, err := helper.NewAPIClientV1(address)
	if err != nil {
		return
	}
	query := `prometheus_build_info{job="prometheus"}`
	startTime := time.Now()
	result, err := apiClient.Query(context.TODO(), query, time.Now())
	endTime := time.Now()
	if err != nil {
		return
	}
	vector := result.(model.Vector)
	if len(vector) != 1 {
		return fmt.Errorf("The query '%s' did not return a vector with a single entry" , query)
	}
	sample := vector[0]
	helper.CheckTimestampFreshness(sample.Timestamp)
	jsonBytes, err := sample.MarshalJSON()
	if err != nil {
		return
	}
	var dat buildInfo
	if err = json.Unmarshal(jsonBytes, &dat); err != nil {
		return
	}
	check_x.NewPerformanceData("duration", endTime.Sub(startTime).Seconds()).Unit("s").Min(0)
	check_x.Exit(check_x.OK, fmt.Sprintf("Version: %s, Instance %s", dat.Metric.Version, dat.Metric.Instance))
	return err
}
