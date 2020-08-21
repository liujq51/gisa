package controllers

import (
	"fmt"
	"gisa/common"
	"gisa/logic"
	"gisa/models"
	"html/template"
	"strings"

	"github.com/astaxie/beego"
)

const MESSAGE_TYPE_SUCCESS = "success"
const MESSAGE_TYPE_ERROR = "error"

type BaseController struct {
	beego.Controller
	controllerName     string          //当前控制名称
	actionName         string          //当前action名称
	curUser            models.User     //当前用户信息
	globalUrlWriteList map[string]bool //全局url白名单，所有用户都有权限访问的 controller.action
	loginUrlWriteList  map[string]bool //登录用户url访问白名单，登录用户都可以访问controller.action
	common.Breadcrumbs
}
type JsonData struct {
	Code    int
	Message string
	Data    interface{}
}

func init() {
	//libs.ExceptMethodAppend("ShowHtml", "ShowJSON")
	beego.AddFuncMap("unixTimeFormat", common.UnixTimeFormat)
	beego.AddFuncMap("pagination", common.PaginationRender)
}
func (c *BaseController) Prepare() {
	fmt.Println("prepare 2:")
	//附值
	c.controllerName, c.actionName = c.GetControllerAndAction()
	c.controllerName = strings.ToLower(c.controllerName[0 : len(c.controllerName)-10])
	c.actionName = strings.ToLower(c.actionName)
	//从Session里获取数据 设置用户信息
	c.adapterUserInfo()
	//c.Layout = "layouts/common/layout.html"
	//c.LayoutSections = make(map[string]string)
	c.globalUrlWriteList = map[string]bool{
		"home.login":   true,
		"home.dologin": true,
	}
	c.loginUrlWriteList = map[string]bool{
		"home.dashboard": true,
	}
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	conAct := c.controllerName + "." + c.actionName
	//fmt.Println("router:", c.Ctx.Request.Method, c.Ctx.Request.URL.Path, c.Ctx.Input.URL(), conAct)
	if !c.globalUrlWriteList[conAct] {
		c.checkAuthor()
	}
	c.AddBreadcrumbs("首页", "/")
	//pjax提交时，设置布局为空
	if c.Ctx.Request.Header.Get("X-PJAX") == "true" {
		c.Layout = ""
	} else {
		c.Layout = "layouts/home/layout.html"
	}
	//xsrf
	c.Data["xsrf_token"] = c.XSRFToken()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlContentHeader"] = "layouts/home/content_header.html"
	c.LayoutSections["HtmlNav"] = "layouts/home/nav.html"
	c.LayoutSections["HtmlSidebar"] = "layouts/home/sidebar.html"
	c.LayoutSections["HtmlFooter"] = "layouts/home/footer.html"
}

//从session里取用户信息
func (c *BaseController) adapterUserInfo() {
	a := c.GetSession("user")
	if a != nil {
		c.curUser = a.(models.User)
		c.Data["user"] = a
		//从session获取用户信息
		var menuTree interface{}
		if menuTree = c.GetSession("menuTree"); menuTree == nil {
			menuTree = models.GetGlobalMenu(a).List
			c.SetSession("menuTree", menuTree)
		}
		httpPath := c.Ctx.Request.URL.Path
		menuTreeStr := models.GetMenuTreeHtml(menuTree, httpPath)
		c.Data["menuTree"] = menuTreeStr
	}
}

func (c *BaseController) jsonResult(code logic.JsonResultCode, msg string, obj interface{}) {
	r := &models.JsonResult{code, msg, obj}
	c.Data["json"] = r
	c.ServeJSON()
	c.StopRun()
}

// checkLogin判断用户是否登录，未登录则跳转至登录页面
// 一定要在BaseController.Prepare()后执行
func (c *BaseController) checkLogin() {
	if c.curUser.Id == 0 {
		//登录页面地址
		urlstr := c.URLFor("HomeController.Login") + "?url="
		//登录成功后返回的址为当前
		returnURL := c.Ctx.Request.URL.Path
		//如果ajax请求则返回相应的错码和跳转的地址
		if c.Ctx.Input.IsAjax() {
			//由于是ajax请求，因此地址是header里的Referer
			returnURL := c.Ctx.Input.Refer()
			c.jsonResult(logic.JRCode302, "请登录", urlstr+returnURL)
		}
		c.Redirect(urlstr+returnURL, 302)
		c.StopRun()
	}
}

