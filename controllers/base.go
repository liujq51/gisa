package controllers

import (
	"fmt"
	"gisa/common"
	"gisa/logic"
	"gisa/models"
	"html/template"
	"strings"

	"github.com/beego/i18n"

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
	Lang string
}
type JsonData struct {
	Code    int
	Message string
	Data    interface{}
}
type langType struct {
	Lang string
	Name string
}

var (
	langTypes []*langType
)

func init() {
	//libs.ExceptMethodAppend("ShowHtml", "ShowJSON")
	beego.AddFuncMap("unixTimeFormat", common.UnixTimeFormat)
	beego.AddFuncMap("pagination", common.PaginationRender)
	beego.AddFuncMap("i18n", common.T)

	// Initialized language type list.
	langs := strings.Split(beego.AppConfig.String("lang::types"), "|")
	names := strings.Split(beego.AppConfig.String("lang::names"), "|")
	langTypes = make([]*langType, 0, len(langs))
	for i, v := range langs {
		langTypes = append(langTypes, &langType{
			Lang: v,
			Name: names[i],
		})
	}

	for _, lang := range langs {
		beego.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/l18n/"+lang+".ini"); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			return
		}
	}
}
func (this *BaseController) Prepare() {
	//l18n
	this.setLangVer()
	//附值
	this.controllerName, this.actionName = this.GetControllerAndAction()
	this.controllerName = strings.ToLower(this.controllerName[0 : len(this.controllerName)-10])
	this.actionName = strings.ToLower(this.actionName)
	//从Session里获取数据 设置用户信息
	this.adapterUserInfo()
	this.globalUrlWriteList = map[string]bool{
		"home.login":   true,
		"home.dologin": true,
	}
	this.loginUrlWriteList = map[string]bool{
		"home.dashboard": true,
	}
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	conAct := this.controllerName + "." + this.actionName
	if !this.globalUrlWriteList[conAct] {
		this.checkAuthor()
	}
	this.AddBreadcrumbs("首页", "/")
	//pjax提交时，设置布局为空
	if this.Ctx.Request.Header.Get("X-PJAX") == "true" {
		this.Layout = ""
	} else {
		this.Layout = "layouts/home/layout.html"
	}
	//xsrf
	this.Data["xsrf_token"] = this.XSRFToken()
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlContentHeader"] = "layouts/home/content_header.html"
	this.LayoutSections["HtmlNav"] = "layouts/home/nav.html"
	this.LayoutSections["HtmlSidebar"] = "layouts/home/sidebar.html"
	this.LayoutSections["HtmlFooter"] = "layouts/home/footer.html"
}

