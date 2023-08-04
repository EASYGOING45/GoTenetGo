package main

import(
	"fmt"
	"html/template"
	"net/http"
	"time"
	"going"
)

type student struct {
	Name string
	Age int
}

func FormatAsDate(t time.Time) string{
	year,month,day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d",year,month,day)
}

func main(){
	r = going.New()
	r.Use(going.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate":FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets","./static")

	stu1 := &student{Name:"tenet",Age:21}
	stu2 := &student{Name:"k",Age:100}
	r.GET("/",func(c *going.Context){
		c.HTML(http.StatusOK,"css.tmpl",nil)
	})

	r.GET("/students",func(c *going.Context){
		c.HTML(http.StatusOK,"arr.tmpl",going.H{
			"title":"going",
			"stuArr":[2]*student{stu1,stu2},
		})
	})

	r.GET("/date",func(c *going.Context){
		c.HTML(http.StatusOK,"custon_func.tmpl",going.H{
			"title":"going",
			"now":time.Date(2019,8,17,0,0,0,0,time.UTC),
		})
	})

	fmt.Println("Server is running..."")
	r.Run(":9999")
}