// 判断某 Controller.Action 当前用户是否有权访问
func (c *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	if c.curUser.Id == 0 {
		return false
	}
	//从session获取用户信息
	user := c.GetSession("user")
	//类型断言
	v, ok := user.(models.User)
	if ok {
		//如果是超级管理员，则直接通过
		if v.IsSuper == true {
			return true
		}
		conAct := ctrlName + "." + ActName
		if !c.loginUrlWriteList[conAct] {
			return true
		}
		//遍历用户所负责的资源列表
		for i, _ := range v.ResourceUrlForList {
			urlfor := strings.TrimSpace(v.ResourceUrlForList[i])
			if len(urlfor) == 0 {
				continue
			}
			// TestController.Get,:last,xie,:first,asta
			strs := strings.Split(urlfor, ",")
			if len(strs) > 0 && strs[0] == (ctrlName+"."+ActName) {
				return true
			}
		}
	}
	return false
}

// 判断用户是否有权访问某地址，无权则会跳转到错误页面
func (c *BaseController) checkAuthor() {
	var (
		httpPath   string
		httpMethod string
	)
	//先判断是否登录
	c.checkLogin()

	//从session获取用户信息
	user := c.GetSession("user")
	//类型断言
	v, ok := user.(models.User)
	if ok {
		httpPath = c.Ctx.Request.URL.Path
		httpMethod = c.Ctx.Request.Method
		//fmt.Println("router:", c.Ctx.Request.Method, c.Ctx.Request.URL.Path, c.Ctx.Input.URL(), conAct)

		hasAuthor := v.CheckPermission(httpPath, httpMethod)
		if !hasAuthor {
			//如果没有权限
			fmt.Println(fmt.Sprintf("author control: path=%s.%s userid=%v  无权访问", c.controllerName, c.actionName, c.curUser.Id))
			fmt.Println(c.Ctx.Input.IsAjax())
			if c.Ctx.Input.IsAjax() {
				c.jsonResult(logic.JRCode401, "无权访问", "")
			} else {
				c.Abort("401")
			}
		}
	}

}

//SetUserSession 获取用户信息（包括资源UrlFor）保存至Session
func (c *BaseController) setUserSession(userId int) error {
	user, err := models.UserOne(userId)
	if err != nil {
		return err
	}
	//获取这个用户能获取到的所有资源列表
	//resourceList := models.ResourceTreeGridByUserId(userId, 1000)
	//for _, item := range resourceList {
	//	m.ResourceUrlForList = append(m.ResourceUrlForList, strings.TrimSpace(item.UrlFor))
	//}
	user = user.WithRoles().WithPermissions().WithMenus()
	c.SetSession("user", *user)
	return nil
}

// 重定向去错误页
func (c *BaseController) pageError(msg string) {
	errorurl := c.URLFor("ErrorController.Error") + "/" + msg
	c.Redirect(errorurl, 302)
	c.StopRun()
}

//server Json
func (c *BaseController) ShowJSON(data *JsonData) {
	c.Data["json"] = data
	c.Controller.ServeJSON()
	c.StopRun()
}

// 是否POST提交
func (c *BaseController) IsPost() bool {
	return c.Ctx.Request.Method == "POST"
}

func (c *BaseController) RedirectMessage(url, message, messageType string) {
	c.AddBreadcrumbs("消息提示", "")
	c.Data["redirect_url"] = url
	c.Data["message"] = message
	c.Data["message_type"] = messageType
	c.ShowHtml("layouts/tip.html")
}

//重新定义beego的render
func (c *BaseController) ShowHtml(tpl ...string) {
	if len(tpl) > 0 {
		c.TplName = tpl[0]
	} else {
		c.TplName = c.controllerName + "/" + c.actionName + ".html"
	}

	//c.Data["homeUrl"] = c.homeUrl
	c.Data["breadcrumbs"] = c.ShowBreadcrumbs()
	//c.Data["menus"] = c.ShowMenu(this.controllerName, this.actionName)

	c.Render()
	c.StopRun()
}
