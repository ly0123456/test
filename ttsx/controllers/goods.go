package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ttsx/models"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"math"
	"github.com/smartwalle/alipay"
	"fmt"

)

type GoodsController struct {
	beego.Controller
}

//获取购物车条目
func cartCount(this *beego.Controller)(int,error){
	username:=GetUserName(this)
	user,err:=GetUser(username.(string))
	conn,_:=redis.Dial("tcp","192.168.110.111:6379")
	count,_:=redis.Int(conn.Do("hlen","cart_"+strconv.Itoa(user.Id)))
	return count,err
}
//主页展示
func (this *GoodsController) ShowIndex() {
	GetUserName(&this.Controller)
	o := orm.NewOrm()
	//类型
	var types []models.GoodsType
	o.QueryTable("GoodsType").All(&types)
	this.Data["types"] = types
	//轮波图
	var lunbotus []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&lunbotus)
	beego.Info(lunbotus)
	//促销商品
	var promotionBanner []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&promotionBanner)
	//
	// beego.Info(promotionBanner)
	goods := make([]map[string]interface{}, len(types))
	for i, value := range types {
		temp := make(map[string]interface{})
		temp["type"] = value
		goods[i] = temp
	}
	for _, vlaue := range goods {
		var textbanner []models.IndexTypeGoodsBanner
		var imagebanner []models.IndexTypeGoodsBanner
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsSKU", "GoodsType").OrderBy("Index").Filter("GoodsType", vlaue["type"]).Filter("DisplayType", 0).All(&textbanner)
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsSKU", "GoodsType").OrderBy("Index").Filter("GoodsType", vlaue["type"]).Filter("DisplayType", 1).All(&imagebanner)
		vlaue["textbanner"] = textbanner
		vlaue["imagebanner"] = imagebanner
	}

	this.Data["goods"] = goods
	this.Data["lunbotus"] = lunbotus
	this.Data["promotionBanner"] = promotionBanner
	count,err:= cartCount(&this.Controller)
	if err!=nil{
		count=0
	}
	this.Data["cartcount"]=count
	this.Data["title"] = "天天生鲜 首页"
	this.Layout = "layout.html"
	this.TplName = "index.html"
}

//用户中心展示
func (this *GoodsController) ShowCenter() {
	userName := GetUserName(&this.Controller)
	o := orm.NewOrm()
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")
	var addr models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Id", user.Id).Filter("Isdefault", true).One(&addr)

	var goods []models.GoodsSKU
	conn, _ := redis.Dial("tcp", "192.168.110.111:6379")
	reply, err := conn.Do("lrange", "history"+strconv.Itoa(user.Id), 0, 4)
	replyInts, _ := redis.Ints(reply, err)
	for _, val := range replyInts {
		var temp models.GoodsSKU
		o.QueryTable("GoodsSKU").Filter("Id", val).One(&temp)
		goods = append(goods, temp)
	}
	this.Data["goods"] = goods
	this.Data["username"] = user.Name
	this.Data["phone"] = addr.Phone
	this.Data["addr"] = addr.Addr
	this.Data["active1"] = "active"
	this.Data["title"] = "用户中心"
	this.Layout = "layout.html"
	this.TplName = "user_center_info.html"
}

