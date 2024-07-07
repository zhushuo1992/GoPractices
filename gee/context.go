package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter //接口类型，不需要指针
	Req        *http.Request       //struct，需要指针
	Path       string
	Method     string
	Params     map[string]string //url参数
	StatusCode int
	handlers   []HandlerFunc //中间件Middleware
	index      int           //中间件索引，默认值是-1
	engine     *Engine       //全局1个的引擎
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Path:   r.URL.Path,
		Method: r.Method,
		Req:    r,
		Writer: w,
		index:  -1,
	}
}

// 调用下面所有的中间件
func (c *Context) NextMiddle() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// 设置错误码、错误信息
func (c *Context) SetFailInfo(code int, err string) {
	c.index = len(c.handlers)
	c.SetJsonRsp(code, H{"message": err})
}

// 根据Key获取url参数(/a.html)的值
func (c *Context) GetUrlParam(key string) string {
	value := c.Params[key]
	return value
}

// 根据key获取Form中的值
func (c *Context) GetPostForm(key string) string {
	return c.Req.FormValue(key)
}

// 根据Key获取请求参数(?key=val)的值
func (c *Context) GetQueryVal(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置返回码
func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置消息头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 设置text格式返回  values:接受零个或多个参数切片
func (c *Context) SetStringRsp(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 设置json格式返回
func (c *Context) SetJsonRsp(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500) //失败的话，直接把encode的错误信息返回给客户端
	}
}

// 返回二进制数据？
func (c *Context) SetDataRsp(code int, data []byte) {
	c.SetStatus(code)
	c.Writer.Write(data)
}

// 返回html数据
func (c *Context) SetHTMLRsp(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.SetFailInfo(500, err.Error())
	}

}
