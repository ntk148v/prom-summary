module github.com/ntk148v/prom-summary

go 1.14

require (
	github.com/alecthomas/units v0.0.0-20210208195552-ff826a37aa15 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.15.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/prometheus/client_golang v1.9.0 => github.com/ntk148v/client_golang v1.9.1-0.20210306080221-9fa2d6cf89bd
