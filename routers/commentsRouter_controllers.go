package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["gisa/controllers:AuthController"] = append(beego.GlobalControllerRouter["gisa/controllers:AuthController"],
        beego.ControllerComments{
            Method: "MenuList",
            Router: "/auth/menu/list",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:ErrorController"] = append(beego.GlobalControllerRouter["gisa/controllers:ErrorController"],
        beego.ControllerComments{
            Method: "Error",
            Router: "/error",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("msg"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:ErrorController"] = append(beego.GlobalControllerRouter["gisa/controllers:ErrorController"],
        beego.ControllerComments{
            Method: "Error404",
            Router: "/error/404",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:HomeController"] = append(beego.GlobalControllerRouter["gisa/controllers:HomeController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:HomeController"] = append(beego.GlobalControllerRouter["gisa/controllers:HomeController"],
        beego.ControllerComments{
            Method: "Dashboard",
            Router: "/dashboard",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:HomeController"] = append(beego.GlobalControllerRouter["gisa/controllers:HomeController"],
        beego.ControllerComments{
            Method: "Login",
            Router: "/login",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:HomeController"] = append(beego.GlobalControllerRouter["gisa/controllers:HomeController"],
        beego.ControllerComments{
            Method: "DoLogin",
            Router: "/login",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:HomeController"] = append(beego.GlobalControllerRouter["gisa/controllers:HomeController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: "/logout",
            AllowHTTPMethods: []string{"*"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:JobController"] = append(beego.GlobalControllerRouter["gisa/controllers:JobController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/job",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:JobController"] = append(beego.GlobalControllerRouter["gisa/controllers:JobController"],
        beego.ControllerComments{
            Method: "Save",
            Router: "/job/add",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:LogController"] = append(beego.GlobalControllerRouter["gisa/controllers:LogController"],
        beego.ControllerComments{
            Method: "List",
            Router: "/log/list",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:TestController"] = append(beego.GlobalControllerRouter["gisa/controllers:TestController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/test",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers:UserController"] = append(beego.GlobalControllerRouter["gisa/controllers:UserController"],
        beego.ControllerComments{
            Method: "UserList",
            Router: "/admin/user",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
