package main

import (
	"database/sql"
	"flag"
	"fmt"
	"gen/data"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	dbUser     = "root"
	dbPassword = "test789"
	dbHost     = "42.192.129.44"
	dbPort     = 3306
	dbName     = "balance"
)

func main() {
	dbType := flag.String("db", "mysql", "the database type")
	host := flag.String("host", dbHost, "the database host")
	port := flag.Int("port", dbPort, "the database port")
	user := flag.String("user", dbUser, "the database user")
	password := flag.String("password", dbPassword, "the database password")
	db_name := flag.String("schema", dbName, "the database schema")
	packageName := flag.String("package", *db_name, "the protocol buffer package. defaults to the database schema.")
	ignoreTableStr := flag.String("ignore_tables", "", "a comma spaced list of tables to ignore")
	flag.Parse()
	fmt.Println(packageName, ignoreTableStr)
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", *user, *password, *host, *port, *db_name)
	db, err := sql.Open(*dbType, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	go data.Data(db)
	
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v\n", err)
	}
	fmt.Printf("Current working directory: %s\n", workDir)
	
	// 确保输出目录存在（使用相对路径，与手动执行时一致）
	outputDir := "./gen_proto_api_param_models/models"
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v\n", err)
	}
	fmt.Printf("Output directory (relative): %s\n", outputDir)
	fmt.Printf("Output directory (absolute): %s\n", absOutputDir)
	
	if err := os.MkdirAll(absOutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v\n", err)
	}
	
	// 构建数据库连接URL
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Printf("Database URL: %s\n", dbUrl)
	
	// 检查 goctl 是否可用
	goctlPath, err := exec.LookPath("goctl")
	if err != nil {
		log.Fatalf("goctl not found in PATH: %v\nPlease install goctl: go install github.com/zeromicro/go-zero/tools/goctl@latest\n", err)
	}
	fmt.Printf("Found goctl at: %s\n", goctlPath)
	
	// 在 Windows 上，使用 cmd /c 执行命令
	// 使用相对路径，与手动执行时保持一致
	cmd := exec.Command(
		"cmd", "/c",
		"goctl", "model", "mysql", "datasource",
		"-url="+dbUrl,
		"-t=*",
		"-dir="+outputDir,  // 使用相对路径
		"-c",
	)
	
	// 设置工作目录为当前目录
	cmd.Dir = workDir
	
	// 打印完整命令
	fmt.Printf("Executing command: cmd /c goctl model mysql datasource -url=%s -t=* -dir=%s -c\n", dbUrl, outputDir)
	
	output, err := cmd.CombinedOutput()
	
	// 立即打印输出和错误，便于调试
	if cmd.ProcessState != nil {
		fmt.Printf("Command exit code: %d\n", cmd.ProcessState.ExitCode())
	}
	fmt.Printf("Command Output:\n%s\n", string(output))
	if err != nil {
		log.Fatalf("Command failed: %v\nOutput: %s\n", err, string(output))
	}
	
	// 等待一下，确保文件写入完成
	time.Sleep(500 * time.Millisecond)
	
	// 检查输出目录是否有文件生成（使用绝对路径检查）
	files, err := os.ReadDir(absOutputDir)
	if err != nil {
		log.Fatalf("Failed to read output directory: %v\n", err)
	}
	if len(files) == 0 {
		log.Fatalf("No files generated in %s. Command output: %s\n", absOutputDir, string(output))
	}
	fmt.Printf("Successfully generated %d files in %s\n", len(files), absOutputDir)
	for _, f := range files {
		fmt.Printf("  - %s\n", f.Name())
	}
	
	gen_models_advanced()
	gen_models_base()
	gen_advancedTomodels()
}
func gen_models_advanced() {
	p, _ := os.Getwd()
	p = p + `\gen_proto_api_param_models\models`
	genFiles := []string{}
	filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_gen.go") {
			genFiles = append(genFiles, info.Name())
		}
		return nil
	})
	destPath, _ := os.Getwd()
	// 使用项目目录，而不是父目录
	destDir := filepath.Join(destPath, "models_advanced")
	// 确保目标目录存在
	if err := os.MkdirAll(destDir, 0755); err != nil {
		log.Printf("Failed to create directory %s: %v\n", destDir, err)
		return
	}
	fmt.Printf("Copying _gen.go files to %s\n", destDir)
	for _, f := range genFiles {
		src := filepath.Join(p, f)
		dst := filepath.Join(destDir, f)
		if err := copyFile(src, dst); err != nil {
			log.Printf("Failed to copy %s to %s: %v\n", src, dst, err)
			continue
		}
		fmt.Printf("  Copied: %s -> %s\n", f, dst)
		os.Remove(src)
	}
	// 复制 vars.go
	varsSrc := filepath.Join(p, "vars.go")
	varsDst := filepath.Join(destDir, "vars.go")
	if _, err := os.Stat(varsSrc); err == nil {
		if err := copyFile(varsSrc, varsDst); err != nil {
			log.Printf("Failed to copy vars.go: %v\n", err)
		} else {
			fmt.Printf("  Copied: vars.go -> %s\n", varsDst)
			os.Remove(varsSrc)
		}
	}
}
func gen_models_base() {
	p, _ := os.Getwd()
	p = p + `\gen_proto_api_param_models\models`
	genFiles := []string{}
	filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			genFiles = append(genFiles, info.Name())
		}
		return nil
	})
	destPath, _ := os.Getwd()
	// 使用项目目录，而不是父目录
	destDir := filepath.Join(destPath, "models_base")
	// 确保目标目录存在
	if err := os.MkdirAll(destDir, 0755); err != nil {
		log.Printf("Failed to create directory %s: %v\n", destDir, err)
		return
	}
	fmt.Printf("Copying .go files to %s\n", destDir)
	for _, f := range genFiles {
		src := filepath.Join(p, f)
		dst := filepath.Join(destDir, f)
		if err := copyFile(src, dst); err != nil {
			log.Printf("Failed to copy %s to %s: %v\n", src, dst, err)
			continue
		}
		fmt.Printf("  Copied: %s -> %s\n", f, dst)
		os.Remove(src)
	}
}
func gen_advancedTomodels() {
	p, _ := os.Getwd()
	// 使用项目目录，而不是父目录
	advanced := filepath.Join(p, "models_advanced")
	base := filepath.Join(p, "models_base")
	mdls := filepath.Join(p, "models")
	
	// 确保目标目录存在
	if err := os.MkdirAll(mdls, 0755); err != nil {
		log.Printf("Failed to create directory %s: %v\n", mdls, err)
		return
	}
	
	// 1. 先复制 models_base 中的基础模型文件
	fmt.Printf("Copying files from %s to %s\n", base, mdls)
	baseFiles := []string{}
	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			baseFiles = append(baseFiles, info.Name())
		}
		return nil
	})
	for _, f := range baseFiles {
		src := filepath.Join(base, f)
		dst := filepath.Join(mdls, f)
		if err := copyFile(src, dst); err != nil {
			log.Printf("Failed to copy %s to %s: %v\n", src, dst, err)
			continue
		}
		fmt.Printf("  Copied from base: %s -> %s\n", f, dst)
	}
	
	// 2. 再复制 models_advanced 中的高级文件（_gen.go 和 vars.go）
	fmt.Printf("Copying files from %s to %s\n", advanced, mdls)
	advancedFiles := []string{}
	filepath.Walk(advanced, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			advancedFiles = append(advancedFiles, info.Name())
		}
		return nil
	})
	for _, f := range advancedFiles {
		src := filepath.Join(advanced, f)
		dst := filepath.Join(mdls, f)
		if err := copyFile(src, dst); err != nil {
			log.Printf("Failed to copy %s to %s: %v\n", src, dst, err)
			continue
		}
		fmt.Printf("  Copied from advanced: %s -> %s\n", f, dst)
		// 不删除源文件，保留在 models_advanced 中
	}
}
func copyFile(src, dst string) error {
	// 打开源文件
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// 创建目标文件
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	// 复制内容
	_, err = io.Copy(destination, source)
	return err
}
func getParentDir(path string) string {
	// 获取绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return ""
	}

	// 获取父目录
	parent := filepath.Dir(absPath)

	// 如果父目录和当前目录相同（例如根目录），返回错误或空字符串
	if parent == absPath {
		return ""
	}

	return parent
}
