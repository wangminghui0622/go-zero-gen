
var (
	{{.lowerStartCamelObject}}FieldNames        = builder.RawFieldNames(&{{.upperStartCamelObject}}{}{{if .postgreSql}},true{{end}})
	{{.lowerStartCamelObject}}Rows              = strings.Join({{.lowerStartCamelObject}}FieldNames, ",")
	{{if .postgreSql}}
	{{.lowerStartCamelObject}}RowsExpectAutoSet = strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}",{{end}} "create_time", "update_time", "created_at", "updated_at"), ",")
	{{.lowerStartCamelObject}}RowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove({{.lowerStartCamelObject}}FieldNames, "{{.originalPrimaryKey}}", "create_time", "update_time", "created_at", "updated_at"))
	{{else}}
	{{.lowerStartCamelObject}}RowsExpectAutoSet = strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}",{{end}} "`create_time`", "`update_time`", "`created_at`", "`updated_at`"), ",")
	{{.lowerStartCamelObject}}RowsWithPlaceHolder = strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, "{{.originalPrimaryKey}}", "`create_time`", "`update_time`", "`created_at`", "`updated_at`"), "=?,") + "=?"
	{{end}}

	{{if .withCache}}{{.cacheKeys}}{{end}}
)
