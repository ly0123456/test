package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ttsx/models"
	"encoding/base64"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowRegister() {
	this.TplName = "register1.html"
}
func (this *UserController) Hendlerregister() {
	userName := this.GetString("userName")
	password := this.GetString("password")
	if userName == "" || password == "" {
		beego.Info("数据不完整")
		this.Data["Eorrmsg"] = "数据不完整"
		this.TplName = "register1.html"
		return
	}
	o := orm.NewOrm()
	var User models.User
	User.Name = userName
	err := o.Read(&User, "Name")
	if err == nil {
		beego.Info("用户名已存在")
		this.Data["Eorrmsg"] = "用户名已存在"
		this.TplName = "register1.html"
		return
	}
	User.PassWord=password
	User.Power=1

	o.Insert(&User)
	//this.Ctx.WriteString("注册成功")
	this.Redirect("/login1", 302)
}
func (this *UserController) ShowLogin() {
	username := this.Ctx.GetCookie("userName")
	beego.Info(username)
	data, _ := base64.StdEncoding.DecodeString(username)
	if string(data) == "" {
		this.Data["checked"] = ""
		this.Data["userName"] = ""
	} else {
		this.Data["checked"] = "checked"
		this.Data["userName"] = string(data)
	}
	this.TplName = "login1.html"
}
func (this *UserController) HendleLogin() {
	//获取数据
	userName := this.GetString("userName")
	password := this.GetString("password")
	//校验数据
	if userName == "" || password == "" {
		beego.Info("数据不完整")
		this.Data["Eorrmsg"] = "数据不完整"
		this.TplName = "login1.html"
		return
	}
	//处理数据
	o := orm.NewOrm()
	var User models.User
	User.Name = userName
	err := o.Read(&User, "Name")
	if err != nil {
		beego.Info("用户名不存在")
		this.Data["Eorrmsg"] = "用户名不存在"
		this.TplName = "login1.html"
		return
	}
	if User.PassWord != password {
		beego.Info("密码错误")
		this.Data["Eorrmsg"] = "密码错误"
		this.TplName = "login1.html"
		return
	}
	if User.Power!=1 {
		this.Data["Eorrmsg"] = "权限不够"
		this.TplName = "login1.html"
		return
	}
	data := this.GetString("remember")
	beego.Info(data, userName)
	if data == "on" {
		temp := base64.StdEncoding.EncodeToString([]byte(userName))//把中文加密
		this.Ctx.SetCookie("userName", temp, 2000)           
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}
	this.SetSession("userName", userName)
	this.Redirect("/goods/comeindex", 302)

}
//退出清除session
func (this *UserController) Showlogout() {
	this.DelSession("userName")
	this.Redirect("/login1", 302)
}
