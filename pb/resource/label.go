package resource

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

// 当Label的值为空或者*时匹配所有
func (l *LabelRequirement) IsMatchAll() bool {
	for i := range l.Values {
		v := l.Values[i]
		if v == "" || v == "*" {
			return true
		}
	}
	return false
}

func (l *LabelRequirement) MakeLabelKey(prefix string) string {
	if prefix == "" {
		return l.Key
	}

	return fmt.Sprintf("%s.%s", prefix, l.Key)
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
