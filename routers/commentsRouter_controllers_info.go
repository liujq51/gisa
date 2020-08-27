package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["gisa/controllers/info:MenuController"] = append(beego.GlobalControllerRouter["gisa/controllers/info:MenuController"],
        beego.ControllerComments{
            Method: "Index",
            Router: "/info/menu",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers/info:MenuController"] = append(beego.GlobalControllerRouter["gisa/controllers/info:MenuController"],
        beego.ControllerComments{
            Method: "Add",
            Router: "/info/menu/add",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers/info:MenuController"] = append(beego.GlobalControllerRouter["gisa/controllers/info:MenuController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/info/menu/delete",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers/info:MenuController"] = append(beego.GlobalControllerRouter["gisa/controllers/info:MenuController"],
        beego.ControllerComments{
            Method: "DoUpdate",
            Router: "/info/menu/doupdate",
            AllowHTTPMethods: []string{"post","put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers/info:MenuController"] = append(beego.GlobalControllerRouter["gisa/controllers/info:MenuController"],
        beego.ControllerComments{
            Method: "SaveMenuOrder",
            Router: "/info/menu/order",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers/info:MenuController"] = append(beego.GlobalControllerRouter["gisa/controllers/info:MenuController"],
        beego.ControllerComments{
            Method: "Update",
            Router: "/info/menu/update/:menu_id",
            AllowHTTPMethods: []string{"get","post","put"},
            MethodParams: param.Make(
				param.New("menu_id", param.InPath),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers/info:PermissionController"] = append(beego.GlobalControllerRouter["gisa/controllers/info:PermissionController"],
        beego.ControllerComments{
            Method: "Index",
            Router: "/info/permission",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gisa/controllers/info:RoleController"] = append(beego.GlobalControllerRouter["gisa/controllers/info:RoleController"],
        beego.ControllerComments{
            Method: "Index",
            Router: "/info/menu",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
