/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-01-01 10:56:14
 * @LastEditors: neozhang
 * @LastEditTime: 2022-01-01 10:58:09
 */
package main

import (
	_ "mbook/routers"
	_ "mbook/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
