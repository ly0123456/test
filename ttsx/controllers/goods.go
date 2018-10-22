package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ttsx/models"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"math"
)

type GoodsController struct {
	beego.Controller
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
	GetUserName(&this.Controller)
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
	var addr models.Address
	addr.Isdefault = true
	err := o.Read(&addr, "Isdefault")
	if err == nil {
		addr.Isdefault = false
		o.Update(&addr)
	}
	user, err := GetUser(username.(string))
	if err != nil {
		beego.Info(err)
		return
	}
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
	GetUserName(&this.Controller)
	this.Data["title"] = "我的订单"
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
