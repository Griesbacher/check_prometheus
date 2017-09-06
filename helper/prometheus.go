package helper

import (
	"fmt"
	"github.com/griesbacher/check_x"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"io/ioutil"
	"net/http"
	"time"
)

//NewAPIClientV1 will create an prometheus api client v1
func NewAPIClientV1(address string) (v1.API, error) {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		return nil, err
	}
	return v1.NewAPI(client), nil
}

//DoAPIRequest does the http handling for an api request
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

//TimestampFreshness is the amount of second a result is treated as valid
var TimestampFreshness int

//CheckTimestampFreshness tests if the data is still valid
func CheckTimestampFreshness(timestamp model.Time) {
	CheckTimeFreshness(time.Unix(int64(timestamp), 0))
}

//CheckTimeFreshness tests if the data is still valid
func CheckTimeFreshness(timestamp time.Time) {
	if TimestampFreshness == 0 {
		return
	}
	timeDiff := time.Now().Sub(timestamp)
	if int(timeDiff.Seconds()) > TimestampFreshness {
		check_x.Exit(check_x.Unknown, fmt.Sprintf("One of the scraped data exceed the freshness by %ds", int(timeDiff.Seconds())-TimestampFreshness))
	}
}
