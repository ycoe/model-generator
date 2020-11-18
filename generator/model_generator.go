package generator

import (
	"fmt"
	"github.com/bigkucha/model-generator/helper"
	"github.com/dave/jennifer/jen"
	"github.com/jinzhu/inflection"
	"os"
	"strings"
)

func GenerateModel(tableName string, columns []map[string]string, dir string) {
	index := strings.LastIndex(dir, "/")
	daoPackage := dir[index+1 : len(dir)]
	var codes []jen.Code
	for _, col := range columns {
		t := col["Type"]
		column := col["Field"]
		var st *jen.Statement
		if column == "id" {
			st = jen.Id("ID").Uint().Tag(map[string]string{"json": "id"})
		} else {
			st = jen.Id(helper.SnakeCase2CamelCase(column, true))
			getCol(st, t)
			st.Tag(map[string]string{"json": column})
		}
		codes = append(codes, st)
	}
	f := jen.NewFile(daoPackage)
	f.HeaderComment("Code generated by model-generator. DO NOT EDIT.")
	f.ImportAlias("time", "time")
	f.Type().Id(helper.SnakeCase2CamelCase(inflection.Singular(tableName), true)).Struct(codes...)
	_ = os.MkdirAll(dir, os.ModePerm)
	fileName := dir + "/" + inflection.Singular(tableName) + ".go"
	fmt.Println(fileName)
	if err := f.Save(fileName); err != nil {
		fmt.Println(err.Error())
	}
}

func getCol(st *jen.Statement, t string) {
	prefix := strings.Split(t, "(")[0]
	switch prefix {
	case "int", "tinyint", "smallint", "bigint", "mediumint":
		st.Int()
	case "float":
		st.Float32()
	case "decimal":
		st.Float32()
	case "date", "time", "timestamp", "year", "datetime":
		st.Id("*").Qual("time", "Time")
	default:
		st.String()
	}
}
