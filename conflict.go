package sqrl

import (
	"io"
	"strings"
)

type conflictHandler struct {
	target Sqlizer
	action []Sqlizer
}

func (ch *conflictHandler) OnConflictKeys(conflictKeys ...string) {
	ch.target = Expr(strings.Join(conflictKeys, ","))
}

func (ch *conflictHandler) DoUpdateSet(clause ...Sqlizer) {
	ch.action = clause
}

func (ch *conflictHandler) DoUpdateSetKeys(keys ...string) {
	exprs := make([]Sqlizer, len(keys))
	for i, key := range keys {
		exprs[i] = Expr(key + "=EXCLUDED." + key)
	}
	ch.DoUpdateSet(exprs...)
}

func (ch *conflictHandler) AppendToSql(w io.Writer, args []interface{}) ([]interface{}, error) {
	if ch.target == nil && ch.action == nil {
		return args, nil
	}
	io.WriteString(w, " ON CONFLICT ")

	if ch.target != nil {
		targetSql, targetArgs, err := ch.target.ToSql()
		if err != nil {
			return args, err
		}
		io.WriteString(w, "("+targetSql+") ")
		args = append(args, targetArgs...)
	}

	if ch.action == nil {
		io.WriteString(w, "DO NOTHING")
	} else {
		io.WriteString(w, "DO UPDATE SET ")
		args, err := appendToSql(ch.action, w, ",", args)
		if err != nil {
			return args, err
		}
	}
	return args, nil
}
