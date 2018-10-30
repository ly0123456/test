package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ttsx/models"
	"encoding/base64"
	"regexp"
	"github.com/astaxie/beego/utils"
	"strconv"
)

type UserControler struct {
	beego.Controller
}
//显示登录信息
func (this *UserControler)ShowReg()  {
	this.TplName="register.html"
}
//处理注册信息
func (this *UserControler)HandleReg()  {
	userName:=this.GetString("user_name")
	pwd:=this.GetString("pwd")
	email:=this.GetString("email")
	cpwd:=this.GetString("cpwd")
	if userName==""||pwd==""||cpwd==""||email=="" {
		this.Data["errmsg"]="数据不完整,请重新输入"
		this.TplName="register.html"
		return
	}
	if cpwd!=pwd {
		this.Data["errmsg"]="两次密码不一致,请重新输入"
		this.TplName="register.html"
		return
	}
	reg,_ :=regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res := reg.FindString(email)
	if res == "" {
		this.Data["errmsg"] = "邮箱格式不正确"
		this.TplName = "register.html"
		return
	}

	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	user.PassWord=pwd
	user.Email=email
	err:=o.Read(&user)
	//判断用户名是否存在
	if err==nil{
		beego.Info("用户名已存在,请重新输入")
		this.Data["errmsg"]="用户名已存在,请重新输入"
		this.TplName="register.html"
		return
	}
	o.Insert(&user)

	emailConfig:=`{"username":"791641402@qq.com","password":"wusujihictgsbahj","host":"smtp.qq.com","post":"587"}`
	emailconn:=utils.NewEMail(emailConfig)
	emailconn.From="791641402@qq.com"
	emailconn.To=[]string{email}
	emailconn.Subject="天天生鲜注册"
	emailconn.Text="192.168.110.111:8080/active?id="+strconv.Itoa(user.Id)
	err=emailconn.Send()
	beego.Info(err)
	if err!=nil {
		this.Data["errmsg"]="发送失败，请重新注册"
		this.TplName="register.html"
		return
	}

	this.Redirect("/login",302)
}
//显示登录页
func (this *UserControler)ShowLog()  {
	userName:=this.Ctx.GetCookie("userName")
	data,err:=base64.StdEncoding.DecodeString(userName)
	if err!=nil {
		beego.Info("解密失败")
		return
	}
	if string(data)=="" {
		this.Data["userName"]=""
		this.Data["checked"]=""
	}else {
		this.Data["userName"]=string(data)
		this.Data["checked"]="checked"
	}
	this.TplName="login.html"

}
//处理登录信息
func (this *UserControler)HandleLog()  {

	username:=this.GetString("username")
	pwd:=this.GetString("pwd")
	o:=orm.NewOrm()
	var user models.User
	user.Name=username
	err:=o.Read(&user,"Name")

	//判断用户名是否存在

	if err!=nil {
		beego.Info("用户名不存在，请重新输入")
		this.Data["errmsg"]="用户名不存在，请重新输入"
		this.TplName="login.html"
		return
	}
	//判断密码
	if user.PassWord!=pwd {
		beego.Info("密码错误，请重新输入")
		this.Data["errmsg"]="密码错误，请重新输入"
		this.TplName="login.html"
		return
	}
	//判断邮箱
	if user.Active!=true {
		beego.Info("邮箱未激活，请激活邮箱")
		this.Data["errmsg"]="邮箱未激活，请激活邮箱"
		this.TplName="login.html"
		return
	}
	data1:=base64.StdEncoding.EncodeToString([]byte(username))

	remember:=this.GetString("remember")
	if remember!="" {
		this.Ctx.SetCookie("userName",data1,1000)
	}else {
		this.Ctx.SetCookie("userName",data1,-1)
	}

	this.SetSession("userName",username)
	this.Redirect("/showindex",302)

}
func (this *UserControler)Handleactive()  {
	id,err:=this.GetInt("id")
	if err==nil {
		this.Data["errmsg"]="激活失败"
		this.TplName="register.html"
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Id=id
	err=o.Read(&user)
	if err!=nil{
		this.Data["errmsg"]="激活失败"
		this.TplName="register.html"
		return
	}
	user.Active=true
	o.Update(&user)
	this.Redirect("/login",302)
}
func (this *UserControler) ShowOut()  {
	this.DelSession("userName")
	this.Redirect("/login",302)
}