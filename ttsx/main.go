package main

import (
	_ "ttsx/routers"
	"github.com/astaxie/beego"
	_"ttsx/models"
)

func main() {
	beego.AddFuncMap("perpage",Prepage)
	beego.AddFuncMap("nextpage",Nextpage)
	beego.AddFuncMap("index",Index)
	beego.Run()
}
func Prepage(page int)  int{
	if page==1 {
		return 1
	}
	return page-1
}
func Nextpage(page,pagecount int ) int {
	if page>=pagecount {
		return pagecount
	}
	return page+1
}
func Index(i int )int   {
	return i+1
}
