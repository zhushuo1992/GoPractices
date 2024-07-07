package main

import (
	"fmt"
	"gee"
	"log"
	"net/http"
	"text/template"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func v2middleware() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.SetFailInfo(500, "v2 is not useable")
		log.Printf("%d %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}

}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.Deafult()
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	r.GET("/", func(c *gee.Context) {
		c.SetStringRsp(http.StatusOK, "Hello zaizai\n")

	})

	r.GET("/panic", func(c *gee.Context) {
		names := []string{"ppp"}
		c.SetStringRsp(http.StatusOK, names[100])
	})

	v2 := r.Group("/v2")
	v2.Use(v2middleware())
	{
		v2.GET("/hello/:name", func(ctx *gee.Context) {
			ctx.SetStringRsp(http.StatusOK, "hello %s, you are at %s\n", ctx.GetUrlParam("name"), ctx.Path)
		})
	}

	studentGroup := r.Group("/student")
	{
		stu1 := &student{
			Name: "zaizai",
			Age:  1,
		}

		stu2 := &student{
			Name: "janson",
			Age:  31,
		}
		studentGroup.GET("/css", func(ctx *gee.Context) {
			ctx.SetHTMLRsp(http.StatusOK, "css.tmpl", nil)
		})

		studentGroup.GET("/all", func(ctx *gee.Context) {
			ctx.SetHTMLRsp(http.StatusOK, "arr.tmpl", gee.H{
				"title":  "gee",
				"stuArr": [2]*student{stu1, stu2},
			})
		})

		studentGroup.GET("/date", func(ctx *gee.Context) {
			ctx.SetHTMLRsp(http.StatusOK, "custom_func.tmpl", gee.H{
				"title": "gee",
				"now":   time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC),
			})
		})
	}

	r.Run(":9988")
}
