package work

import (
	"fmt"
	"gen/schemabuf"
	"os"
)

func param(mc schemabuf.MessageCollection, tablesUnit []schemabuf.UnitNoUnit, db_name string) {
	os.Remove(db_name + ".param")
	file, err := os.Create(db_name + ".param")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for _, one := range mc {
		for _, two := range tablesUnit {
			if one.Name == two.TableName {
				file.WriteString(`{` + "\n")
				for _, ele := range one.Fields {
					t := ele.Typ
					fmt.Println(t)
					if t == "float" || t == "int64" || t == "bool" {
						file.WriteString("		" + `"` + ele.Name + `":` + ` 0,`)
					} else {
						file.WriteString("		" + `"` + ele.Name + `":` + ` "",`)
					}
					file.WriteString("\n")
				}
				file.WriteString("\n" + "}")
				file.WriteString("\n")
			}
		}
	}
}
