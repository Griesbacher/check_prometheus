package helper

import (
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"time"
	"github.com/griesbacher/check_x"
	"fmt"
	"github.com/prometheus/common/model"
	"net/http"
	"io/ioutil"
)

func NewAPIClientV1(address string) (v1.API, error) {
	client, err := api.NewClient(api.Config{
		Address:address,
	})
	if err != nil {
		return nil, err
	}
	return v1.NewAPI(client), nil
}

func DoAPIRequest(address string) ([]byte, error) {
	resp, err := http.DefaultClient.Get(address)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

var TimestampFreshness = 120

func CheckTimestampFreshness(timestamp model.Time) {
	CheckTimeFreshness(time.Unix(int64(timestamp), 0))
}

func CheckTimeFreshness(timestamp time.Time) {
	if TimestampFreshness == 0 {
		return
	}
	timeDiff := time.Now().Sub(timestamp)
	if int(timeDiff.Seconds()) > TimestampFreshness {
		check_x.Exit(check_x.Unknown, fmt.Sprintf("One of the scraped data exceed the freshness by %ds", int(timeDiff.Seconds()) - TimestampFreshness))
	}
}