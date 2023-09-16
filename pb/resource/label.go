package resource

import (
	"fmt"
	"strings"
)

func ParseLabelRequirementListFromString(str string) (
	lables []*LabelRequirement) {
	items := strings.Split(str, "&")
	for i := range items {
		l := ParseLabelRequirementFromString(items[i])
		lables = append(lables, l)
	}
	return
}

// key=value1,value2,value3
func ParseLabelRequirementFromString(str string) *LabelRequirement {
	op := OPERATOR_IN
	var kvs = []string{}
	if strings.Contains(str, OPERATOR_NOT_IN.Expr()) {
		op = OPERATOR_NOT_IN
		kvs = strings.Split(str, OPERATOR_NOT_IN.Expr())
	} else {
		op = OPERATOR_IN
		kvs = strings.Split(str, OPERATOR_IN.Expr())
	}

	key := kvs[0]
	var value = []string{}
	if len(kvs) > 1 {
		v := strings.Join(kvs[1:], op.Expr())
		value = strings.Split(v, ",")
	}

	return NewLabelRequirement(key, value...)
}

// key1=value1,key2=value2
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
		return "nil"
	}
}

func ParseMapFromString(kvItems string) map[string]string {
	m := map[string]string{}
	kvs := strings.Split(kvItems, ",")
	for _, kv := range kvs {
		kv := strings.TrimSpace(kv)
		if kv == "" {
			continue
		}

		kvList := strings.Split(kv, "=")
		key, value := kvList[0], ""
		if len(kvList) > 1 {
			value = strings.Join(kvList[1:], "=")
		}
		m[key] = value
	}

	return m
}