// setLangVer sets site language version.
func (this *BaseController) setLangVer() bool {
	isNeedRedir := false
	hasCookie := false

	// 1. Check URL arguments.
	lang := this.Input().Get("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = this.Ctx.GetCookie("lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir = false
		hasCookie = false
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := this.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. Default language is English.
	if len(lang) == 0 {
		lang = "zh"
		isNeedRedir = false
	}

	curLang := langType{
		Lang: lang,
	}

	// Save language information in cookies.
	if !hasCookie {
		this.Ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")
	}

	restLangs := make([]*langType, 0, len(langTypes)-1)
	for _, v := range langTypes {
		if lang != v.Lang {
			restLangs = append(restLangs, v)
		} else {
			curLang.Name = v.Name
		}
	}

	// Set language properties.
	this.Lang = lang
	this.Data["Lang"] = curLang.Lang
	this.Data["CurLang"] = curLang.Name
	this.Data["RestLangs"] = restLangs

	return isNeedRedir
}

//从session里取用户信息
func (this *BaseController) adapterUserInfo() {
	a := this.GetSession("user")
	if a != nil {
		this.curUser = a.(models.User)
		this.Data["user"] = a
		//从session获取用户信息
		var menuTree interface{}
		if menuTree = this.GetSession("menuTree"); menuTree == nil {
			menuTree = models.GetGlobalMenu(a).List
			this.SetSession("menuTree", menuTree)
		}
		httpPath := this.Ctx.Request.URL.Path
		menuTreeStr := models.GetMenuTreeHtml(menuTree, httpPath, this.Lang)
		this.Data["menuTree"] = menuTreeStr
	}
}

func (this *BaseController) jsonResult(code logic.JsonResultCode, msg string, obj interface{}) {
	r := &models.JsonResult{code, msg, obj}
	this.Data["json"] = r
	this.ServeJSON()
	this.StopRun()
}

// checkLogin判断用户是否登录，未登录则跳转至登录页面
// 一定要在BaseController.Prepare()后执行
func (this *BaseController) checkLogin() {
	if this.curUser.Id == 0 {
		//登录页面地址
		urlstr := this.URLFor("HomeController.Login") + "?url="
		//登录成功后返回的址为当前
		returnURL := this.Ctx.Request.URL.Path
		//如果ajax请求则返回相应的错码和跳转的地址
		if this.Ctx.Input.IsAjax() {
			//由于是ajax请求，因此地址是header里的Referer
			returnURL := this.Ctx.Input.Refer()
			this.jsonResult(logic.JRCode302, "请登录", urlstr+returnURL)
		}
		this.Redirect(urlstr+returnURL, 302)
		this.StopRun()
	}
}

// 判断某 Controller.Action 当前用户是否有权访问
func (this *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	if this.curUser.Id == 0 {
		return false
	}
	//从session获取用户信息
	user := this.GetSession("user")
	//类型断言
	v, ok := user.(models.User)
	if ok {
		//如果是超级管理员，则直接通过
		if v.IsSuper == true {
			return true
		}
		conAct := ctrlName + "." + ActName
		if !this.loginUrlWriteList[conAct] {
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
func (this *BaseController) checkAuthor() {
	var (
		httpPath   string
		httpMethod string
	)
	//先判断是否登录
	this.checkLogin()

	//从session获取用户信息
	user := this.GetSession("user")
	//类型断言
	v, ok := user.(models.User)
	if ok {
		httpPath = this.Ctx.Request.URL.Path
		httpMethod = this.Ctx.Request.Method
		//fmt.Println("router:", this.Ctx.Request.Method, this.Ctx.Request.URL.Path, this.Ctx.Input.URL(), conAct)

		hasAuthor := v.CheckPermission(httpPath, httpMethod)
		if !hasAuthor {
			//如果没有权限
			fmt.Println(fmt.Sprintf("author control: path=%s.%s userid=%v  无权访问", this.controllerName, this.actionName, this.curUser.Id))
			fmt.Println(this.Ctx.Input.IsAjax())
			if this.Ctx.Input.IsAjax() {
				this.jsonResult(logic.JRCode401, "无权访问", "")
			} else {
				this.Abort("401")
			}
		}
	}

}

//SetUserSession 获取用户信息（包括资源UrlFor）保存至Session
func (this *BaseController) setUserSession(userId int) error {
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
	this.SetSession("user", *user)
	return nil
}

// 重定向去错误页
func (this *BaseController) pageError(msg string) {
	errorurl := this.URLFor("ErrorController.Error") + "/" + msg
	this.Redirect(errorurl, 302)
	this.StopRun()
}

//server Json
func (this *BaseController) ShowJSON(data *JsonData) {
	this.Data["json"] = data
	this.Controller.ServeJSON()
	this.StopRun()
}

// 是否POST提交
func (this *BaseController) IsPost() bool {
	return this.Ctx.Request.Method == "POST"
}

func (this *BaseController) RedirectMessage(url, message, messageType string) {
	this.AddBreadcrumbs("消息提示", "")
	this.Data["redirect_url"] = url
	this.Data["message"] = message
	this.Data["message_type"] = messageType
	this.ShowHtml("layouts/tip.html")
}

//重新定义beego的render
func (this *BaseController) ShowHtml(tpl ...string) {
	if len(tpl) > 0 {
		this.TplName = tpl[0]
	} else {
		this.TplName = this.controllerName + "/" + this.actionName + ".html"
	}

	//this.Data["homeUrl"] = this.homeUrl
	this.Data["breadcrumbs"] = this.ShowBreadcrumbs()
	//this.Data["menus"] = this.ShowMenu(this.controllerName, this.actionName)

	this.Render()
	this.StopRun()
}
