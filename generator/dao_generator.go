package generator

import (
	"fmt"
	"github.com/bigkucha/model-generator/helper"
	"github.com/dave/jennifer/jen"
	"github.com/jinzhu/inflection"
	"os"
	"strings"
)

func GenerateDao(appId, tableName, dir string) {
	index := strings.LastIndex(dir, "/")
	daoPackage := dir[index+1 : len(dir)]
	f := jen.NewFile(daoPackage)
	f.HeaderComment("Code generated by model-generator.")
	f.ImportAlias(appId+"/model", "model")
	f.ImportAlias(appId+"/proto", appId)
	f.ImportAlias("time", "time")
	createFun(f, appId, tableName)
	_ = os.MkdirAll(dir, os.ModePerm)
	fileName := dir + "/" + inflection.Singular(tableName) + ".dao.go"
	if err := f.Save(fileName); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(fileName)
}

func createFun(f *jen.File, appId, tableName string) {
	entityName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true)
	f.Func().Params(
		jen.Id("d *Dao"),
	).Id("Create" + entityName).Params(
		jen.Id("entity").Id("*").Qual(appId+"/model", entityName),
	).Params(
		jen.Id("uint"),
		jen.Id("error"),
	).Block(
		jen.Id("result").Op(":=").Id("d").Dot("client").Dot("Table").Call(
			jen.Lit(tableName),
		).Dot("Create").Call(
			jen.Id("entity"),
		),
		jen.Return(
			jen.Id("entity").Dot("ID"),
			jen.Id("result").Dot("Error"),
		),
	)
}