package sqrl

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// withExpr helps to alias part of the query generated with underlying "expr"
type withExpr struct {
	alias string
	expr  Sqlizer
}

func (e withExpr) ToSql() (sql string, args []interface{}, err error) {
	sql, args, err = e.expr.ToSql()
	if err == nil {
		sql = fmt.Sprintf("%s AS (%s)", e.alias, sql)
	}
	return
}

type WithBuilder struct {
	StatementBuilderType

	parts []withExpr
}

func NewWithBuilder(b StatementBuilderType) *WithBuilder {
	return &WithBuilder{
		StatementBuilderType: b,
		parts:                make([]withExpr, 0),
	}
}

func (b *WithBuilder) With(alias string, part Sqlizer) *WithBuilder {
	b.parts = append(b.parts, withExpr{alias: alias, expr: part})
	return b
}

func (b *WithBuilder) mergeWith(other *WithBuilder) *WithBuilder {
	b.parts = append(b.parts, other.parts...)
	return b
}

func (b *WithBuilder) ToSql() (sqlStr string, args []interface{}, err error) {
	sql := &bytes.Buffer{}
	sql.WriteString("WITH ")

	partsStr := make([]string, 0)

	var partSql string
	var partArgs []interface{}

	for _, e := range b.parts {
		partSql, partArgs, err = e.ToSql()
		if err != nil {
			return
		}
		partsStr = append(partsStr, partSql)
		args = append(args, partArgs...)
	}
	sql.WriteString(strings.Join(partsStr, ","))

	sqlStr, err = b.placeholderFormat.ReplacePlaceholders(sql.String())
	return
}

func (b *WithBuilder) AppendToSql(w io.Writer, args []interface{}) ([]interface{}, error) {
	cteSql, cteArgs, err := b.ToSql()
	if err != nil {
		return nil, err
	}
	io.WriteString(w, cteSql)
	args = append(args, cteArgs...)
	return args, nil
}

func (b *WithBuilder) Select(columns ...string) *SelectBuilder {
	return NewSelectBuilder(b.StatementBuilderType).Select(columns...).With(b)
}
