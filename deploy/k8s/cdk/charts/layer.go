package charts

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	. "github.com/medmouine/thesis/deploy/helpers"
	"github.com/medmouine/thesis/deploy/internal"
)

func NewLayer(scope constructs.Construct, id *string, tier *string) cdk8s.Chart {
	location := id

	chart := cdk8s.NewChart(scope, id, &cdk8s.ChartProps{
		Labels: &map[string]*string{
			"location": location,
			"tier":     tier,
		},
		Namespace: id,
	})
	NewNS(chart, location)

	internal.NewMapper(chart, S("mapper"), location)
	db := internal.NewInfluxDB(chart, S("influxdb"), location, tier)
	internal.NewTelegraf(chart, S("telegraf"), location, tier, internal.TelegrafProps{
		InfluxDBURL: db.URL,
		Bucket:      S("sensor_data"),
		BrokerURL:   S("tcp://edgebroker.%s:1883", *location),
	})

	return chart
}
