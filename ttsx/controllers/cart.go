package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"github.com/astaxie/beego/orm"
	"ttsx/models"
	"time"
	"strings"
)

type CartController struct {
	beego.Controller
}

//添加购物车
func (this *CartController) HandleCart() {
	skuid, err1 := this.GetInt("skuid")
	count, err2 := this.GetInt("count")
	userName := GetUserName(&this.Controller)
	user, err := GetUser(userName.(string))
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	if err1 != nil || err2 != nil {
		resp["code"] = 1
		resp["msg"] = "发送数据错误"
		this.Data["json"] = resp
		return
	}
	if err != nil {
		resp["code"] = 2
		resp["msg"] = "链接错误"
		this.Data["json"] = resp
		return
	}
	conn, err := redis.Dial("tcp", "192.168.110.111:6379")
	if err != nil {
		resp["code"] = 3
		resp["msg"] = "redis链接错误"
		this.Data["json"] = resp
		return
	}
	defer conn.Close()
	preCount, _ := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), skuid))
	conn.Do("hset", "cart_"+strconv.Itoa(user.Id), skuid, count+preCount)
	resp["code"] = 5
	resp["msg"] = "OK"
	resp["cartCount"] = preCount + count
	this.Data["json"] = resp
}

//添加商品
func (this *CartController) HandleUpdateCart() {
	skuid, err1 := this.GetInt("skuid")
	count, err2 := this.GetInt("count")
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	if err1 != nil || err2 != nil {
		resp["code"] = 1
		resp["msg"] = "数据传输错误"
		this.Data["json"] = resp
		beego.Error(err1, err2)
		return
	}
	usnername := GetUserName(&this.Controller)
	user, err := GetUser(usnername.(string))
	if err != nil {
		resp["code"] = 2
		resp["msg"] = "用户未登录"
		this.Data["json"] = resp
		return
	}
	conn, err := redis.Dial("tcp", "192.168.110.111:6379")
	if err != nil {
		resp["code"] = 3
		resp["msg"] = "redis数据库链接错误"
		this.Data["json"] = resp
		return
	}
	defer conn.Close()
	conn.Do("hset", "cart_"+strconv.Itoa(user.Id), skuid, count)
	conn.Do("hlen", "cart_"+strconv.Itoa(user.Id))
	resp["code"] = 5
	resp["msg"] = "Ok"
	this.Data["json"] = resp
}

//删除
func (this *CartController) HandleDeleteCart() {
	skuid, err := this.GetInt("skuid")
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	if err != nil {
		resp["code"] = 1
		resp["msg"] = "数据传输失败"
		this.Data["json"] = resp
	}
	username := GetUserName(&this.Controller)
	user, err := GetUser(username.(string))
	if err != nil {
		resp["code"] = 2
		resp["msg"] = "未登录"
		this.Data["json"] = resp
	}
	conn, err := redis.Dial("tcp", "192.168.110.111:6379")
	if err != nil {
		resp["code"] = 3
		resp["msg"] = "redis链接失败"
		this.Data["json"] = resp
	}
	conn.Do("hdel", "cart_"+strconv.Itoa(user.Id), skuid)

	resp["code"] = 5
	resp["msg"] = "OK"
	this.Data["json"] = resp

}

//提交订单
func (this *CartController) HandleAddOrder() {
	addrid, err1 := this.GetInt("addrid")
	payid, err2 := this.GetInt("payid")
	skuids := this.GetStrings("skuids")
	totalCount, _ := this.GetInt("totalCount")
	transitprice, _ := this.GetInt("transitprice")
	totalprice, _ := this.GetInt("totalprice")
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	if err1 != nil || err2 != nil {
		resp["code"] = 1
		resp["msg"] = "链接失败"
		this.Data["json"] = resp
		return
	}
	username := GetUserName(&this.Controller)
	user, err := GetUser(username.(string))
	if err != nil {
		resp["code"] = 2
		resp["msg"] = "未登录"
		this.Data["json"] = resp
		return
	}
	o := orm.NewOrm()
	o.Begin()
	var addr models.Address
	addr.Id = addrid
	o.Read(&addr)
	var order models.OrderInfo
	order.User = &user
	order.OrderId = time.Now().Format("20060102150405") + strconv.Itoa(user.Id)
	order.Orderstatus = 1
	order.PayMethod = payid
	order.TotalCount = totalCount
	order.TotalPrice = totalprice
	order.TransitPrice = transitprice
	order.Address = &addr
	o.Insert(&order)

	conn, _ := redis.Dial("tcp", "192.168.110.111:6379")
	defer conn.Close()

	skuids = strings.Fields(skuids[0][1 : len(skuids[0])-1])
	for _, skuid := range skuids {
		id, _ := strconv.Atoi(skuid)
		var goodsSku models.GoodsSKU
		goodsSku.Id = id
		i := 3
		for i > 0 {

			o.Read(&goodsSku)
			count, _ := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), id))
			var ordergoods models.OrderGoods
			if count > goodsSku.Stock {
				resp["code"] = 2
				resp["msg"] = "库存不足"
				this.Data["json"] = resp
				o.Rollback()
				return
			}
			ordergoods.Price = goodsSku.Price * count
			ordergoods.Count = count
			ordergoods.GoodsSKU = &goodsSku
			ordergoods.OrderInfo = &order
			o.Insert(&ordergoods)
			preCount := goodsSku.Stock
			goodsSku.Stock -= count
			goodsSku.Sales += count
			updatecount, _ := o.QueryTable("GoodsSKU").Filter("Id", goodsSku.Id).Filter("Stock", preCount).Update(orm.Params{"Stock": goodsSku.Stock, "Sales": goodsSku.Sales})
			if updatecount == 0 {
				if i > 0 {
					i -= 1
					continue
				}
				resp["code"] = 3
				resp["msg"] = "库存变化，订单提交失败"
				this.Data["json"] = resp
				o.Rollback()
				return
			} else {
				conn.Do("hdel", "cart_"+strconv.Itoa(user.Id), goodsSku.Id)
				break
			}
		}
	}
	o.Commit()
	resp["code"] = 5
	resp["msg"] = "OK"
	this.Data["json"] = resp
}
