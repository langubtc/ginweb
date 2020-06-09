package main
import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)


func webRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	// 添加两个多模板继承, 初始模板必须写在前面。
	r.AddFromFiles("server", "templates/base.html", "templates/server.html")
	r.AddFromFiles("login", "templates/base.html", "templates/login.html")
	return r
}

func main() {
	r := gin.Default()
	//自定义模版变量必须放在解析模版之前
	r.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML{
			return template.HTML(str)
		},
	})

	// 加载静态文件,以static开头指向statics
	r.Static("/static","./statics")
	r.StaticFile("/favicon.ico","./statics/favicon.ico")


	//r.LoadHTMLFiles("templates")  //模版渲染 如果只有一层的情况下
	//r.LoadHTMLGlob("templates/**/*")  //**表示文件夹/*表示文件

	//使用模版
	r.HTMLRender=webRender()


	r.GET("/server",func(c *gin.Context){

		c.HTML(http.StatusOK,"server",gin.H{
			"name":"roddy",
		})
	})

	r.GET("/login",func(c *gin.Context){

		c.HTML(http.StatusOK,"login",gin.H{
			"name":"roddy",
		})
	})


	// json数据
	r.GET("/json",func(c *gin.Context){
		data := map[string]interface{}{
			"name":"roddy",
			"message":"hello world!",
			"age":28,
		}

		c.JSON(http.StatusOK,data)
	})

	// 结构体方式返回json数据
	r.GET("/structJson",func(c *gin.Context){
		type Msg struct{
			Name string `json:"name"`
			Age int
			Message string

		}
		data := Msg{"roddy",23,"this is "}

		c.JSON(http.StatusOK,data)
	})

	r.Run("127.0.0.1:8001")
}
