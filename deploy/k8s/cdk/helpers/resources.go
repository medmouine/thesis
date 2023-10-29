package helpers

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	kplus "github.com/cdk8s-team/cdk8s-plus-go/cdk8splus24/v2"
	"github.com/samber/lo"
)

func NewNS(scope constructs.Construct, id *string) kplus.Namespace {
	return kplus.NewNamespace(scope, id, &kplus.NamespaceProps{Metadata: &cdk8s.ApiObjectMetadata{
		Name: id,
	}})
}

func GetNodes(selectors ...lo.Entry[string, *string]) kplus.LabeledNode {
	nodeQueries := lo.Compact(
		lo.Map(selectors, func(s lo.Entry[string, *string], _ int) kplus.NodeLabelQuery {
			if s.Value == nil {
				return nil
			}
			return kplus.NodeLabelQuery_Is(S(s.Key), s.Value)
		}))

	return kplus.NewLabeledNode(&nodeQueries)
}

func EnvStr(v string, f ...any) kplus.EnvValue {
	if len(f) > 0 {
		v = fmt.Sprintf(v, f...)
	}
	return kplus.EnvValue_FromValue(S(v))
}
