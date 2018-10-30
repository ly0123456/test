package routers

import (
	"练习2/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/comeindex",&controllers.ArticleController{},"get:ShowIndex")
    beego.Router("/addarticle",&controllers.ArticleController{},"get:ShowAddArticle;post:Handleadd")
    beego.Router("/dateil",&controllers.ArticleController{},"get:Showdateil")
    beego.Router("/update",&controllers.ArticleController{},"get:ShowUpdate;post:Handleupdate")

    beego.Router("/delete",&controllers.ArticleController{},"get:ShowDelete")

    beego.Router("/addtype",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    beego.Router("/deleteType",&controllers.ArticleController{},"get:ShowDateletype")
}
