package label

import (
	"fmt"
	"strings"
)

func NewLabelRequirement(key string, values ...string) *LabelRequirement {
	return &LabelRequirement{
		Key:      key,
		Operator: OPERATOR_IN,
		Values:   values,
	}
}

func (l *LabelRequirement) Expr() string {
	return fmt.Sprintf("%s%s%s",
		l.Key,
		l.Operator.Expr(),
		strings.Join(l.Values, ","),
	)
}

func (o OPERATOR) Expr() string {
	switch o {
	case OPERATOR_IN:
		return "="
	case OPERATOR_NOT_IN:
		return "!="
	default:
		return ""
	}
}
