package client

//
//import (

//)
//
//type ClientV5[T interface{}] struct {
//	cm        *autopaho.ConnectionManager
//	d         device.Device[T]
//	SubTopics []string
//	DataTopic string
//	mux       sync.Mutex
//}
//
//type OptionsV5 struct {
//	MQTTOptions config.MqttConfig
//	SubTopics   []string
//	DataTopic   string
//	StateTopics []string
//}
//
//func NewConnectionV5[T interface{}](ctx context.Context, d device.Device[T], opts config.MqttConfig) (*ClientV5[T], error) {
//	cm, err := autopaho.NewConnection(ctx, opts.ToV5Config())
//	c := &ClientV5[T]{
//		cm:        cm,
//		d:         d,
//		SubTopics: opts.SubTopics,
//		DataTopic: opts.DataTopic,
//	}
//	c.cm.Subscribe(ctx, &paho.Subscribe{
//		Subscriptions: map[string]paho.SubscribeOptions{
//			opts.DataTopic: {QoS: 0},
//		},
//	})
//	return c, err
//}
