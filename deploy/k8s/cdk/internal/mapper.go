package internal

import (
	"strconv"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	kplus "github.com/cdk8s-team/cdk8s-plus-go/cdk8splus24/v2"
	. "github.com/medmouine/thesis/deploy/helpers"
	"github.com/samber/lo"
)

type MapperProps struct {
	cdk8s.ChartProps
	Name      *string
	Version   *string
	BrokerURL *string

	Replicas          *float64
	ExtraNodeSelector *map[string]*string
	ExtraEnvVars      *map[string]kplus.EnvValue
}

func NewMapper(scope constructs.Construct, id *string, location *string, mapperProps ...*MapperProps) constructs.Construct {
	construct := constructs.NewConstruct(scope, id)

	props := &MapperProps{}
	if len(mapperProps) > 0 {
		props = mapperProps[0]
	}

	v, _ := lo.Coalesce(props.Version, S("latest"))
	replicas, _ := lo.Coalesce(props.Replicas, N(1))
	name, _ := lo.Coalesce(props.Name, S("mapper"))
	labels, _ := lo.Coalesce(props.Labels, &map[string]*string{})
	brokerURL, _ := lo.Coalesce(props.BrokerURL, S("mqtt://edgebroker.%s:1883", *location))
	envVars, _ := lo.Coalesce(props.ExtraEnvVars, &map[string]kplus.EnvValue{})
	nodeSelector, _ := lo.Coalesce(props.ExtraNodeSelector, &map[string]*string{})

	labels = lo.ToPtr(lo.Assign(map[string]*string{
		LabelName:     name,
		LabelLocation: location,
	}, *labels))

	envVars = lo.ToPtr(lo.Assign(map[string]kplus.EnvValue{
		"DEVICE_LOCATION":       EnvStr(*location),
		"STAGE_STATUS":          EnvStr("dev"),
		"SERVER_HOST":           EnvStr("0.0.0.0"),
		"SERVER_PORT":           EnvStr("3000"),
		"SERVER_ADDR":           EnvStr("0.0.0.0:3000"),
		"SERVER_READ_TIMEOUT":   EnvStr("5s"),
		"SERVER_WRITE_TIMEOUT":  EnvStr("10s"),
		"SERVER_IDLE_TIMEOUT":   EnvStr("120s"),
		"MQTT_SUB_TOPICS":       EnvStr(TopicState+":"+TopicData, *location, *location),
		"MQTT_STATE_TOPICS":     EnvStr(TopicState, *location),
		"MQTT_DATA_TOPIC":       EnvStr(TopicData, *location),
		"MQTT_BROKER_URL":       EnvStr(*brokerURL),
		"MQTT_PUBLISH_INTERVAL": EnvStr("2s"),
		"MQTT_CLIENT_ID":        kplus.EnvValue_FromFieldRef(kplus.EnvFieldPaths_POD_UID, nil),
	}, *envVars))

	port, _ := lo.TryOr(F(strconv.ParseFloat((*envVars)["SERVER_PORT"].Value().(string), 64)), 3000.)
	image := S(ImageMapper + ":" + *v)

	depl := kplus.NewDeployment(construct, S("depl"), &kplus.DeploymentProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Labels: labels,
			Name:   name,
			//Namespace: ns,
		},
		Containers: &[]*kplus.ContainerProps{{
			Image:        image,
			EnvVariables: envVars,
			Ports: &[]*kplus.ContainerPort{{
				Number: N(port),
			}},
		}},
		Replicas: replicas,
	})

	nodeQueries := append(
		lo.Entries(*nodeSelector),
		E("location", location),
	)
	nodes := GetNodes(nodeQueries...)
	depl.Scheduling().Attract(nodes, nil)

	svc := kplus.NewService(construct, S("svc"), &kplus.ServiceProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Labels: labels,
			Name:   name,
			//Namespace: ns,
		},
	})

	svc.SelectLabel(S(LabelName), name)
	svc.SelectLabel(S(LabelLocation), location)
	svc.Bind(N(port), nil)
	return construct
}
