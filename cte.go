package sqrl

import (
	"bytes"
	"io"
	"strings"
)

type WithBuilder struct {
	StatementBuilderType

	parts map[string]Sqlizer
}

func NewWithBuilder(b StatementBuilderType) *WithBuilder {
	return &WithBuilder{
		StatementBuilderType: b,
		parts:                make(map[string]Sqlizer),
	}
}

func (b *WithBuilder) With(alias string, part Sqlizer) *WithBuilder {
	b.parts[alias] = part
	return b
}

func (b *WithBuilder) mergeWith(other *WithBuilder) *WithBuilder {
	for alias, part := range other.parts {
		b.parts[alias] = part
	}
	return b
}

func (b *WithBuilder) ToSql() (sqlStr string, args []interface{}, err error) {
	sql := &bytes.Buffer{}
	sql.WriteString("WITH ")

	partsStr := make([]string, 0)

	var partSql string
	var partArg []interface{}

	for alias, expr := range b.parts {
		partSql, partArg, err = expr.ToSql()
		if err != nil {
			return
		}
		partsStr = append(partsStr, alias+" AS ("+partSql+")")
		args = append(args, partArg...)
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
