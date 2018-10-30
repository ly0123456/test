package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"练习2/models"
	"math"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowIndex() {

	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var artcles []models.Article
	//qs.All(&artcles)
	page, err := this.GetInt("id")
	if err != nil {
		page = 1
	}
	count, _ := qs.Count()
	pagesize := 2
	pagecount := math.Ceil(float64(count) / float64(pagesize))
	start := (page - 1) * pagesize
	//qs.Limit(pagesize, start).All(&artcles)
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types
	typename:=this.GetString("select")
  //var  typeName models.ArticleType
  //typeName.Tname=typename
  //o.Read(&typeName)
  qs.Limit(pagesize,start).RelatedSel("ArticleType").Filter("ArticleType__Tname",typename).All(&artcles)
	this.Data["count"] = count
	this.Data["pagecount"] = int(pagecount)
	this.Data["page"] = page
	this.Data["articles"] = artcles
	this.TplName = "index.html"
}
func (this *ArticleController) ShowAddArticle() {
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types
	this.TplName = "add.html"
}
func (this *ArticleController) Handleadd() {
	TpyeName := this.GetString("select")
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	if articleName == "" || content == "" {
		beego.Info("数据不完整，请重新输入")
		this.Data["errmsg"] = "数据不完整，请重新输入"
		this.TplName = "add.html"
		return
	}

	file, head, err := this.GetFile("uploadname")
	defer file.Close()
	if err != nil {
		beego.Info("上传失败，请重新上传")
		this.Data["errmsg"] = "上传失败，请重新上传"
		this.TplName = "add.html"
		return
	}
	if head.Size > 500000 {
		beego.Info("上传文件大了，请重新上传")
		this.Data["errmsg"] = "上传文件大了，请重新上传"
		this.TplName = "add.html"
		return
	}
	filename := head.Filename
	ext := path.Ext(filename)
	if ext != ".jpg" && ext != ".png" && ext != "" {
		beego.Info("上传文件格式错误，请重新上传")
		this.Data["errmsg"] = "上传文件格式错误，请重新上传"
		this.TplName = "add.html"
		return
	}
	filename = time.Now().Format("2006-01-02-15-04-05")
	this.SaveToFile("uploadname", "./static/img/"+filename+ext)
	o := orm.NewOrm()
	var article models.Article
	article.ArtiName = articleName
	article.Acontent = content
	article.Aimg = "/static/img/" + filename + ext
	var Typename models.ArticleType
	Typename.Tname = TpyeName
	o.Read(&Typename, "Tname")
	article.ArticleType = &Typename
	o.Insert(&article)
	this.Redirect("/comeindex", 302)

}
func (this *ArticleController) Showdateil() {
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)
	article.Acount += 1
	o.Update(&article)
	this.Data["article"] = article
	this.TplName = "content.html"
}
func (this *ArticleController) ShowUpdate() {
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)
	this.Data["article"] = article
	this.TplName = "update.html"
}
func Uploadname(this *beego.Controller, filepath string, tplname string) (string, string) {
	file, head, err := this.GetFile(filepath)
	defer file.Close()
	if err != nil {
		beego.Info("上传失败，请重新上传")
		this.Data["errmsg"] = "上传失败，请重新上传"
		this.TplName = tplname
		return "", ""
	}
	if head.Size > 500000 {
		beego.Info("上传文件大了，请重新上传")
		this.Data["errmsg"] = "上传文件大了，请重新上传"
		this.TplName = tplname
		return "", ""
	}
	filename := head.Filename
	ext := path.Ext(filename)
	if ext != ".jpg" && ext != ".png" && ext != "" {
		beego.Info("上传文件格式错误，请重新上传")
		this.Data["errmsg"] = "上传文件格式错误，请重新上传"
		this.TplName = tplname
		return "", ""
	}
	filename = time.Now().Format("2006-01-02-15-04-05")
	this.SaveToFile(filepath, "./static/img/"+filename+ext)
	return "/static/img/" + filename + ext, ext
}
func (this *ArticleController) Handleupdate() {
	id, _ := this.GetInt("id")
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	filepath, ext := Uploadname(&this.Controller, "uploadname", "update.html")

	if articleName == "" || content == "" {
		beego.Info("数据不完整，请重新输入")
		this.Data["errmsg"] = "数据不完整，请重新输入"
		this.TplName = "update.html"
		return
	}
	/*	file,head,err:=this.GetFile("uploadname")
		defer file.Close()
		if err!=nil {
			beego.Info("上传失败，请重新上传")
			this.Data["errmsg"]="上传失败，请重新上传"
			this.TplName="update.html"
			return
		}
		if head.Size>500000 {
			beego.Info("上传文件大了，请重新上传")
			this.Data["errmsg"]="上传文件大了，请重新上传"
			this.TplName="update.html"
			return
		}
		filename:=head.Filename
		ext:=path.Ext(filename)
		if ext!=".jpg"&&ext!=".png"&&ext!="" {
			beego.Info("上传文件格式错误，请重新上传")
			this.Data["errmsg"]="上传文件格式错误，请重新上传"
			this.TplName="update.html"
			return
		}
		filename=time.Now().Format("2006-01-02-15-04-05")
		this.SaveToFile("uploadname","./static/img/"+filename+ext)*/
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err := o.Read(&article)
	if err != nil {
		beego.Info("链接错误")
		this.Data["errmsg"] = "链接错误"
		this.TplName = "update.html"
		return
	}
	article.ArtiName = articleName
	article.Acontent = content
	article.Atime = time.Now()
	if ext != "" {
		article.Aimg = filepath
	}
	o.Update(&article)
	this.Redirect("/comeindex", 302)
}
func (this *ArticleController) ShowDelete() {
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Delete(&article)
	this.Redirect("/comeindex", 302)
}
func (this *ArticleController) ShowAddType() {
	o := orm.NewOrm()
	qs := o.QueryTable("ArticleType")
	var ArticleTypes []models.ArticleType
	qs.All(&ArticleTypes)
	this.Data["ArticleTypes"] = ArticleTypes
	this.TplName = "addType.html"
}
func (this *ArticleController) HandleAddType() {
	typeName := this.GetString("typeName")
	if typeName == "" {
		beego.Info("类型名不能为空,请重新输入")
		this.Data["errmsg"] = "类型名不能为空,请重新输入"
		this.TplName = "addType.html"
	}
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Tname = typeName
	err := o.Read(&articleType, "Tname")
	if err == nil {
		beego.Info("类型已存在,请重新输入")
		this.Data["errmsg"] = "类型已存在,请重新输入"
		this.TplName = "addType.html"
		return
	}
	o.Insert(&articleType)
	this.Redirect("/addtype", 302)
}
func (this *ArticleController) ShowDateletype() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Info("请求路径错误")
		return
	}
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	o.Delete(&articleType)
	this.Redirect("/addtype", 302)

}
