/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-01-02 10:34:14
 * @LastEditors: neozhang
 * @LastEditTime: 2022-01-02 10:36:13
 */
package sysinit

import (
	_ "mbook/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//调用方式
//dbinit() 或 dbinit("w") 或 dbinit("default") //初始化主库
//dbinit("w","r")	//同时初始化主库和从库
//dbinit("w")
func dbinit(aliases ...string) {
	//如果是开发模式，则显示命令信息
	isDev := (beego.AppConfig.String("runmode") == "dev")

	if len(aliases) > 0 {
		for _, alias := range aliases {
			registDatabase(alias)
			//主库 自动建表
			if "w" == alias {
				orm.RunSyncdb("default", false, isDev)
			}
		}
	} else {
		registDatabase("w")
		orm.RunSyncdb("default", false, isDev)
	}

	if isDev {
		orm.Debug = isDev
	}
}

func registDatabase(alias string) {
	if len(alias) == 0 {
		return
	}
	//连接名称
	dbAlias := alias
	if "w" == alias || "default" == alias {
		dbAlias = "default"
		alias = "w"
	}
	//数据库名称
	dbName := beego.AppConfig.String("db_" + alias + "_database")
	//数据库连接用户名
	dbUser := beego.AppConfig.String("db_" + alias + "_username")
	//数据库连接用户名
	dbPwd := beego.AppConfig.String("db_" + alias + "_password")
	//数据库IP（域名）
	dbHost := beego.AppConfig.String("db_" + alias + "_host")
	//数据库端口
	dbPort := beego.AppConfig.String("db_" + alias + "_port")

	orm.RegisterDataBase(dbAlias, "mysql", dbUser+":"+dbPwd+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8", 30)

}
