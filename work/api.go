package work

import (
	"fmt"
	"gen/schemabuf"
	"gen/utils"
	"os"
)

func api(s *schemabuf.Schema, tablesUnit []schemabuf.UnitNoUnit, db_name string) {
	os.Remove(db_name + ".api")
	file, err := os.Create(db_name + ".api")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("syntax = " + `"v1"` + "\n")
	file.WriteString("\n")
	file.WriteString("info ( \n")
	file.WriteString("	title: " + `""` + "\n")
	file.WriteString("	author: " + `""` + "\n")
	file.WriteString(")")
	file.WriteString("\n")

	file.WriteString("type (")
	file.WriteString("\n")
	for _, one := range s.Messages {
		file.WriteString("	" + one.Name + "Unit" + " {")
		file.WriteString("\n")
		for _, two := range tablesUnit {
			if one.Name == two.TableName {
				for _, ele := range one.Fields {
					if ele.Name == "updatedAt" || ele.Name == "deletedAt" {
						continue
					}
					t := ele.Typ
					if ele.Typ == "float" {
						t = "float64"
					}
					if ele.Typ == "bytes" {
						t = "string"
					}
					name := ele.Name
					if ele.Name == "index" {
						//name = ele.Name + "," + "optional"
						continue
					}
					if ele.Name == "id" || ele.Name == "createdAt" {
						file.WriteString("		" + utils.FirstLetterUpper(ele.Name) + " " + t + " " + "`json:" + `"` + name + ",optional" + `"` + "`")
						file.WriteString("\n")
					} else {
						file.WriteString("		" + utils.FirstLetterUpper(ele.Name) + " " + t + " " + "`json:" + `"` + name + `"` + "`")
						file.WriteString("\n")
					}
				}
				file.WriteString("	}")
			}
		}
		file.WriteString("\n")
	}

	file.WriteString("\n")
	file.WriteString(")")
	file.WriteString("\n")
	file.WriteString("////////////////////////////////////////////")
	file.WriteString("\n")

	file.WriteString("type (")
	file.WriteString("\n")

	for _, one := range tablesUnit {
		file.WriteString("	" + one.AddReq + " {" + "\n")
		file.WriteString("		" + one.Unit + " " + one.Unit + " " + "`json:" + `"` + utils.FirstLetterLower(one.Unit) + `"` + "`" + "\n")
		file.WriteString("	}")
		file.WriteString("\n")

		file.WriteString("	" + one.AddResp + " {" + "\n")
		//file.WriteString("		" + one.UnitId + " int32 " + " " + "`json:" + `"` + utils.FirstLetterLower(one.UnitId) + `"` + "`" + "\n")
		file.WriteString("	}")
		file.WriteString("\n")

		file.WriteString("	" + one.DeleteReq + " {" + "\n")
		file.WriteString("		" + one.UnitId + " int64 " + " " + "`json:" + `"` + utils.FirstLetterLower(one.UnitId) + `"` + "`" + "\n")
		file.WriteString("	}")
		file.WriteString("\n")

		file.WriteString("	" + one.DeleteResp + " {" + "\n")
		//file.WriteString("		" + one.UnitId + " int32 " + " " + "`json:" + `"` + utils.FirstLetterLower(one.UnitId) + `"` + "`" + "\n")
		file.WriteString("	}")
		file.WriteString("\n")

		file.WriteString("	" + one.ModifyReq + " {" + "\n")
		file.WriteString("		" + one.Unit + " " + one.Unit + " " + "`json:" + `"` + utils.FirstLetterLower(one.Unit) + `"` + "`" + "\n")
		file.WriteString("	}")
		file.WriteString("\n")

		file.WriteString("	" + one.ModifyResp + " {" + "\n")
		//file.WriteString("		" + one.UnitId + " int64 " + " " + "`json:" + `"` + utils.FirstLetterLower(one.UnitId) + `"` + "`" + "\n")
		file.WriteString("	}")
		file.WriteString("\n")

		file.WriteString("	" + one.QueryReq + " {" + "\n")
		file.WriteString("		" + one.UnitId + " int64 " + " " + "`json:" + `"` + utils.FirstLetterLower(one.UnitId) + `"` + "`" + "\n")
		file.WriteString("	}")
		file.WriteString("\n")

		file.WriteString("	" + one.QueryResp + " {" + "\n")
		file.WriteString("		" + one.Unit + " " + one.Unit + " " + "`json:" + `"` + utils.FirstLetterLower(one.Unit) + `"` + "`" + "\n")
		file.WriteString("	}")
		file.WriteString("\n")
		file.WriteString("\n")
		file.WriteString("\n")
		file.WriteString("\n")
	}

	file.WriteString("\n")
	file.WriteString(")")
	fmt.Println(tablesUnit, s)
	file.WriteString("////////////////////////////////////////////")
	file.WriteString("\n")
	for _, one := range tablesUnit {
		file.WriteString("@server (" + "\n")
		for _, two := range s.Messages {
			if one.TableName == two.Name {
				file.WriteString("	prefix: " + `v1/` + utils.FirstLetterLower(one.TableName) + "\n")
				file.WriteString("	group: " + utils.FirstLetterLower(one.TableName) + "\n")
				file.WriteString(")")
				file.WriteString("\n")
				file.WriteString("service " + utils.FirstLetterLower(db_name) + " {" + "\n")

				file.WriteString("	@doc " + `"` + utils.FirstLetterLower(one.Add) + `"` + "\n")
				file.WriteString("	@handler " + utils.FirstLetterLower(one.Add) + "\n")
				_, second1 := utils.SplitString(one.Add, len(one.TableName))
				file.WriteString("	post " + "/" + utils.FirstLetterLower(second1))
				file.WriteString(" (" + one.AddReq + ") " + "returns" + " (" + one.AddResp + ")" + "\n")
				file.WriteString("\n")

				file.WriteString("	@doc " + `"` + utils.FirstLetterLower(one.Delete) + `"` + "\n")
				file.WriteString("	@handler " + utils.FirstLetterLower(one.Delete) + "\n")
				_, second2 := utils.SplitString(one.Delete, len(one.TableName))
				file.WriteString("	post " + "/" + utils.FirstLetterLower(second2))
				file.WriteString(" (" + one.DeleteReq + ") " + "returns" + " (" + one.DeleteResp + ")" + "\n")
				file.WriteString("\n")

				file.WriteString("	@doc " + `"` + utils.FirstLetterLower(one.Modify) + `"` + "\n")
				file.WriteString("	@handler " + utils.FirstLetterLower(one.Modify) + "\n")
				_, second3 := utils.SplitString(one.Modify, len(one.TableName))
				file.WriteString("	post " + "/" + utils.FirstLetterLower(second3))
				file.WriteString(" (" + one.ModifyReq + ") " + "returns" + " (" + one.ModifyResp + ")" + "\n")
				file.WriteString("\n")

				file.WriteString("	@doc " + `"` + utils.FirstLetterLower(one.Query) + `"` + "\n")
				file.WriteString("	@handler " + utils.FirstLetterLower(one.Query) + "\n")
				_, second4 := utils.SplitString(one.Query, len(one.TableName))
				file.WriteString("	post " + "/" + utils.FirstLetterLower(second4))
				file.WriteString(" (" + one.QueryReq + ") " + "returns" + " (" + one.QueryResp + ")" + "\n")
				file.WriteString("\n")

				file.WriteString("}")

			}
		}
		file.WriteString("\n\n")
	}
}
