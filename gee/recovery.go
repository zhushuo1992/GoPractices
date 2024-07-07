package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 打印调用栈
func trace(message string) string {

	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) //跳过前3个caller

	var str strings.Builder

	str.WriteString(message + "\nTraceback:")

	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}

	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				//代表有异常被recover捕获到了
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.SetFailInfo(http.StatusInternalServerError, "internal Server Error")
			}
		}()

		c.NextMiddle() //如果有人调用了Recovery，那Recovery负责后续的调用
	}

}
