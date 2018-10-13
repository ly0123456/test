package main

import (
	_ "练习2/routers"
	"github.com/astaxie/beego"
	_"练习2/models"
)

func main() {
	beego.AddFuncMap("perpage",Perpage)
	beego.AddFuncMap("nextpage",Nextpage)
	beego.AddFuncMap("index" ,Index)
	beego.Run()
}
func Perpage(page int) int {
	if page==1 {
		return 1
	}
	return page-1
}
func Nextpage(page ,pagecount int)int  {
	if page==pagecount {
		return page
	}
	return page+1
}
func Index(index int ) int {
	return index+1
}