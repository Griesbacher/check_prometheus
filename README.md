[![Circle CI](https://circleci.com/gh/Griesbacher/check_prometheus/tree/master.svg?style=svg)](https://circleci.com/gh/Griesbacher/check_prometheus/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Griesbacher/check_prometheus)](https://goreportcard.com/report/github.com/Griesbacher/check_prometheus)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)

# check_prometheus
Monitoring Plugin to check the health of a Prometheus server and its data

## Usage
### Global Options
```
$ check_prometheus -h
NAME:
   check_prometheus - Checks different prometheus stats as well the data itself
   Copyright (c) 2017 Philip Griesbacher
   https://github.com/Griesbacher/check_prometheus

USAGE:
   check_prometheus [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     mode, m  check mode
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -t value       Seconds till check returns unknown, 0 to disable (default: 10)
   -f value       If the checked data is older then this in seconds, unknown will be returned. Set to 0 to disable. (default: 300)
   --help, -h     show help
   --version, -v  print the version
```

### Command options

```
$ check_prometheus mode -h
   NAME:
      check_prometheus mode - check mode
   
   USAGE:
      check_prometheus mode command [command options] [arguments...]
   
   COMMANDS:
        ping, p         Returns the build informations
        query, q        Checks collected data
        targets_health  Returns the health of the targets
   
   OPTIONS:
      --help, -h  show help
```

### Subcommand options example

```
$ check_prometheus mode query -h
NAME:
   check_prometheus mode query - Checks collected data

USAGE:
   check_prometheus mode query [command options] [arguments...]

DESCRIPTION:
   Your Promqlquery has to return a vector / scalar / matrix result. The warning and critical values are applied to every value.
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
       --> OK - Query: 'up'|'prometheus'=1;;;; 'iapetos'=0;;;;

OPTIONS:
   --address value  Prometheus address: Protocol + IP + Port. (default: "http://localhost:9100")
   -q value         Query to be executed
   -a value         Alias, will replace the query within the output, if set
   -w value         Warning value. Use nagios-plugin syntax here.
   -c value         Critical value. Use nagios-plugin syntax here.
   --search value   If this variable is set, the given Golang regex will be used to search and replace the result with the 'replace' flag content. This will be appied on the perflabels.
   --replace value  See search flag. If the 'search' flag is empty this flag will be ignored.

```