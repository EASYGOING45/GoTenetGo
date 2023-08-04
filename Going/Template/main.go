package main

/*
(1) render array
$ curl http://localhost:9999/date
<html>
<body>
    <p>hello, going</p>
    <p>Date: 2019-08-17</p>
</body>
</html>
*/

/*
(2) custom render function
$ curl http://localhost:9999/students
<html>
<body>
    <p>hello, going</p>
    <p>0: goingktutu is 20 years old</p>
    <p>1: Jack is 22 years old</p>
</body>
</html>
*/

/*
(3) serve static files
$ curl http://localhost:9999/assets/css/goingktutu.css
p {
    color: orange;
    font-weight: 700;
    font-size: 20px;
}
*/

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"going"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := going.New()
	r.Use(going.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "goingktutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *going.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *going.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", going.H{
			"title":  "going",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *going.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", going.H{
			"title": "going",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})
	fmt.Println("Server running on port 9999.......")
	r.Run(":9999")
}