//全部订单
func (this *GoodsController) ShowOrder() {
	unsename:=GetUserName(&this.Controller)
	user,_:=GetUser(unsename.(string))
	page,err:=this.GetInt("page")
	if err!=nil {
		page=1
	}
	pageSize:=2
	start:=(page-1)*pageSize
	o:=orm.NewOrm()
	var goodsInfos []models.OrderInfo
count,_:=o.QueryTable("OrderInfo").RelatedSel("User").Filter("User__Id",user.Id).Count()
	o.QueryTable("OrderInfo").RelatedSel("User").Filter("User__Id",user.Id).Limit(pageSize,start).All(&goodsInfos)
	pagecount:=math.Ceil(float64(count)/float64(pageSize))
	goods:=make([]map[string]interface{},len(goodsInfos))
	for i,goodsInfo:=range goodsInfos{
		var orderGoods []models.OrderGoods
		o.QueryTable("OrderGoods").RelatedSel("OrderInfo","GoodsSKU").Filter("OrderInfo__Id",goodsInfo.Id).All(&orderGoods)
		temp:=make(map[string]interface{})
		temp["goodsInfo"]=goodsInfo
		temp["orderGoods"]=orderGoods
		goods[i]=temp
	}
	parpage:=page-1
	if parpage<1 {
		parpage=1
	}
	nextpage:=page+1
	if nextpage>int(pagecount) {
		nextpage=int(pagecount)
	}
	pages:=showpage(page,int(pagecount))
	beego.Info(pages)
	beego.Info(pagecount)
	this.Data["parpage"]=parpage
	this.Data["nextpage"]=nextpage
	this.Data["page"]=page
	this.Data["pages"]=pages
	this.Data["goods"]=goods
	this.Data["active2"] = "active"
	this.Data["title"] = "全部订单"
	this.Layout = "layout.html"
	this.TplName = "user_center_order.html"
}

//用户名展示
func GetUserName(this *beego.Controller) interface{} {
	userName := this.GetSession("userName")

	if userName != nil {
		this.Data["sblock"] = "block"
		this.Data["name"] = userName.(string)
		this.Data["none"] = "none"
		this.Data["html"] = "/logout"
		this.Data["zuce"] = "退出"
	} else {
		this.Data["html"] = "/register"
		this.Data["zuce"] = "注册"
	}
	if userName != nil {
		return userName
	} else {
		return ""
	}

}

//查询用户
func GetUser(username string) (models.User, error) {
	o := orm.NewOrm()
	var user models.User
	user.Name = username
	err := o.Read(&user, "Name")
	return user, err
}

//地址展示
func (this *GoodsController) ShowSite() {
	username := GetUserName(&this.Controller)
	o := orm.NewOrm()
	user, _ := GetUser(username.(string))
	var address models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Id", user.Id).Filter("Isdefault", true).One(&address)
	this.Data["addr"] = address.Addr
	beego.Info(address)
	this.Data["username"] = address.Receiver
	this.Data["phone"] = address.Phone
	this.Data["active3"] = "active"
	this.Data["title"] = "全部地址"
	this.Layout = "layout.html"
	this.TplName = "user_center_site.html"
}

//输入地址
func (this *GoodsController) HandleSite() {
	username := GetUserName(&this.Controller)
	receiver := this.GetString("receiver")
	address1 := this.GetString("addr")
	zipcode := this.GetString("zipcode")
	phone := this.GetString("phone")
	o := orm.NewOrm()
	user, err := GetUser(username.(string))
	if err != nil {
		beego.Info(err)
		return
	}
	var addr models.Address
	addr.Isdefault = true
	o.QueryTable("Address").RelatedSel("User").Filter("User__Id",user.Id).Filter("Isdefault",true).Update(orm.Params{"Isdefault":false})
	/*err = o.Read(&addr, "Isdefault")
	if err == nil {
		addr.Isdefault = false
		o.Update(&addr)
	}*/
	var address models.Address
	address.User = &user
	address.Addr = address1
	address.Phone = phone
	address.Receiver = receiver
	address.Zipcode = zipcode
	address.Isdefault = true
	o.Insert(&address)
	this.Redirect("/goods/user_center_site", 302)

}

// 购物车
func (this *GoodsController) ShowCart() {
	username:=GetUserName(&this.Controller)
	user,_:=GetUser(username.(string))
	conn,err:=redis.Dial("tcp","192.168.110.111:6379")
	if err != nil {
		beego.Info("redis链接错误")
		return
	}
	defer conn.Close()
	goodsMap,_:=redis.IntMap(conn.Do("hgetall","cart_"+strconv.Itoa(user.Id)))
	 goods:=make([]map[string]interface{},len(goodsMap))
	 i:=0
	totalprice:=0
	totalcount:=0
	for k,v:=range goodsMap {
		skuid,_:=strconv.Atoi(k)
		o:=orm.NewOrm()
		var goodsSku models.GoodsSKU
		goodsSku.Id=skuid
		o.Read(&goodsSku)
		temp:=make(map[string]interface{})
		temp["goodsSku"]=goodsSku
		temp["count"]=v
		totalprice+=goodsSku.Price*v
		totalcount+=v
		temp["price"]=goodsSku.Price*v
		goods[i]=temp
		i+=1
	}
	this.Data["totalprice"]=totalprice
	this.Data["totalcount"]=totalcount
	this.Data["goods"]=goods
	this.Data["title"] = "购物车"
	this.Layout = "layout.html"
	this.TplName = "cart.html"
}

