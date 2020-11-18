package generator

import (
	"fmt"
	"github.com/bigkucha/model-generator/helper"
	"github.com/dave/jennifer/jen"
	"github.com/jinzhu/inflection"
	"os"
)

func GenBaseDao(appId, packageName string, tableNames []string) {
	f := jen.NewFile(packageName)
	genDaoStruct(f)
	genVarDefaultDao(f)
	genGetDao(f)
	genInit(appId, f)
	genNewDao(f, tableNames)
	genPing(f)
	genDisconnect(f)

	filename := "./" + packageName + "/dao.go"
	fmt.Println(filename)
	_ = os.MkdirAll("./"+packageName, os.ModePerm)
	if err := f.Save(filename); err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("%#v\n", f)
}

func genDisconnect(f *jen.File) {
	/**
	func (d *Dao) Disconnect() error {
		return d.client.DB().Close()
	}
	*/
	f.Func().Params(
		jen.Id("d *Dao"),
	).Id("Disconnect").Params().Id("error").Block(
		jen.Return(
			jen.Id("d").Dot("client").Dot("DB").Call().Dot("Close").Call(),
		),
	)
}

func genPing(f *jen.File) {
	/**
	func (d *Dao) Ping() error {
		return d.client.DB().Ping()
	}
	*/
	f.Func().Params(
		jen.Id("d *Dao"),
	).Id("Ping").Params().Id("error").Block(
		jen.Return(
			jen.Id("d").Dot("client").Dot("DB").Call().Dot("Ping").Call(),
		),
	)
}

/**
// newDao 创建 Dao 实例
func newDao(c *conf.Config) (*Dao, error) {
	var (
		d   Dao
		err error
	)
	if d.client, err = gorm.Open(c.DB.DriverName, c.DB.URL); err != nil {
		return nil, err
	}
	d.client.SingularTable(true)       //表名采用单数形式
	d.client.DB().SetMaxOpenConns(100) //SetMaxOpenConns用于设置最大打开的连接数
	d.client.DB().SetMaxIdleConns(10)  //SetMaxIdleConns用于设置闲置的连接数
	//d.client.LogMode(true)

	if err = d.client.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.FinanceAccount{},
	).Error; err != nil {
		_ = d.client.Close()
		return nil, err
	}

	return &d, nil
}
*/
func genNewDao(f *jen.File, tableNames []string) {
	var codes []jen.Code
	for _, tableName := range tableNames {
		entityName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true)
		codes = append(codes, jen.Line().Id("&").Qual("finance/model", entityName).Block())
	}

	f.Func().Id("newDao").Params(
		jen.Id("c *conf.Config"),
	).Params(
		jen.Id("*Dao"),
		jen.Id("error"),
	).Block(
		jen.Var().Id("d Dao").Line().Var().Id("err error"),
		jen.If(
			jen.Id("d.client, err").Op("=").Id("gorm").Dot("Open").Call(
				jen.Id("c.DB.DriverName"),
				jen.Id("c.DB.URL"),
			),
			jen.Id("err").Op("!=").Nil(),
		).Block(
			jen.Return(
				jen.Nil(),
				jen.Id("err"),
			),
		).Line(),
		jen.Id("d").Dot("client").Dot("SingularTable").Call(
			jen.True(),
		).Comment("表名采用单数形式"),
		jen.Id("d").Dot("client").Dot("DB").Call().Dot("SetMaxOpenConns").Call(
			jen.Id("100"),
		).Comment("SetMaxOpenConns用于设置最大打开的连接数"),
		jen.Id("d").Dot("client").Dot("DB").Call().Dot("SetMaxIdleConns").Call(
			jen.Id("10"),
		).Comment("SetMaxIdleConns用于设置闲置的连接数"),
		jen.Comment("d.client.LogMode(true)").Line(),

		jen.If(
			jen.Id("err").Op("=").Id("d").Dot("client").Dot("Set").Call(
				jen.Lit("gorm:table_options"),
				jen.Lit("ENGINE=InnoDB"),
			).Dot("AutoMigrate").Call(codes...).Dot("Error"),
			jen.Id("err").Op("!=").Nil(),
		).Block(
			jen.Id("_").Op("=").Id("d").Dot("client").Dot("Close").Call(),
			jen.Return(
				jen.Nil(),
				jen.Id("err"),
			),
		),
		jen.Return(
			jen.Id("&d"),
			jen.Nil(),
		),
	)
}

/**
func Init(c *conf.Config) (err error) {
	defaultDao, err = newDao(c)
	return
}
*/
func genInit(appId string, f *jen.File) {
	f.Line().Func().Id("Init").Params(
		jen.Id("c *").Qual(appId+"/conf", "Config"),
	).Params(
		jen.Id("err").Id("error"),
	).Block(
		jen.Id("defaultDao, err").Op("=").Id("newDao").Call(
			jen.Id("c"),
		).Line().Return(),
	)
}

/**
func GetDao() *Dao {
	return defaultDao
}
*/
func genGetDao(f *jen.File) {
	f.Line().Func().Id("GetDao").Params().Params(
		jen.Id("*Dao"),
	).Block(
		jen.Return(
			jen.Id("defaultDao"),
		),
	)
}

func genVarDefaultDao(f *jen.File) *jen.Statement {
	return f.Var().Id("defaultDao").Id("*Dao")
}

func genDaoStruct(f *jen.File) *jen.Statement {
	return f.Type().Id("Dao").Struct(
		jen.Id("client").Id("*").Qual("github.com/jinzhu/gorm", "DB"),
	)
}
