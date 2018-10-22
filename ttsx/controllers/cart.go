package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type CartController struct {
	beego.Controller
}
//添加购物车
func (this *CartController)HandleCart(){
	skuid,err1:=this.GetInt("skuid")
	count,err2:=this.GetInt("count")
	userName:=GetUserName(&this.Controller)
	user,err:=GetUser(userName.(string))
	resp:=make(map[string]interface{})
	defer  this.ServeJSON()
	if err1!=nil||err2!=nil {
		resp["code"]=1
		resp["msg"]="发送数据错误"
		this.Data["json"]=resp
		return
	}
	if err!=nil {
		resp["code"]=2
		resp["msg"]="链接错误"
		this.Data["json"]=resp
		return
	}
	conn,err:=redis.Dial("tcp","192.168.110.111:6379")
	if err!=nil {
		resp["code"]=3
		resp["msg"]="redis链接错误"
		this.Data["json"]=resp
		return
	}
	defer conn.Close()
	preCount,_:=redis.Int(conn.Do("hget","cart_"+strconv.Itoa(user.Id),skuid))
	conn.Do("hset","cart_"+strconv.Itoa(user.Id),skuid,count+preCount)
	resp["code"]=5
	resp["msg"]="OK"
	resp["cartCount"]=preCount+count
	this.Data["json"]=resp
}
//添加商品
func (this *CartController)HandleUpdateCart()  {
	skuid,err1:=this.GetInt("skuid")
	count,err2:=this.GetInt("count")
	resp:=make(map[string]interface{})
	defer this.ServeJSON()
	if err1!=nil||err2!=nil {
		resp["code"]=1
		resp["msg"]="数据传输错误"
		this.Data["json"]=resp
		beego.Error(err1,err2)
		return
	}
	usnername:=GetUserName(&this.Controller)
	user,err:=GetUser(usnername.(string))
	if err!=nil {
		resp["code"]=2
		resp["msg"]="用户未登录"
		this.Data["json"]=resp
		return
	}
	conn,err:=redis.Dial("tcp","192.168.110.111:6379")
	if err!=nil {
		resp["code"]=3
		resp["msg"]="redis数据库链接错误"
		this.Data["json"]=resp
		return
	}
	defer conn.Close()
	conn.Do("hset","cart_"+strconv.Itoa(user.Id),skuid,count)
	conn.Do("hlen","cart_"+strconv.Itoa(user.Id))
	resp["code"]=5
	resp["msg"]="Ok"
	this.Data["json"]=resp







}
