package internal

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	. "github.com/medmouine/thesis/deploy/helpers"
	"github.com/samber/lo"
)

type TelegrafProps struct {
	cdk8s.ChartProps
	InfluxDBURL *string
	Bucket      *string
	BrokerURL   *string
}

func NewTelegraf(
	scope constructs.Construct,
	id *string,
	location *string,
	tier *string,
	props TelegrafProps,
) constructs.Construct {
	construct := constructs.NewConstruct(scope, id)

	dburl := lo.FromPtrOr(props.InfluxDBURL, "http://influxdb2:8086")
	bucket := lo.FromPtrOr(props.Bucket, "temperature-sensor")
	brokerURL := lo.FromPtrOr(props.BrokerURL, "tcp://broker:1883")

	cdk8s.NewHelm(construct, S("telegraf"), &cdk8s.HelmProps{
		Chart:       S("telegraf"),
		ReleaseName: S("telegraf-" + *location),
		Repo:        S("https://helm.influxdata.com/"),
		Values: &map[string]interface{}{
			"replicaCount": 1,
			"env": []map[string]interface{}{
				{
					"name": "INFLUX_TOKEN",
					"valueFrom": map[string]interface{}{
						"secretKeyRef": map[string]interface{}{
							"name": "influxdb2-auth",
							"key":  "admin-token",
						},
					},
				},
			},
			"resources": map[string]interface{}{
				"limits": map[string]interface{}{
					"memory": "128Mi",
					"cpu":    "100m",
				},
			},
			"nodeSelector": map[string]interface{}{
				"tier":        *tier,
				LabelLocation: *location,
			},
			"config": map[string]interface{}{
				"agent": map[string]interface{}{
					"interval":            "10s",
					"round_interval":      true,
					"metric_batch_size":   1000,
					"metric_buffer_limit": 10000,
					"collection_jitter":   "0s",
					"flush_interval":      "10s",
					"flush_jitter":        "0s",
					"precision":           "",
					"debug":               false,
					"quiet":               false,
					"logfile":             "",
					"hostname":            "",
					"omit_hostname":       false,
				},
				"outputs": []map[string]interface{}{
					{
						"influxdb_v2": map[string]interface{}{
							"urls":         []string{dburl},
							"token":        "$INFLUX_TOKEN",
							"organization": *location + "-org",
							"Bucket":       bucket,
						},
					},
				},
				"inputs": []map[string]interface{}{
					{
						"mqtt_consumer": map[string]interface{}{
							"servers":          []string{brokerURL},
							"topics":           []string{fmt.Sprintf(TopicData, *location)},
							"data_format":      "json",
							"json_time_key":    "timestamp",
							"json_time_format": "unix",
						},
					},
				},
			},
			"service": map[string]interface{}{
				"enabled": false,
			},
		},
	})

	return construct
}
