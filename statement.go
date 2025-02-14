package sqrl

// StatementBuilderType is the type of StatementBuilder.
type StatementBuilderType struct {
	placeholderFormat PlaceholderFormat
	runWith           BaseRunner
}

// With returns a WithBuilder for this StatementBuilder.
func (b StatementBuilderType) With(alias string, expr Sqlizer) *WithBuilder {
	return NewWithBuilder(b).With(alias, expr)
}

// Select returns a SelectBuilder for this StatementBuilder.
func (b StatementBuilderType) Select(columns ...string) *SelectBuilder {
	return NewSelectBuilder(b).Columns(columns...)
}

// Insert returns a InsertBuilder for this StatementBuilder.
func (b StatementBuilderType) Insert(into string) *InsertBuilder {
	return NewInsertBuilder(b).Into(into)
}

// Update returns a UpdateBuilder for this StatementBuilder.
func (b StatementBuilderType) Update(table string) *UpdateBuilder {
	return NewUpdateBuilder(b).Table(table)
}

// Delete returns a DeleteBuilder for this StatementBuilder.
func (b StatementBuilderType) Delete(what ...string) *DeleteBuilder {
	return NewDeleteBuilder(b).What(what...)
}

// PlaceholderFormat sets the PlaceholderFormat field for any child builders.
func (b StatementBuilderType) PlaceholderFormat(f PlaceholderFormat) StatementBuilderType {
	b.placeholderFormat = f
	return b
}

// RunWith sets the RunWith field for any child builders.
func (b StatementBuilderType) RunWith(runner BaseRunner) StatementBuilderType {
	b.runWith = wrapRunner(runner)
	return b
}

// StatementBuilder is a basic statement builder, holds global configuration options
// like placeholder format or SQL runner
var StatementBuilder = StatementBuilderType{placeholderFormat: Question}

// With returns a new WithBuilder, taking a first alias and expression
//
// See WithBuilder.With
func With(alias string, expr Sqlizer) *WithBuilder {
	return StatementBuilder.With(alias, expr)
}

// Select returns a new SelectBuilder, optionally setting some result columns.
//
// See SelectBuilder.Columns.
func Select(columns ...string) *SelectBuilder {
	return StatementBuilder.Select(columns...)
}

// Insert returns a new InsertBuilder with the given table name.
//
// See InsertBuilder.Into.
func Insert(into string) *InsertBuilder {
	return StatementBuilder.Insert(into)
}

// Update returns a new UpdateBuilder with the given table name.
//
// See UpdateBuilder.Table.
func Update(table string) *UpdateBuilder {
	return StatementBuilder.Update(table)
}

// Delete returns a new DeleteBuilder for given table names.
//
// See DeleteBuilder.Table.
func Delete(what ...string) *DeleteBuilder {
	return StatementBuilder.Delete(what...)
}

// Case returns a new CaseBuilder
// "what" represents case value
func Case(what ...interface{}) *CaseBuilder {
	b := &CaseBuilder{}

	switch len(what) {
	case 0:
	case 1:
		b = b.what(what[0])
	default:
		b = b.what(newPart(what[0], what[1:]...))

	}
	return b
}

// Values returns a new ValuesBuilder, with the provided arguments as the first row.
// Added to support `SelectBuilder.FromValues`
func Values(values ...interface{}) *ValuesBuilder {
	intialVals := make([][]interface{}, 0)
	if len(values) > 0 {
		intialVals = append(intialVals, values)
	}
	return &ValuesBuilder{
		StatementBuilderType: StatementBuilder,
		values:               intialVals,
	}
}
