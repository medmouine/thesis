package main

import (
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/medmouine/thesis/deploy/charts"
	"github.com/medmouine/thesis/deploy/helpers"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

func main() {
	app := cdk8s.NewApp(nil)
	charts.NewLayer(app, helpers.S("edge1"))

	app.Synth()
}
