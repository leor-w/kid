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

type Wheres struct {
	Wheres []*Where
	isOr   bool
}

func And() *Wheres {
	return &Wheres{
		Wheres: make([]*Where, 0),
	}
}

func Or() *Wheres {
	return &Wheres{
		Wheres: make([]*Where, 0),
		isOr:   true,
	}
}

func (w *Wheres) And() *Wheres {
	return &Wheres{
		Wheres: make([]*Where, 0),
	}
}

func (w *Wheres) Or() *Wheres {
	return &Wheres{
		Wheres: make([]*Where, 0),
		isOr:   true,
	}
}

func (w *Wheres) add(where *Where) *Wheres {
	w.Wheres = append(w.Wheres, where)
	return w
}

func (w *Wheres) Eq(column string, value interface{}) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: "=",
		V1:        value,
		Or:        w.isOr,
	})
}

func (w *Wheres) Neq(column string, value interface{}) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: "<>",
		V1:        value,
		Or:        w.isOr,
	})
}

func (w *Wheres) Gt(column string, value interface{}) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: ">",
		V1:        value,
		Or:        w.isOr,
	})
}

func (w *Wheres) Gte(column string, value interface{}) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: ">=",
		V1:        value,
		Or:        w.isOr,
	})
}

func (w *Wheres) Lt(column string, value interface{}) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: "<",
		V1:        value,
		Or:        w.isOr,
	})
}

func (w *Wheres) Lte(column string, value interface{}) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: "<=",
		V1:        value,
		Or:        w.isOr,
	})
}

func (w *Wheres) In(column string, value interface{}) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: "IN",
		V1:        value,
		Or:        w.isOr,
	})
}

func (w *Wheres) NotIn(column string, value interface{}) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: "NOT IN",
		V1:        value,
		Or:        w.isOr,
	})
}

func (w *Wheres) Like(column string, value interface{}) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: "LIKE",
		V1:        value,
		Or:        w.isOr,
	})
}

func (w *Wheres) Between(column string, v1, v2 interface{}) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: "BETWEEN",
		V1:        v1,
		V2:        v2,
		Or:        w.isOr,
	})
}

func (w *Wheres) IsNull(column string) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: "IS NULL",
		Or:        w.isOr,
	})
}

func (w *Wheres) IsNotNull(column string) *Wheres {
	return w.add(&Where{
		Column:    column,
		Condition: "IS NOT NULL",
		Or:        w.isOr,
	})
}
