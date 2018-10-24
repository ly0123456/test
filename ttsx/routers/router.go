package routers

import (
	"ttsx/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//路由1
	beego.InsertFilter("/goods/*",beego.BeforeRouter, func(context *context.Context) {
		userName:=context.Input.Session("userName")
		if userName==nil {
			context.Redirect(302,"/login")
		}
	})
	//路由2

    beego.Router("/", &controllers.MainController{})
    //注册
    beego.Router("/register",&controllers.UserControler{},"get:ShowReg;post:HandleReg")
    //登录
    beego.Router("/login",&controllers.UserControler{},"get:ShowLog;post:HandleLog")
    //处理激活
    beego.Router("/active",&controllers.UserControler{},"get:Handleactive")
	//退出
	beego.Router("/logout",&controllers.UserControler{},"get:ShowOut")
	//主页
	beego.Router("/showindex",&controllers.GoodsController{},"get:ShowIndex")
    //用户中心
    beego.Router("/goods/user_center_info",&controllers.GoodsController{},"get:ShowCenter")
	//全部订单
	beego.Router("/goods/user_center_order",&controllers.GoodsController{},"get:ShowOrder")
	//全部地址
	beego.Router("/goods/user_center_site",&controllers.GoodsController{},"get:ShowSite;post:HandleSite")
	//购物车
	beego.Router("/goods/cart",&controllers.GoodsController{} ,"get:ShowCart")
	//我的订单
	beego.Router("/goods/Order",&controllers.GoodsController{},"post:ShowplaceOder")
	//查看商品详情
	beego.Router("/showcentent",&controllers.GoodsController{},"get:Showcentent")
	//查看同一类型的商品
	beego.Router("/list",&controllers.GoodsController{},"get:Showlisst")
	//搜索
	beego.Router("/searchGoods",&controllers.GoodsController{},"post:HandleSearch")
	//添加购物车
	beego.Router("/goods/cart",&controllers.CartController{},"post:HandleCart")
	//改变商品数量
	beego.Router("/goods/updateCart",&controllers.CartController{},"post:HandleUpdateCart")
	//删除购物车
	beego.Router("/goods/deleteCart",&controllers.CartController{},"post:HandleDeleteCart")
	//提交订单
	beego.Router("/goods/addOrder",&controllers.CartController{},"post:HandleAddOrder")
	//付款
	beego.Router("/goods/pay",&controllers.GoodsController{},"get:HandlePay")
	//返回数据
	beego.Router("/goods/payok",&controllers.GoodsController{},"get:ShowPay")
	beego.InsertFilter("/goods/*",beego.BeforeExec, func(context *context.Context) {
		userName:=context.Input.Session("userName")
		if userName==nil {
			context.Redirect(302,"/login1")
		}
	})
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register1",&controllers.UserController{},"get:ShowRegister;post:Hendlerregister")
	beego.Router("/login1",&controllers.UserController{},"get:ShowLogin;post:HendleLogin")
	beego.Router("/goods/comeindex",&controllers.ArticleController{},"get:ShowIndex")
	beego.Router("/goods/addgoods",&controllers.ArticleController{},"get:ShowAddGoods;post:HendleAddGoods")
	beego.Router("/goods/AddGoodsSPU",&controllers.ArticleController{},"get:ShowAddSPU;post:HandleAddSPU")
	beego.Router("/goods/showDateil",&controllers.ArticleController{},"get:ShowDateil")
	beego.Router("/goods/updategoods",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
	//删除数据
	beego.Router("/goods/delete",&controllers.ArticleController{},"get:ShowDelete")
	//添加类型
	beego.Router("/goods/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddtype")
	beego.Router("/goods/deleteType" , &controllers.ArticleController{},"get:ShowDeleteType")
	beego.Router("/goods/logout",&controllers.UserController{},"get:Showlogout")
}

