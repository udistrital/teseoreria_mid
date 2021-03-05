package main

import (
	"github.com/astaxie/beego"
	_ "github.com/udistrital/tesoreria_mid/routers"
	"github.com/udistrital/utils_oas/customerrorv2"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.ErrorController(&customerrorv2.CustomErrorController{})
	beego.Run()
}
