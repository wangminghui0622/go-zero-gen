package work

import (
	"fmt"
	"gen/schemabuf"
	"log"
	"os/exec"
	"sync"
)

func Worker(id int, wg *sync.WaitGroup, dbUser, dbPassword,
	dbHost string, dbPort int, dbName, tb string, s *schemabuf.Schema,
	tablesUnit []schemabuf.UnitNoUnit, db_name string) {
	defer wg.Done()
	if id == 1 {
		dbUrl := "-url=" + fmt.Sprintf("%s", dbUser) + ":" + fmt.Sprintf("%s", dbPassword) + "@" + "tcp(" + fmt.Sprintf("%s:%d", dbHost, dbPort) + ")" + fmt.Sprintf("/%s", dbName)
		cmd := exec.Command(
			"cmd", "/c",
			"goctl", "model", "mysql", "datasource",
			dbUrl,
			//"-t=box_device,box_apk_version",
			"-t="+fmt.Sprintf("%s", tb),
			"-dir=./models_base",
			"-c",
		)
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Command failed: %v\n", err)
		}
	} else {
		if id == 2 {
			api(s, tablesUnit, db_name)
		} else {
			param(s.Messages, tablesUnit, db_name)
		}
	}
}
