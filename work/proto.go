package work

import (
	"fmt"
	"gen/schemabuf"
	"gen/utils"
	"os"
)

func Proto(tablesUnit []schemabuf.UnitNoUnit, DBName string, file *os.File) {
	for _, one := range tablesUnit {
		a := `message ` + one.AddReq + ` {` + "\r\n" + `  ` + one.Unit + ` ` + one.Unit + ` = 1;` + "\r\n" + `}` + "\r\n"
		b := `message ` + one.AddResp + ` {` + "\r\n" + `}` + "\r\n"

		c := `message ` + one.QueryReq + ` {` + "\r\n" + `  int64 ` + one.UnitId + ` = 1;` + "\r\n" + `}` + "\r\n"
		d := `message ` + one.QueryResp + ` {` + "\r\n" + `  ` + one.Unit + ` ` + one.Unit + ` = 1;` + "\r\n" + `}` + "\r\n"

		e := `message ` + one.ModifyReq + ` {` + "\r\n" + `  ` + one.Unit + ` ` + one.Unit + ` = 1;` + "\r\n" + `}` + "\r\n"
		f := `message ` + one.ModifyResp + ` {` + "\r\n" + `}` + "\r\n"

		x := `message ` + one.DeleteReq + ` {` + "\r\n" + `  int64 ` + one.UnitId + ` = 1;` + "\r\n" + `}` + "\r\n"
		y := `message ` + one.DeleteResp + ` {` + "\r\n" + `}` + "\r\n" + `/////////////////////////////////////` + "\r\n"
		file.WriteString(a)
		file.WriteString(b)
		file.WriteString(c)
		file.WriteString(d)
		file.WriteString(e)
		file.WriteString(f)
		file.WriteString(x)
		file.WriteString(y)
		fmt.Print(a, b, c, d, e, f, x, y)
	}
	r := `service ` + utils.FirstToLower(DBName) + ` {` + "\r\n"
	for _, one := range tablesUnit {
		r += `  ` + `rpc ` + one.Add + ` (` + one.AddReq + `) ` + `returns` + ` (` + one.AddResp + `);` + "\r\n"
		r += `  ` + `rpc ` + one.Delete + ` (` + one.DeleteReq + `) ` + `returns` + ` (` + one.DeleteResp + `);` + "\r\n"
		r += `  ` + `rpc ` + one.Modify + ` (` + one.ModifyReq + `) ` + `returns` + ` (` + one.ModifyResp + `);` + "\r\n"
		r += `  ` + `rpc ` + one.Query + ` (` + one.QueryReq + `) ` + `returns` + ` (` + one.QueryResp + `);` + "\r\n" + "\r\n"
	}
	r += `}` + "\r\n"
	fmt.Println(r)
	file.WriteString(r)
}
