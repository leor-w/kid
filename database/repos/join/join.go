package join

import "fmt"

type Joins []*Join

type Join struct {
	LeftTable  string
	RightTable string
	LeftField  string
	RightField string
	JoinType   Type
}

type Type string

const (
	LEFT_JOIN  Type = "LEFT"
	RIGHT_JOIN Type = "RIGHT"
	INNER_JOIN Type = "INNER"
)

func (j *Join) Build() string {
	return fmt.Sprintf("%s JOIN %s ON %s.%s = %s.%s", j.JoinType, j.RightTable, j.LeftTable, j.LeftField, j.RightTable, j.RightField)
}

func JoinOn(leftTable, rightTable, leftField, rightField string, joinType Type) *Join {
	return &Join{
		LeftTable:  leftTable,
		RightTable: rightTable,
		LeftField:  leftField,
		RightField: rightField,
		JoinType:   joinType,
	}
}
