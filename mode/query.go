package mode

import (
	"context"
	"fmt"
	"github.com/griesbacher/check_prometheus/helper"
	"github.com/griesbacher/check_x"
	"github.com/prometheus/common/model"
	"regexp"
	"strconv"
	"time"
)

//Query allows the user to test data in the prometheus server
func Query(address, query, warning, critical, alias, search, replace string) (err error) {
	warn, err := check_x.NewThreshold(warning)
	if err != nil {
		return
	}

	crit, err := check_x.NewThreshold(critical)
	if err != nil {
		return
	}
	var re *regexp.Regexp
	if search != "" {
		re, err = regexp.Compile(search)
		if err != nil {
			return
		}
	}

	apiClient, err := helper.NewAPIClientV1(address)
	if err != nil {
		return
	}

	result, err := apiClient.Query(context.TODO(), query, time.Now())
	if err != nil {
		return
	}
	switch result.Type() {
	case model.ValScalar:
		scalar := result.(*model.Scalar)
		scalarValue := float64(scalar.Value)
		helper.CheckTimestampFreshness(scalar.Timestamp)

		check_x.NewPerformanceData(replaceLabel("scalar", re, replace), scalarValue).Warn(warn).Crit(crit)
		state := check_x.Evaluator{Warning: warn, Critical: warn}.Evaluate(scalarValue)

		resultAsString := strconv.FormatFloat(scalarValue, 'f', -1, 64)
		if alias == "" {
			check_x.Exit(state, fmt.Sprintf("Query: '%s' returned: '%s'", query, resultAsString))
		} else {
			check_x.Exit(state, fmt.Sprintf("Alias: '%s' returned: '%s'", alias, resultAsString))
		}
	case model.ValVector:
		vector := result.(model.Vector)
		states := check_x.States{}
		for _, sample := range vector {
			helper.CheckTimestampFreshness(sample.Timestamp)

			sampleValue := float64(sample.Value)
			check_x.NewPerformanceData(replaceLabel(sample.Metric.String(), re, replace), sampleValue).Warn(warn).Crit(crit)
			states = append(states, check_x.Evaluator{Warning: warn, Critical: warn}.Evaluate(sampleValue))

		}
		return evalStates(states, alias, query)
	case model.ValMatrix:
		matrix := result.(model.Matrix)
		states := check_x.States{}
		for _, sampleStream := range matrix {
			for _, value := range sampleStream.Values {
				helper.CheckTimestampFreshness(value.Timestamp)
				states = append(states, check_x.Evaluator{Warning: warn, Critical: warn}.Evaluate(float64(value.Value)))
			}
		}
		return evalStates(states, alias, query)
	default:
		check_x.Exit(check_x.Unknown, fmt.Sprintf("The query did not return a suppoted type(scalar, vector, matrix), instead: '%s'. Query: '%s'", result.Type().String(), query))
		return nil
	}
	return err
}

func replaceLabel(label string, re *regexp.Regexp, replace string) string {
	if re != nil {
		label = re.ReplaceAllString(label, replace)
	}
	return label
}

func evalStates(states check_x.States, alias, query string) error {
	state, err := states.GetWorst()
	if err != nil {
		return err
	}
	if alias == "" {
		check_x.Exit(*state, fmt.Sprintf("Query: '%s'", query))
	} else {
		check_x.Exit(*state, fmt.Sprintf("Alias: '%s'", alias))
	}
	return nil
}
