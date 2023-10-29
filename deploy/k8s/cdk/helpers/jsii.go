package helpers

import (
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/samber/lo"
)

func E(k string, v *string) lo.Entry[string, *string] {
	return lo.Entry[string, *string]{Key: k, Value: v}
}

func F[T any](v T, err error) func() (T, error) {
	return func() (T, error) {
		return v, err
	}
}

func N(n float64) *float64 {
	return jsii.Number(n)
}

func S(s string, f ...any) *string {
	if len(f) > 0 {
		s = fmt.Sprintf(s, f...)
	}
	return jsii.String(s)
}
