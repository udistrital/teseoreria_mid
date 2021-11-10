package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:CuentaBancariaBancoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:CuentaBancariaBancoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:CuentaBancariaBancoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:CuentaBancariaBancoController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:SolicitudAvanceController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:SolicitudAvanceController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:SolicitudAvanceController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:SolicitudAvanceController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:SolicitudAvanceController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:SolicitudAvanceController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:SolicitudAvanceController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:SolicitudAvanceController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:SolicitudAvanceController"] = append(beego.GlobalControllerRouter["github.com/udistrital/tesoreria_mid/controllers:SolicitudAvanceController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
