package data

import (
	"database/sql"
	"fmt"
	"gen/utils"
	"github.com/scylladb/go-set/strset"
	"os"
	"strings"
)

type Column struct {
	TableName              string
	ColumnName             string
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
}
type T struct {
	TableName  string `json:"tableName"`
	TableNamex string `json:"tableNamex"`
}

func Data(db *sql.DB) {
	var schema string
	err := db.QueryRow("SELECT SCHEMA()").Scan(&schema)

	q := "SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, " +
		"CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, COLUMN_TYPE " +
		"FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION"

	rows, err := db.Query(q, schema)
	defer rows.Close()
	if nil != err {
		return
	}
	cols := []Column{}
	set := strset.New()
	var setOrder []T
	for rows.Next() {
		var cs Column
		rows.Scan(&cs.TableName, &cs.ColumnName, &cs.IsNullable, &cs.DataType,
			&cs.CharacterMaximumLength, &cs.NumericPrecision, &cs.NumericScale, &cs.ColumnType)
		if cs.ColumnName == "index" || cs.ColumnName == "created_at" || cs.ColumnName == "updated_at" || cs.ColumnName == "deleted_at" {
			continue
		}
		cs.DataType = dataType(cs.DataType)
		cols = append(cols, cs)
		if !set.Has(cs.TableName) {
			set.Add(cs.TableName)
			setOrder = append(setOrder, T{
				TableName:  cs.TableName,
				TableNamex: utils.FirstLetterUpper(utils.SnakeToCamel(cs.TableName)),
			})
		}
	}
	dir := DealDir()
	apiComfile, err := os.Create(dir + `\api\` + "common" + ".api")
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer apiComfile.Close()
	apiHead(apiComfile, setOrder)

	protoComfile, err := os.Create(dir + `\proto\` + "common" + ".proto")
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer protoComfile.Close()
	protoHead(protoComfile, "zero")

	paramComfile, err := os.Create(dir + `\param\` + "common" + ".param")
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer paramComfile.Close()

	for _, table := range setOrder {
		var cs []Column
		for _, c := range cols {
			if table.TableName == c.TableName {
				cs = append(cs, c)
			}
		}
		apiUnit(apiComfile, table.TableName, cs)
		protoUnit(protoComfile, table.TableName, cs)
		paramUnit(paramComfile, table.TableName, cs)

		apiReqResp(dir, table.TableName, table.TableNamex)
		protoReqResp(dir, table.TableName, table.TableNamex)
	}
	return
}
func apiReqResp(dir, tableName, tableNamex string) {
	file, err := os.Create(dir + `\api\` + tableName + ".api")
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer file.Close()
	file.WriteString(`import "common.api"` + "\n")
	file.WriteString("\n")
	file.WriteString("type (" + "\n")

	file.WriteString("	" + tableNamex + "AddReq" + " {" + "\n")
	file.WriteString("		" + tableNamex + "Unit" + "\n")
	file.WriteString("	}")
	file.WriteString("\n")

	file.WriteString("	" + tableNamex + "AddResp" + " {" + "\n")
	file.WriteString("	}")
	file.WriteString("\n")

	file.WriteString("	" + tableNamex + "DeleteReq" + " {" + "\n")
	file.WriteString("		" + "Id int64 " + "`" + `json:"id"` + "`" + "\n")
	file.WriteString("	}")
	file.WriteString("\n")

	file.WriteString("	" + tableNamex + "DeleteResp" + " {" + "\n")
	file.WriteString("	}")
	file.WriteString("\n")

	file.WriteString("	" + tableNamex + "ModifyReq" + " {" + "\n")
	file.WriteString("		" + tableNamex + "Unit" + "\n")
	file.WriteString("	}")
	file.WriteString("\n")

	file.WriteString("	" + tableNamex + "ModifyResp" + " {" + "\n")
	file.WriteString("	}")
	file.WriteString("\n")

	file.WriteString("	" + tableNamex + "QueryReq" + " {" + "\n")
	file.WriteString("		" + "Id int64 " + "`" + `json:"id"` + "`" + "\n")
	file.WriteString("	}")
	file.WriteString("\n")

	file.WriteString("	" + tableNamex + "QueryResp" + " {" + "\n")
	file.WriteString("		" + tableNamex + "Unit" + "\n")
	file.WriteString("	}")
	file.WriteString("\n")

	file.WriteString(")")

	file.WriteString("////////////////////////////////////////////")
	file.WriteString("\n")

	file.WriteString("@server (" + "\n")

	file.WriteString("	prefix: " + `v1/` + utils.FirstLetterLower(tableNamex) + "\n")
	file.WriteString("	jwt: Auth" + "\n")
	file.WriteString("	tags: " + utils.FirstLetterLower(tableNamex) + "\n")
	file.WriteString("	group: " + utils.FirstLetterLower(tableNamex) + "\n")

	file.WriteString("	middleware: TokenParseMiddleware, TokenRefreshMiddleware, TokenBlacklistMiddleware " + "\n")
	file.WriteString(")")
	file.WriteString("\n")
	file.WriteString("service " + "zero" + " {" + "\n")

	file.WriteString("	@doc " + `"` + utils.FirstLetterLower(tableNamex) + "Add" + `"` + "\n")
	file.WriteString("	@handler " + utils.FirstLetterLower(tableNamex) + "Add" + "\n")

	file.WriteString("	post " + "/" + "add")
	file.WriteString(" (" + tableNamex + "AddReq" + ") " + "returns" + " (" + tableNamex + "AddResp" + ")" + "\n")
	file.WriteString("\n")

	file.WriteString("	@doc " + `"` + utils.FirstLetterLower(tableNamex) + "Delete" + `"` + "\n")
	file.WriteString("	@handler " + utils.FirstLetterLower(tableNamex) + "Delete" + "\n")

	file.WriteString("	post " + "/" + "delete")
	file.WriteString(" (" + tableNamex + "DeleteReq" + ") " + "returns" + " (" + tableNamex + "DeleteResp" + ")" + "\n")
	file.WriteString("\n")

	file.WriteString("	@doc " + `"` + utils.FirstLetterLower(tableNamex) + "Modify" + `"` + "\n")
	file.WriteString("	@handler " + utils.FirstLetterLower(tableNamex) + "Modify" + "\n")

	file.WriteString("	post " + "/" + "edit")
	file.WriteString(" (" + tableNamex + "ModifyReq" + ") " + "returns" + " (" + tableNamex + "ModifyResp" + ")" + "\n")
	file.WriteString("\n")

	file.WriteString("	@doc " + `"` + utils.FirstLetterLower(tableNamex) + "Query" + `"` + "\n")
	file.WriteString("	@handler " + utils.FirstLetterLower(tableNamex) + "Query" + "\n")

	file.WriteString("	post " + "/" + "query")
	file.WriteString(" (" + tableNamex + "QueryReq" + ") " + "returns" + " (" + tableNamex + "QueryResp" + ")" + "\n")
	file.WriteString("\n")

	file.WriteString("}")
}
func protoReqResp(dir, tableName, tableNamex string) {
	file, err := os.Create(dir + `\proto\` + tableName + ".proto")
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer file.Close()

	file.WriteString(`syntax = 'proto3';` + "\n")
	file.WriteString(`option go_package="./zero";` + "\n")
	file.WriteString(`package pb;` + "\n")
	file.WriteString(`import "common.proto";` + "\n")
	file.WriteString("\n")

	file.WriteString("message")

	file.WriteString(" " + tableNamex + "AddReq" + " {" + "\n")
	file.WriteString("		" + tableNamex + "Unit " + utils.FirstLetterLower(tableNamex) + "Unit = 1;" + "\n")
	file.WriteString("}")
	file.WriteString("\n")

	file.WriteString("message")
	file.WriteString(" " + tableNamex + "AddResp" + " {" + "\n")
	file.WriteString("}")
	file.WriteString("\n")

	file.WriteString("message")
	file.WriteString(" " + tableNamex + "DeleteReq" + " {" + "\n")
	file.WriteString("		" + "int64 id = 1;" + "\n")
	file.WriteString("}")
	file.WriteString("\n")

	file.WriteString("message")
	file.WriteString(" " + tableNamex + "DeleteResp" + " {" + "\n")
	file.WriteString("}")
	file.WriteString("\n")

	file.WriteString("message")
	file.WriteString(" " + tableNamex + "ModifyReq" + " {" + "\n")
	file.WriteString("		" + tableNamex + "Unit " + utils.FirstLetterLower(tableNamex) + "Unit = 1;" + "\n")
	file.WriteString("}")
	file.WriteString("\n")

	file.WriteString("message")
	file.WriteString(" " + tableNamex + "ModifyResp" + " {" + "\n")
	file.WriteString("}")
	file.WriteString("\n")

	file.WriteString("message")
	file.WriteString(" " + tableNamex + "QueryReq" + " {" + "\n")
	file.WriteString("		" + "int64 id = 1;" + "\n")
	file.WriteString("}")
	file.WriteString("\n")

	file.WriteString("message")
	file.WriteString(" " + tableNamex + "QueryResp" + " {" + "\n")
	file.WriteString("		" + tableNamex + "Unit " + utils.FirstLetterLower(tableNamex) + "Unit = 1;" + "\n")
	file.WriteString("}")
	file.WriteString("\n")

	file.WriteString("\n")

	file.WriteString("service " + utils.FirstLetterLower(tableNamex) + "{" + "\n")

	file.WriteString("  rpc " + utils.FirstLetterLower(tableNamex) + "Add" + "(" + tableNamex + "AddReq)" + " returns " + "(" + tableNamex + "AddResp" + ");" + "\n")

	file.WriteString("  rpc " + utils.FirstLetterLower(tableNamex) + "Delete" + "(" + tableNamex + "DeleteReq)" + " returns " + "(" + tableNamex + "DeleteResp" + ");" + "\n")

	file.WriteString("  rpc " + utils.FirstLetterLower(tableNamex) + "Modify" + "(" + tableNamex + "ModifyReq)" + " returns " + "(" + tableNamex + "ModifyResp" + ");" + "\n")

	file.WriteString("  rpc " + utils.FirstLetterLower(tableNamex) + "Query" + "(" + tableNamex + "QueryReq)" + " returns " + "(" + tableNamex + "QueryResp" + ");" + "\n")

	file.WriteString("}")
}

func protoUnit(file *os.File, table string, cs []Column) {
	file.WriteString("message ")
	file.WriteString(utils.FirstLetterUpper(utils.SnakeToCamel(table)) + "Unit" + " {")
	file.WriteString("\n")
	for n, ele := range cs {
		if ele.ColumnName == "updatedAt" || ele.ColumnName == "deletedAt" {
			continue
		}
		t := ele.DataType
		if ele.DataType == "bytes" {
			t = "string"
		}
		if ele.ColumnName == "index" {
			continue
		}
		file.WriteString("		" + t + " " + utils.SnakeToCamel(ele.ColumnName) + " =" + fmt.Sprintf("%d;", n+1))
		file.WriteString("\n")
	}
	file.WriteString("}")
	file.WriteString("\n")
	file.WriteString("////////////////////////////////////////////")
	file.WriteString("\n")
}
func apiHead(file *os.File, t []T) {
	file.WriteString("syntax = " + `"v1"` + "\n")
	file.WriteString("\n")
	file.WriteString("info ( \n")
	file.WriteString("	title: " + `""` + "\n")
	file.WriteString("	author: " + `""` + "\n")
	file.WriteString(")")
	file.WriteString("\n")
	for _, v := range t {
		file.WriteString(`//import "desc` + "/" + v.TableName + ".api" + `"`)
		file.WriteString("\n")
	}
}
func protoHead(file *os.File, table string) {
	file.WriteString(fmt.Sprintf("syntax = '%s';\n", "proto3"))
	file.WriteString(`option go_package="./zero";` + "\n")
	file.WriteString(fmt.Sprintf("package pb;\n"))
	file.WriteString("\n")
}
func apiUnit(file *os.File, table string, cs []Column) {
	file.WriteString("\n")
	file.WriteString("type (")
	file.WriteString("\n")
	file.WriteString("	" + utils.FirstLetterUpper(utils.SnakeToCamel(table)) + "Unit" + " {")
	file.WriteString("\n")
	for _, ele := range cs {
		if ele.ColumnName == "updatedAt" || ele.ColumnName == "deletedAt" {
			continue
		}
		t := ele.DataType
		if ele.DataType == "float" {
			t = "float32"
		}
		if ele.DataType == "bytes" {
			t = "string"
		}
		name := ele.ColumnName
		if ele.ColumnName == "index" {
			//name = ele.Name + "," + "optional"
			continue
		}
		if ele.ColumnName == "id" || ele.ColumnName == "createdAt" {
			file.WriteString("		" + utils.FirstLetterUpper(utils.SnakeToCamel(ele.ColumnName)) + " " + t + " " + "`json:" + `"` + utils.SnakeToCamel(name) + ",optional" + `"` + "`")
			file.WriteString("\n")
		} else {
			file.WriteString("		" + utils.FirstLetterUpper(utils.SnakeToCamel(ele.ColumnName)) + " " + t + " " + "`json:" + `"` + utils.SnakeToCamel(name) + `"` + "`")
			file.WriteString("\n")
		}
	}
	file.WriteString("	}")
	file.WriteString("\n")
	file.WriteString(")")
	file.WriteString("\n")
	file.WriteString("////////////////////////////////////////////")
}
func paramUnit(file *os.File, table string, cs []Column) {
	file.WriteString("\n")
	file.WriteString("	" + utils.FirstLetterUpper(utils.SnakeToCamel(table)) + "Unit" + " {")
	file.WriteString("\n")
	for n, ele := range cs {
		if ele.ColumnName == "updatedAt" || ele.ColumnName == "deletedAt" {
			continue
		}
		if ele.ColumnName == "index" {
			//name = ele.Name + "," + "optional"
			continue
		}
		if n == len(cs)-1 {
			if ele.DataType == "float" || ele.DataType == "int64" {
				file.WriteString("		" + `"` + utils.SnakeToCamel(ele.ColumnName) + `"` + ":" + " 0")
				file.WriteString("\n")
			} else {
				file.WriteString("		" + `"` + utils.SnakeToCamel(ele.ColumnName) + `"` + ":" + ` ""`)
				file.WriteString("\n")
			}
		} else {
			if ele.DataType == "float" || ele.DataType == "int64" {
				file.WriteString("		" + `"` + utils.SnakeToCamel(ele.ColumnName) + `"` + ":" + " 0,")
				file.WriteString("\n")
			} else {
				file.WriteString("		" + `"` + utils.SnakeToCamel(ele.ColumnName) + `"` + ":" + ` "",`)
				file.WriteString("\n")
			}
		}

	}
	file.WriteString("	}")
	file.WriteString("\n")
	file.WriteString(")")
	file.WriteString("\n")
	file.WriteString("////////////////////////////////////////////")
}
func DealDir() string {
	path, _ := os.Getwd()
	dir := path + `\gen_proto_api_param_models`
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Printf("remove dir %s err: %v\n", dir, err)
		return ""
	}
	fmt.Printf("remove dir %s sucess!", dir)
	os.Mkdir(dir, 0755)
	os.Mkdir(dir+`\`+"api", 0755)
	os.Mkdir(dir+`\`+"proto", 0755)
	os.Mkdir(dir+`\`+"param", 0755)
	return dir
}
func dataType(dataType string) string {
	typ := strings.ToLower(dataType)
	var dataType_ string
	switch typ {
	case "char", "varchar", "text", "longtext", "mediumtext", "tinytext":
		dataType_ = "string"
	case "json":
		dataType_ = "bytes"
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		dataType_ = "bytes"
	case "date", "time", "datetime", "timestamp":
		dataType_ = "int64"
	case "bool":
		dataType_ = "bool"
	case "tinyint":
		dataType_ = "int64"
	case "bigint":
		dataType_ = "int64"
	case "smallint", "int", "mediumint":
		dataType_ = "int64"
	case "float", "decimal", "double":
		dataType_ = "float"
	}
	return dataType_
}
