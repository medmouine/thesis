package internal

const (
	LabelName     = "app.kubernetes.io/name"
	LabelLocation = "location"

	ImageMapper = "medmouine/mapper"

	TopicData  = "%s/sensor/mapper/+/temperature/data"
	TopicState = "%s/sensor/mapper/+/state"
)
