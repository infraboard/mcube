package cmd

import (
{{ if $.EnableMySQL -}}
	"context"
	"fmt"
	"io/ioutil"
	"time"
{{- end }}

	"github.com/spf13/cobra"

{{ if $.EnableMySQL -}}
	"{{.PKG}}/conf"
{{- end }}
)

var (
	createTableFilePath string
)

// initCmd represents the start command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "{{.Name}} 服务初始化",
	Long:  "{{.Name}} 服务初始化",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 初始化全局变量
		if err := loadGlobalConfig(confType); err != nil {
			return err
		}

{{ if $.EnableMySQL -}}
		err := createTables()
		if err != nil {
			return err
		}
{{- end }}

		return nil
	},
}

{{ if $.EnableMySQL -}}
func createTables() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	// 读取SQL文件
	sqlFile, err := ioutil.ReadFile(createTableFilePath)
	if err != nil {
		return err
	}

	fmt.Println("执行的SQL: ")
	fmt.Println(string(sqlFile))

	// 执行SQL文件
	_, err = db.ExecContext(ctx, string(sqlFile))
	if err != nil {
		return err
	}

	return nil
}
{{- end }}

func init() {
{{ if $.EnableMySQL -}}
	initCmd.PersistentFlags().StringVarP(&createTableFilePath, "sql-file-path", "s", "docs/schema/tables.sql", "the sql file path")
{{- end }}
	RootCmd.AddCommand(initCmd)
}