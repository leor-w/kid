package where

import "fmt"

const (
	IN        = "IN"
	NotIN     = "NOT IN"
	BETWEEN   = "BETWEEN"
	LIKE      = "LIKE"
	IsNULL    = "IS NULL"
	IsNotNULL = "IS NOT NULL"
)

type Where struct {
	Column    string
	Condition string
	V1        interface{}
	V2        interface{}
	Or        bool
}

func (w *Where) Build() (string, []interface{}) {
	if w.Condition == IN || w.Condition == NotIN {
		return fmt.Sprintf("%s %s (?)", w.Column, w.Condition), []interface{}{w.V1}
	}
	if w.Condition == BETWEEN {
		return fmt.Sprintf("%s %s ? AND ?", w.Column, w.Condition), []interface{}{w.V1, w.V2}
	}
	if w.Condition == LIKE {
		return fmt.Sprintf("%s %s ?", w.Column, w.Condition), []interface{}{w.V1}
	}
	if w.Condition == IsNULL || w.Condition == IsNotNULL {
		return fmt.Sprintf("%s %s", w.Column, w.Condition), nil
	}
	return fmt.Sprintf("%s %s ?", w.Column, w.Condition), []interface{}{w.V1}
}

type Wheres []*Where

func New() Wheres {
	return make(Wheres, 0)
}

func (w Wheres) And(wheres ...*Where) Wheres {
	return append(w, wheres...)
}

func (w Wheres) Or(wheres ...*Where) Wheres {
	for _, where := range wheres {
		where.Or = true
	}
	return append(w, wheres...)
}

func Eq(column string, value interface{}) *Where {
	return &Where{
		Column:    column,
		Condition: "=",
		V1:        value,
	}
}

func Neq(column string, value interface{}) *Where {
	return &Where{
		Column:    column,
		Condition: "<>",
		V1:        value,
	}
}

func Gt(column string, value interface{}) *Where {
	return &Where{
		Column:    column,
		Condition: ">",
		V1:        value,
	}
}

func Gte(column string, value interface{}) *Where {
	return &Where{
		Column:    column,
		Condition: ">=",
		V1:        value,
	}
}

func Lt(column string, value interface{}) *Where {
	return &Where{
		Column:    column,
		Condition: "<",
		V1:        value,
	}
}

func Lte(column string, value interface{}) *Where {
	return &Where{
		Column:    column,
		Condition: "<=",
		V1:        value,
	}
}

func In(column string, value interface{}) *Where {
	return &Where{
		Column:    column,
		Condition: "IN",
		V1:        value,
	}
}

func NotIn(column string, value interface{}) *Where {
	return &Where{
		Column:    column,
		Condition: "NOT IN",
		V1:        value,
	}
}

func Like(column string, value interface{}) *Where {
	return &Where{
		Column:    column,
		Condition: "LIKE",
		V1:        value,
	}
}

func Between(column string, v1, v2 interface{}) *Where {
	return &Where{
		Column:    column,
		Condition: "BETWEEN",
		V1:        v1,
		V2:        v2,
	}
}

func IsNull(column string) *Where {
	return &Where{
		Column:    column,
		Condition: "IS NULL",
	}
}

func IsNotNull(column string) *Where {
	return &Where{
		Column:    column,
		Condition: "IS NOT NULL",
		V1:        nil,
	}
}