//我的订单
func (this *GoodsController) ShowplaceOder() {
	userName:=GetUserName(&this.Controller)
	user,_:=GetUser(userName.(string))
	skuids:=this.GetStrings("skuid")
	goods:=make([]map[string]interface{},len(skuids))
	o:=orm.NewOrm()
	conn,_:=redis.Dial("tcp","192.168.110.111:6379")
	totalPrice:=0
	totalCount:=0
	for i,skuid:=range skuids  {
		tmep:=make(map[string]interface{})
		id,_:=strconv.Atoi(skuid)
		var goodsSku models.GoodsSKU
		goodsSku.Id=id
		o.Read(&goodsSku)
		tmep["goods"]=goodsSku
		count,_:=redis.Int(conn.Do("hget","cart_"+strconv.Itoa(user.Id),id))
		tmep["count"]=count
		amount:=goodsSku.Price*count
		tmep["amount"]=amount
		goods[i]=tmep
		totalCount+=count
		totalPrice+=amount
	}
	this.Data["goods"]=goods
	var addrs []models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Id",user.Id).All(&addrs)
	this.Data["addrs"]=addrs
	this.Data["totalCount"]=totalCount
	this.Data["totalPrice"]=totalPrice
	youfei:=10
	this.Data["youfei"]=youfei
	this.Data["shifukuan"]=totalPrice+youfei
	this.Data["title"] = "我的订单"
	this.Data["skuids"]=skuids
	this.Layout = "layout.html"

	this.TplName = "place_order.html"
}

//查看商品详情
func (this *GoodsController) Showcentent() {
	userName := GetUserName(&this.Controller)
	user, _ := GetUser(userName.(string))
	id, err := this.GetInt("id")
	if err != nil {
		beego.Info("链接失败")
		return
	}
	o := orm.NewOrm()
	var goods models.GoodsSKU
	goods.Id = id
	o.QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").Filter("Id", id).One(&goods)
	var goodsType []models.GoodsType
	o.QueryTable("GoodsType").All(&goodsType)

	var goodsnew []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType", goods.GoodsType).OrderBy("Time").Limit(2, 0).All(&goodsnew)
	this.Data["goodsnew"] = goodsnew
	this.Data["goodsType"] = goodsType
	this.Data["goods"] = goods
	//判断用户是否存在
	if userName != nil {
		red, err := redis.Dial("tcp", "192.168.110.111:6379")
		if err != nil {
			beego.Error(err)
		}
		//把相同的记录删除
		red.Do("lrem", "history"+strconv.Itoa(user.Id), 0, id)
		//把最新的记录存起来
		red.Do("lpush", "history"+strconv.Itoa(user.Id), id)

	}
	this.Layout = "layout.html"
	this.TplName = "detail.html"
}
func showpage(page,pagecount int)[]int  {
	var pages []int
	if pagecount<=5 {
		pages=make([]int,pagecount)
		for i,_:=range pages{
			pages[i]=i+1
		}
	}else if page<=3{
		pages=make([]int,5)
		pages=[]int{1,2,3,4,5}
	} else if page>pagecount-3 {
		pages=[]int{pagecount-4,pagecount-3,pagecount-2,pagecount-1,pagecount}
	}else {
		pages=[]int{page-2,page-1,page,page+1,page+2}
	}
	return pages
}
//查看分类商品
func (this *GoodsController) Showlisst() {
	id, _ := this.GetInt("id")
	oder, _ := this.GetInt("oder")
	page, err := this.GetInt("page")
	if err != nil {
		beego.Info(err)
		page = 1
	}
	this.Data["oder"] = oder
	this.Data["id"] = id
	this.Data["page"] = page
	o:=orm.NewOrm()
	count,_:=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).Count()
	pageSize:=2
	start := (page-1) * pageSize
	pagecount:=math.Ceil(float64(count)/float64(pageSize))
	pages:=showpage(page,int(pagecount))
	this.Data["pages"]=pages
	var goods []models.GoodsSKU
	if oder == 1 {
		GetUserName(&this.Controller)
		o := orm.NewOrm()
		Showlist(&this.Controller)
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).Limit(2, start).All(&goods)
		this.Data["active"]="active"
	}else if oder == 2 {
		GetUserName(&this.Controller)
		o := orm.NewOrm()
		Showlist(&this.Controller)
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).OrderBy("Price").Limit(2, start).All(&goods)
		this.Data["active"]="active"

	}else {
		GetUserName(&this.Controller)
		o := orm.NewOrm()
		Showlist(&this.Controller)
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).OrderBy("Sales").Limit(2, start).All(&goods)
		this.Data["active"]="active"
	}

