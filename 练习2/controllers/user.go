package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"练习2/models"
)

type UserController struct {
	beego.Controller
}

func (this  *UserController)ShowRegister() {
	this.TplName="register.html"
}
func (this *UserController)HandleRegister()  {
	userName:=this.GetString("userName")
	password:=this.GetString("password")
	if userName==""||password=="" {
		beego.Info("数据不完整，请重新输入")
		this.Data["errmsg"]="数据不完整，请重新输入"
		this.TplName="rehister.html"
		return
	}
	o:=orm.NewOrm()
	var User models.User
	User.Name=userName
	err:=o.Read(&User,"Name")
	if err==nil {
	beego.Info("用户名已存在")
	this.Data["errmsg"]="用户名已存在"
	this.TplName="register.html"
	return
	}
	User.Password=password
	o.Insert(&User)
	this.Redirect("/login",302)

}
func (this *UserController) ShowLogin()  {
	userName:=this.Ctx.GetCookie("userName")
	//beego.Info(userName)
	if userName=="" {
		this.Data["checked"]=""
		this.Data["userName"]=""
	}else {
		this.Data["checked"]="checked"
		this.Data["userName"]=userName
	}

	this.TplName="login.html"
}
func (this UserController)HandleLogin(){
	userName:=this.GetString("userName")
	password:=this.GetString("password")
	if userName==""||password=="" {
		beego.Info("数据不完整，请重新输入")
		this.Data["errmsg"]="数据不完整，请重新输入"
		this.TplName="login.html"
		return
	}
	o:=orm.NewOrm()
	var User models.User
	User.Name=userName
	err:=o.Read(&User,"Name")
	if err!=nil {
		beego.Info("用户名不存在，请重新输入")
		this.Data["errmsg"]="用户名不存在，请重新输入"
		this.TplName="login.html"
		return
	}
	if password!=User.Password {
		beego.Info("用户名密码错误，请重新输入")
		this.Data["errmsg"]="用户名密码错误，请重新输入"
		this.TplName="login.html"
		return
	}
	remember:=this.GetString("remember")
	//beego.Info(remember)
	if remember =="on"{
		this.Ctx.SetCookie("userName",userName,1000)
	}else {
		this.Ctx.SetCookie("userName",userName,-1)
	}


	this.Redirect("/comeindex",302)
}
