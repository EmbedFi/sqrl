package sqrl

import (
	"io"
	"strings"
)

type ConflictHandler struct {
	target Sqlizer
	action []Sqlizer
}

func (ch *ConflictHandler) Target(targets ...string) {
	ch.target = Expr(strings.Join(targets, ","))
}

func (ch *ConflictHandler) DoNothing() {
	ch.action = nil
}

func (ch *ConflictHandler) DoUpdateSet(clause ...Sqlizer) {
	if ch.action == nil {
		ch.action = make([]Sqlizer, 0)
	}
	ch.action = append(ch.action, clause...)
}

func (ch *ConflictHandler) DoUpdateSetExcluded(keys ...string) {
	if ch.action == nil {
		ch.action = make([]Sqlizer, 0)
	}
	for _, key := range keys {
		ch.action = append(ch.action, Expr(key+"=EXCLUDED."+key))
	}
}

func (ch *ConflictHandler) AppendToSql(w io.Writer, args []interface{}) ([]interface{}, error) {
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
		var err error
		args, err = appendToSql(ch.action, w, ",", args)
		if err != nil {
			return args, err
		}
	}
	return args, nil
}
