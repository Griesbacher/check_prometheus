package main

import (
	"github.com/griesbacher/check_prometheus/helper"
	"github.com/griesbacher/check_prometheus/mode"
	"github.com/griesbacher/check_x"
	"github.com/urfave/cli"
	"os"
	"time"
)

var (
	address  string
	timeout  int
	warning  string
	critical string
	query    string
	alias    string
	search   string
	replace  string
	label    string
)

func startTimeout() {
	if timeout != 0 {
		check_x.StartTimeout(time.Duration(timeout) * time.Second)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "check_prometheus"
	app.Usage = "Checks different prometheus stats as well the data itself\n   Copyright (c) 2017 Philip Griesbacher\n   https://github.com/Griesbacher/check_prometheus"
	app.Version = "0.0.1"
	flagAddress := cli.StringFlag{
		Name:        "address",
		Usage:       "Prometheus address: Protocol + IP + Port.",
		Destination: &address,
		Value:       "http://localhost:9100",
	}
	flagWarning := cli.StringFlag{
		Name:        "w",
		Usage:       "Warning value. Use nagios-plugin syntax here.",
		Destination: &warning,
	}
	flagCritical := cli.StringFlag{
		Name:        "c",
		Usage:       "Critical value. Use nagios-plugin syntax here.",
		Destination: &critical,
	}
	flagQuery := cli.StringFlag{
		Name:        "q",
		Usage:       "Query to be executed",
		Destination: &query,
	}
	flagAlias := cli.StringFlag{
		Name:        "a",
		Usage:       "Alias, will replace the query within the output, if set",
		Destination: &alias,
	}
	flagLabel := cli.StringFlag{
		Name:        "l",
		Usage:       "Prometheus-Label, which will be used for the performance data label. By default job and instance should be available.",
		Destination: &label,
		Value:       mode.DefaultLabel,
	}
	app.Commands = []cli.Command{
		{
			Name:    "mode",
			Aliases: []string{"m"},
			Usage:   "check mode",
			Subcommands: []cli.Command{
				{
					Name:        "ping",
					Aliases:     []string{"p"},
					Usage:       "Returns the build informations",
					Description: `This check requires that the prometheus server itself is listetd as target. Following query will be used: 'prometheus_build_info{job="prometheus"}'`,
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.Ping(address)
					},
					Flags: []cli.Flag{
						flagAddress,
					},
				}, {
					Name:    "query",
					Aliases: []string{"q"},
					Usage:   "Checks collected data",
					Description: `Your Promqlquery has to return a vector / scalar / matrix result. The warning and critical values are applied to every value.
   Examples:
       Vector:
           check_prometheus mode query -q 'up'
       --> OK - Query: 'up'|'up{instance="192.168.99.101:9245", job="iapetos"}'=1;;;; 'up{instance="0.0.0.0:9091", job="prometheus"}'=1;;;;

       Scalar:
           check_prometheus mode query -q 'scalar(up{job="prometheus"})'
       --> OK - OK - Query: 'scalar(up{job="prometheus"})' returned: '1'|'scalar'=1;;;;

       Matrix:
           check_prometheus mode query -q 'http_requests_total{job="prometheus"}[5m]'
       --> OK - Query: 'http_requests_total{job="prometheus"}[5m]'

       Search and Replace:
           check_prometheus m query -q 'up' --search '.*job=\"(.*?)\".*' --replace '$1'
       --> OK - Query: 'up'|'prometheus'=1;;;; 'iapetos'=0;;;;`,

					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.Query(address, query, warning, critical, alias, search, replace)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagQuery,
						flagAlias,
						flagWarning,
						flagCritical,
						cli.StringFlag{
							Name:        "search",
							Usage:       "If this variable is set, the given Golang regex will be used to search and replace the result with the 'replace' flag content. This will be appied on the perflabels.",
							Destination: &search,
						},
						cli.StringFlag{
							Name:        "replace",
							Usage:       "See search flag. If the 'search' flag is empty this flag will be ignored.",
							Destination: &replace,
						},
					},
				}, {
					Name:        "targets_health",
					Usage:       "Returns the health of the targets",
					Description: `The warning and critical thresholds are appied on the health_rate. The health_rate is calculted: sum(healthy) / sum(targets).`,
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.TargetsHealth(address, label, warning, critical)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagWarning,
						flagCritical,
						flagLabel,
					},
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "t",
			Usage:       "Seconds till check returns unknown, 0 to disable",
			Value:       10,
			Destination: &timeout,
		},
		cli.IntFlag{
			Name:        "f",
			Usage:       "If the checked data is older then this in seconds, unknown will be returned. Set to 0 to disable.",
			Value:       300,
			Destination: &helper.TimestampFreshness,
		},
	}

	if err := app.Run(os.Args); err != nil {
		check_x.ErrorExit(err)
	}
}
