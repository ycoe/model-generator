package generator

import (
	"fmt"
	"github.com/bigkucha/model-generator/helper"
	"github.com/dave/jennifer/jen"
	"github.com/jinzhu/inflection"
	"os"
	"strings"
)

/**
package dao

import model "finance/model"

type AccountDao struct {
	dao *Dao
}

func (account *AccountDao) Create (entity *model.Account) (uint, error) {
	result := account.dao.client.Table("accounts").Create(entity)
	return entity.ID, result.Error
}
*/
func GenerateDao(orgTableName string, appId, tableName, dir string) {
	index := strings.LastIndex(dir, "/")
	daoPackage := dir[index+1 : len(dir)]
	f := jen.NewFile(daoPackage)
	f.HeaderComment("Code generated by model-generator.")
	f.ImportAlias(appId+"/model", "model")
	f.ImportAlias(appId+"/proto", appId)
	f.ImportAlias("time", "time")

	genEntityDaoStruct(f, tableName)
	genGetDb(f, orgTableName, tableName)
	genCreateFun(f, appId, tableName)

	_ = os.MkdirAll(dir, os.ModePerm)
	fileName := dir + "/" + inflection.Singular(tableName) + ".dao.go"
	if err := f.Save(fileName); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(fileName)
}

/**
type AccountDao struct {
	Dao *Dao
}
*/
func genEntityDaoStruct(f *jen.File, tableName string) {
	tableEntityDaoName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true) + "Dao"
	f.Type().Id(tableEntityDaoName).Struct(
		jen.Id("Dao").Id("*Dao"),
	)
}

/**
func (dao *AccountDao) getDb(tableName string) *gorm.DB {
	tableDao := Dao.dao.client.Table(tableName)
	return tableDao
}
*/
func genGetDb(f *jen.File, orgTableName string,  tableName string) {
	tableEntityDaoName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true) + "Dao"
	f.Func().Params(
		jen.Id("d").Id("*" + tableEntityDaoName),
	).Id("getDb").Params().Params(
		jen.Id("*").Qual("gorm.io/gorm", "DB"),
	).Block(
		jen.Return(
			jen.Id("d").Dot("Dao").Dot("client").Dot("Table").Call(
				jen.Lit(orgTableName),
			),
		),
	).Line()
}

/**
func (dao *AccountDao) Create(entity *model.Account) (uint, error) {
	result := dao.getDb("account").Create(entity)
	return entity.ID, result.Error
}
*/
func genCreateFun(f *jen.File, appId, tableName string) {
	entityName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true)
	entityDaoName := entityName + "Dao"
	f.Func().Params(
		jen.Id("d").Id("*"+entityDaoName),
	).Id("Create").Params(
		jen.Id("entity").Id("*").Qual(appId+"/model", entityName),
	).Params(
		jen.Id("uint"),
		jen.Id("error"),
	).Block(
		jen.Id("result").Op(":=").Id("d").Dot("getDb").Call().Dot("Create").Call(
			jen.Id("entity"),
		),
		jen.Return(
			jen.Id("entity").Dot("ID"),
			jen.Id("result").Dot("Error"),
		),
	)
}
