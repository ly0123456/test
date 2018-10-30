package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ttsx/models"
	"math"
	"path"
	"time"

	"github.com/weilaihui/fdfs_client"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowIndex() {
  userName:=this.GetSession("useName")
	o := orm.NewOrm()
	qs := o.QueryTable("GoodsSKU")
	var goods []models.GoodsSKU
	//qs.All(&artcles)
	page, err := this.GetInt("id")
	if err != nil {
		page = 1
	}
	var count int64
	pagesize := 2

	start := (page - 1) * pagesize
	//qs.Limit(pagesize, start).All(&artcles)
	var types []models.GoodsType
	o.QueryTable("GoodsType").All(&types)
	this.Data["types"] = types
	typename:=this.GetString("select")
	if typename!="" {
		qs.Limit(pagesize,start).RelatedSel("GoodsType").Filter("GoodsType__Name",typename).All(&goods)
 		count,_=qs.Limit(pagesize,start).RelatedSel("GoodsType").Filter("GoodsType__Name",typename).Count()

	}else {
		qs.Limit(pagesize,start).All(&goods)
		count, _ = qs.Count()
	}
	pagecount := math.Ceil(float64(count) / float64(pagesize))
	this.Data["userName"]=userName
	this.Data["typename"]=typename
	this.Data["count"] = count
	this.Data["pagecount"] = int(pagecount)
	this.Data["page"] = page
	this.Data["goods"] = goods
	this.Layout="layout1.html"
	this.TplName = "index1.html"
}//主页
func Uploadname(this *ArticleController, filepath string, tplname string) string {
	file, head, err := this.GetFile(filepath)
	defer file.Close()
	if err != nil {
		beego.Info("上传失败，请重新上传")
		this.Data["errmsg"] = "上传失败，请重新上传"
		this.TplName = tplname
		return ""
	}
	if head.Size > 500000 {
		beego.Info("上传文件大了，请重新上传")
		this.Data["errmsg"] = "上传文件大了，请重新上传"
		this.TplName = tplname
		return ""
	}
	filename := head.Filename
	ext := path.Ext(filename)
	if ext != ".jpg" && ext != ".png" && ext != "" {
		beego.Info("上传文件格式错误，请重新上传")
		this.Data["errmsg"] = "上传文件格式错误，请重新上传"
		this.TplName = tplname
		return ""
	}
	clint,err:=fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
	if err!=nil {
		beego.Info(err)
	}
	fileBuffer:=make([]byte,head.Size)
	file.Read(fileBuffer)
	beego.Info(fileBuffer)
	res,err:=clint.UploadByBuffer(fileBuffer,ext[1:])
	beego.Info(res)
	if err!=nil {
		beego.Info(err)
	}

	return res.RemoteFileId
}
func (this *ArticleController) ShowAddGoods() {
	o := orm.NewOrm()
	var types []models.GoodsType
	o.QueryTable("GoodsType").All(&types)
	var goods []models.Goods
	o.QueryTable("Goods").All(&goods)
	this.Data["goods"]=goods
	this.Data["types"] = types
	this.Layout="layout1.html"
	this.TplName = "add1.html"
}//查看添加文章页
func Uploadnames(this *ArticleController ,filepath string,tplname string)  []string {
	var files []string
	imgsfilehander,err:=this.GetFiles(filepath)
	if err!=nil {
		this.Data["errmsg"]="文件问题"
		this.TplName=tplname
		return	[]string{}
	}
	for _,imgname:=range imgsfilehander{
		if imgname.Size > 500000 {
			beego.Info("上传文件大了，请重新上传")
			this.Data["errmsg"] = "上传文件大了，请重新上传"
			this.TplName = tplname
			return	[]string{}

		}

		ext := path.Ext(imgname.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".jpng" {
			beego.Info("上传文件格式错误，请重新上传")
			this.Data["errmsg"] = "上传文件格式错误，请重新上传"
			this.TplName = tplname
			return []string{}
		}
		filename:=time.Now().Format("2006-01-02-15-04-05")
		this.SaveToFile(imgname.Filename,"./static/img/"+filename+ext)
		file:="/static/img/"+filename+ext
		files=append(files,file)
	}
		return files
	}

func (this *ArticleController) HendleAddGoods() {
	TpyeName := this.GetString("selectType")
	goodsName := this.GetString("goodsName")
	goodsname:=this.GetString("selectGoodsSPU")
	desc:=this.GetString("desc")
	price,_:=this.GetInt("goodsPrice")
	//unite:=this.GetString("unite")
	stock,_:=this.GetInt("goodsStock")
	//Status,_:=this.GetInt("Status")
	//sales,_:=this.GetInt("sales")

	if goodsName == "" || goodsname == ""||desc=="" {
		beego.Info("数据不完整，请重新输入")
		this.Data["errmsg"] = "数据不完整，请重新输入"
		this.TplName = "add.html"
		return
	}
	filename:=Uploadname(this,"uploadname","/goods/addgoods")
	if filename==""{
		this.Data["errmsg"] = "图片有问题，请重新输入"
		this.TplName="add.html"
		return
	}
	o := orm.NewOrm()
var goods1 models.Goods
goods1.Name=goodsname
o.Read(&goods1)
	/*filepath:=Uploadnames(this,"uploadname","/goods/addgoods")
	if filepath==nil{
		this.Data["errmsg"] = "图片有问题，请重新输入"
		this.TplName="add.html"
		return
	}

*/
	var goods models.GoodsSKU
	goods.Stock=stock
	goods.Name=goodsName
	goods.Desc=desc
	goods.Price=price
	goods.Image=filename
	var Typename models.GoodsType
	Typename.Name = TpyeName
	o.Read(&Typename, "Name")
	goods.GoodsType = &Typename
	goods.Goods=&goods1
	o.Insert(&goods)
	/*var goodsImge models.GoodsImage
	for _,file:=range filepath {
		goodsImge.Image=file
	}
	goodsImge.GoodsSKU=&goods
	o.Insert(&goodsImge)*/

	this.Redirect("/goods/comeindex", 302)

}//处理添加文章

func (this *ArticleController) ShowDateil() {
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	var goods models.GoodsSKU
	goods.Id = id
o.QueryTable("GoodsSKU").RelatedSel("GoodsType"). Filter("Id",id).One(&goods)
	/*article.Acount += 1
	o.Update(&article)
	userName:=this.GetSession("useName")
	if userName==nil{
		this.Redirect("/login",302)
		return
	}

	m2m:=o.QueryM2M(&article,"Users")
	var usr models.User
	usr.Name=userName.(string)
	o.Read(&usr,"Name")
	m2m.Add(usr)
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)*/
	this.Data["goods"] = goods
	//this.Data["users"]=users
	this.Layout="layout1.html"
	this.TplName = "content.html"
}//

func (this *ArticleController) ShowUpdate() {
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	var good models.GoodsSKU
	good.Id = id
	o.Read(&good)
	this.Data["goods"] = good
	this.Layout="layout1.html"
	this.TplName = "update.html"
}//

func (this *ArticleController) HandleUpdate() {
	id, _ := this.GetInt("id")
	goodsName := this.GetString("goodsName")
	desc:=this.GetString("desc")
	price,_:=this.GetInt("goodsPrice")
	//unite:=this.GetString("unite")
	stock,_:=this.GetInt("goodsStock")

	filepath := Uploadname(this ,"uploadname", "update.html")

	if goodsName == "" || desc == "" {
		beego.Info("数据不完整，请重新输入")
		this.Data["errmsg"] = "数据不完整，请重新输入"
		this.TplName = "update.html"
		return
	}

	o := orm.NewOrm()
	var goods models.GoodsSKU
	goods.Id = id
	err := o.Read(&goods)
	if err != nil {
		beego.Info("链接错误")
		this.Data["errmsg"] = "链接错误"
		this.TplName = "update.html"
		return
	}
	goods.Name= goodsName
	goods.Desc= desc
	goods.Price = price
	goods.Stock=stock
	if filepath!= "" {
		goods.Image = filepath
	}
	o.Update(&goods)
	this.Redirect("/goods/comeindex", 302)
}

func (this *ArticleController) ShowDelete() {
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	var goods models.Goods
	goods.Id = id
	o.Delete(&goods)
	this.Redirect("/goods/comeindex", 302)
}
func (this *ArticleController) ShowAddType() {
	o := orm.NewOrm()
	qs := o.QueryTable("GoodsType")
	var GoodsTypes []models.GoodsType
	qs.All(&GoodsTypes)
	this.Data["GoodsTypes"] = GoodsTypes
	this.Layout="layout1.html"
	this.TplName = "addType1.html"
}
func (this *ArticleController) HandleAddtype() {
	typeName := this.GetString("typeName")
	logfilepath:=Uploadname(this,"log","addType.html")
	imgfilepath:=Uploadname(this,"img","addType.html")
	if typeName == "" {
		beego.Info("类型名不能为空,请重新输入")
		this.Data["errmsg"] = "类型名不能为空,请重新输入"
		this.TplName = "addType.html"
	}
	o := orm.NewOrm()
	var GoodsType models.GoodsType
	GoodsType.Name = typeName
	err := o.Read(&GoodsType, "Name")
	if err == nil {
		beego.Info("类型已存在,请重新输入")
		this.Data["errmsg"] = "类型已存在,请重新输入"
		this.TplName = "addType.html"
		return
	}
	GoodsType.Logo=logfilepath
	GoodsType.Image=imgfilepath
	o.Insert(&GoodsType)
	this.Redirect("/goods/addType", 302)
}
func (this *ArticleController) ShowDeleteType() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Info("请求路径错误")
		return
	}
	o := orm.NewOrm()
	var GoodsType models.GoodsType
	GoodsType.Id = id
	o.Delete(&GoodsType)
	this.Redirect("/goods/addType", 302)

}
func (this *ArticleController)ShowAddSPU()  {
	o:=orm.NewOrm()
	var goods []models.Goods
	o.QueryTable("Goods").All(&goods)
	this.Data["goods"]=goods
	this.Layout="layout1.html"
	this.TplName="addGoodsSPU.html"
}
func (this *ArticleController)HandleAddSPU()  {
	spuNamethis:=this.GetString("spuName")
	spuDetail:=this.GetString("spuDetail")
	o:=orm.NewOrm()
	var goods models.Goods
	goods.Name=spuNamethis
	goods.Detail=spuDetail
	o.Insert(&goods)
	this.Redirect("/goods/AddGoodsSPU",302)

}