beego.Info(goods)
	this.Data["pagecount"] = int(pagecount)
	this.Data["goods"] = goods
	this.Layout = "layout.html"
	this.TplName = "list.html"
}
func Showlist(this *beego.Controller) {
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	var goodsType []models.GoodsType
	o.QueryTable("GoodsType").All(&goodsType)
	var goodstype models.GoodsType
	goodstype.Id = id
	o.Read(&goodstype)
	var goodsnew []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodstype.Id).OrderBy("Time").Limit(2, 0).All(&goodsnew)
	this.Data["goodsnew"] = goodsnew
	this.Data["goodsType"] = goodsType
}

//搜索
func (this *GoodsController) HandleSearch() {
	search := this.GetString("goodsName")
	o:=orm.NewOrm()
	var goods []models.GoodsSKU
	if search=="" {
		o.QueryTable("GoodsSKU").All(&goods)

	}else {
		o.QueryTable("GoodsSKU").Filter("Name__icontains",search).All(&goods)
	}
	this.Data["goods"]=goods
	this.Layout="layout.html"
	this.TplName="search.html"
}
//去支付
func (this *GoodsController)HandlePay(){
	appId:="2016092000557265"
	aliPublicKey:="MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAywC7gD+vrKqvgGemDzo3bltinzA4qa2sDJMxThEMa3tFikP1Gr9ZnskAT8mc7eb2x5XlbWVyiXeUfOJAplUzgL/a9HemRTRyLbtrPl/TLj5322y1wFqajU4Rt6MWLvu8wcCdaIZuLwYfbzShfHA5S2DfxgJF0XkOMsYitrosvtKzR78CE2eLpAlfitAKTYg8zart5sUHNXjV+tfuNKTyWiWjLCOjFHVUgwkAJjajU+iaKraeQCYaZuHaCRAIiwq3kI6q6XL8UG8VCZaf8j0P/cdl3M4f/gR5E3S+LkG70E6uCwrVhKJSG8nf6NtAP+cIs8sUYv8wpr4VwKWTZhNVOwIDAQAB"
	privateKey:="MIIEpQIBAAKCAQEAuGDcHWdHNLS25wL015JizgKgXU5RHdSKYnRoQLexQCXb8dTg"+
		"2kF8EZ4PUhwJ0s9fKu7RmHDT1K7sVdQ+WoTGmT5vPjqrKTGiNi0yLavKdZziUEie"+
		"nuChHzCVCm7roC89hCB/6qHVL92NCaqSwOaxdDEmCy3xXMklGrX+X6Zj9IU6WqSj"+
		"etqGYLtJ0Q64czQA2yM7i92NzaDtwqd4Fl40Ahb9WD+hr0BCOCskdepqfs5P6k+q"+
		"cz3FPJonIrYk1G8gGGGL6tg9pcD6Lmvr0mWeecM1uEjVwN+gZeGPzJ9ccViqEP3V"+
		"Pd1Xe665IOwfRm67MaX16wfbkfpOZ7kOj0l+WQIDAQABAoIBAQCJ/TO/bcQE1hrs"+
		"2XGUxKHdvGl4a1yaDq9i7+v2Q4QMlkj9vGxr7AaGyNx+fy168GgxIXsLs6VVz3Rg"+
		"5++inyxjFC79S7s9oT/dfAXJ2IA1dayKmU7daRAs35crr8f4omJPuGMDnwqGQDGF"+
		"wnsCk6TLaN0oEMJKxt9WFk7CFy1Hmgc8pNd1vnMFpNkztHGPMwCZFJ24+SQRCqAI"+
		"iAvjbdGS1S7DhxiChcdImzP381Da+XVMnBJlQsaFC7NdiL13XMp4fPRQYyiRPm/8"+
		"mtERWfYJcvTO8+/4ShrYWYjz72HPWclIHBdv12uU/jhBKDm8H0bV9JpM2yhh1EXf"+
		"YduyqpmBAoGBAO/zgUaFzsuHXwz150FHUjpuxE04wqcuwN/PaO7z1qu3cUHOtX6J"+
		"s7tXbq7a8HSnhKWSosLaJ9IwwoqCJzHjGKVT4Naln50nzdgkt3K5L5XYE2heZ0Jn"+
		"u25F6Uao49aSz0a7LQE3RR2oT2EBPirjgkXX2B75z51hx4Sh57TBrDe1AoGBAMS1"+
		"0za85W/AjvngzaotUc4LDe62h5rIODGXrNRRST/ijwJu/pzapj0QPYrCOtfpWwi4"+
		"CfCm4TFt3xV1skIRZNVuku4SMwbaMRdKfn721LS6W9eHg/dptRHAfrrkuPIliw4P"+
		"w7gyRr51VKtYQUsBYPCAM2LwcIsIgb+bcBZYYAqVAoGBAO+sRUNg4jcPh1SVxqDA"+
		"kZTGEROlD2EYZRTowkJzkshAWkNGKqky+DC3W1oSXD3ZGbicaDDC4SWlCJx69pVw"+
		"5aw1xQ4BrxW1rXko64gPC0Xb5z7HlNKSdHfoIuMuTS2FxL48te5R+5ptBKS7LhJ+"+
		"3x/OQhRmqAbmpPiJE7zL+q5FAoGAFqS3g32LC6omyyzNf+FnoUg0el4YjgCuN0c2"+
		"ZdpVjD0QKT+Nn5Crwiu0adyh2WjLSd2lh0Yudfony9iYhHJsIQVxdGYz6X4EWKIC"+
		"narcIVGycMTws/I/HaQC8pCRmY4oy52U8gcXjaUD8hVerruh5Q1c3O7AhcCc7ul9"+
		"pZTWuWECgYEA1l2UfH37W6ZtV7jcCiG7Y1uX6fRPcMjKngyOHFE56pagjgFZx0AU"+
		"7R6J7VmJdsUZXzGhvJKwBYSeLRGcIh12GAJ7eQLBpYrmS67d5azuMUh0YyKlJ7dK"+
		"LagMHqHf6W7x/AKXS6w2p6SuCX7rSgKlT+TZPUfSeQ4LD3FanED0phs="
	var client = alipay.New(appId, aliPublicKey, privateKey, false)
	orderId := this.GetString("orderId")
	totalPrice := this.GetString("totalPrice")
	//alipay.trade.page.pay
	var p = alipay.AliPayTradePagePay{}
	p.NotifyURL = "http://192.168.110.81:8080/user/payOk"
	p.ReturnURL = "http://192.168.110.111:8080/goods/payok"
	p.Subject = "天天生鲜"
	p.OutTradeNo = orderId
	p.TotalAmount = totalPrice
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err = client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	var payURL = url.String()

	this.Redirect(payURL,302)

}
//支付返回
func (this *GoodsController)ShowPay() {
	orderid:=this.GetString("out_trade_no")
	if orderid=="" {
		beego.Info("返回数据错误")
		return
	}
	o:=orm.NewOrm()
	count,_:=o.QueryTable("OrderInfo").Filter("OrderId",orderid).Update(orm.Params{"Orderstatus":2})
	if count==0 {
		beego.Info("更新数据失败")
		return
	}
	this.Redirect("/goods/user_center_order",302)
}