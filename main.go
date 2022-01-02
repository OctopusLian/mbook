/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-01-02 10:34:14
 * @LastEditors: neozhang
 * @LastEditTime: 2022-01-02 10:35:47
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
