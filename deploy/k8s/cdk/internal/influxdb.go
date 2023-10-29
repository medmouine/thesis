package internal

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	. "github.com/medmouine/thesis/deploy/helpers"
)

type InfluxDB struct {
	constructs.Construct
	URL *string
}

func NewInfluxDB(scope constructs.Construct, id *string, location *string, tier *string) *InfluxDB {
	construct := constructs.NewConstruct(scope, id)
	name := *location + "-influxdb"

	cdk8s.NewHelm(construct, S("influxdb"), &cdk8s.HelmProps{
		Chart:       S("influxdb2"),
		ReleaseName: S(name),
		Repo:        S("https://helm.influxdata.com/"),
		Values: &map[string]interface{}{
			"nameOverride":     name,
			"fullnameOverride": name,
			"nodeSelector": &map[string]*string{
				LabelLocation: location,
				"tier":        tier,
			},
			"adminUser": map[string]interface{}{
				"organization":     name + "-org",
				"Bucket":           "sensor_data",
				"user":             "admin",
				"retention_policy": "0s",
				"password":         "influx_pw",
				"token":            "influx_token",
			},
			"persistence": map[string]interface{}{
				"enabled":    false,
				"accessMode": "ReadWriteOnce",
				"size":       "10Gi",
				"mountPath":  "/var/lib/influxdb2",
				"subPath":    "",
			},
		},
	})
	return &InfluxDB{
		Construct: construct,
		URL:       S("http://%s.%s:%s", name, *location, "8086"),
	}
}
