package generator

import (
	"fmt"
	"github.com/bigkucha/model-generator/database"
	"github.com/urfave/cli"
	"strings"
)

func Generate(c *cli.Context) error {
	dbSns := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		c.String("u"), c.String("p"), c.String("d"))
	db := database.GetDB(dbSns)
	appId := c.String("appid")
	daoDir := c.String("daodir")
	tableName := c.String("t")
	if tableName == "ALL" {
		tableNames := make([]string, 0)
		tables := db.GetDataBySql("show tables")
		for _, table := range tables {
			tableName := table["Tables_in_"+c.String("d")]
			tableNames = append(tableNames, tableName)
			columns := db.GetDataBySql("desc " + tableName)
			GenerateModel(tableName, columns, c.String("dir"))
			GenerateDao(appId, tableName, daoDir)
		}

		//生成dao.go
		index := strings.LastIndex(daoDir, "/")
		daoPackage := daoDir[index+1 : len(daoDir)]
		GenBaseDao(appId, daoPackage, tableNames)
	} else {
		columns := db.GetDataBySql("desc " + tableName)
		GenerateModel(tableName, columns, c.String("dir"))
		GenerateDao(appId, tableName, daoDir)
	}
	return nil
}

