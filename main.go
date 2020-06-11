package main

import (
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
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
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})

	// 加载静态文件,以static开头指向statics
	r.Static("/static", "./statics")
	r.StaticFile("/favicon.ico", "./statics/favicon.ico")

	//r.LoadHTMLFiles("templates")  //模版渲染 如果只有一层的情况下
	//r.LoadHTMLGlob("templates/**/*")  //**表示文件夹/*表示文件

	//使用模版
	r.HTMLRender = webRender()

	r.GET("/server", func(c *gin.Context) {

		queryName := c.DefaultQuery("name", "default")
		queryAge := c.DefaultQuery("age", "18")
		query, ok := c.GetQuery("queryName")

		if ok == false {
			query = "123"
		}

		c.HTML(http.StatusOK, "server", gin.H{
			"name":  queryName,
			"age":   queryAge,
			"query": query,
		})
	})

	r.GET("/login", func(c *gin.Context) {

		c.HTML(http.StatusOK, "login", nil)
	})

	r.POST("/login", func(c *gin.Context) {

		username, ok := c.GetPostForm("username")
		if username == "" {
			username = "rrr"
		}

		password, ok := c.GetPostForm("password")

		fmt.Println(username, password, ok)
		if password == "" {
			password = "xxxxx"
		}

		fmt.Println(username, password)
		c.HTML(http.StatusOK, "server", gin.H{
			"username": username,
			"password": password,
		})
	})

	r.GET("/blog/:year/:month", func(c *gin.Context) {
		year := c.Param("year")
		month := c.Param("month")

		c.JSON(http.StatusOK, gin.H{
			"year":  year,
			"month": month,
		})
	})

	// json数据
	r.GET("/json", func(c *gin.Context) {
		data := map[string]interface{}{
			"name":    "roddy",
			"message": "hello world!",
			"age":     28,
		}

		c.JSON(http.StatusOK, data)
	})

	// 结构体方式返回json数据
	r.GET("/structJson", func(c *gin.Context) {
		type Msg struct {
			Name    string `json:"name"`
			Age     int
			Message string
		}
		data := Msg{"roddy", 23, "this is "}

		c.JSON(http.StatusOK, data)
	})

	type UserInfo struct {
		Username string `json:"username" `
		Password string `json:"password" `
	}

	// shouldbind数据绑定
	r.POST("/json2", func(c *gin.Context) {
		var u UserInfo

		err := c.ShouldBind(&u)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":   "ok",
				"username": u.Username,
				"password": u.Password,
			})

		}

	})

	// 上传文件
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		log.Print(file.Filename)
		dst := fmt.Sprintf("./upload/%s", file.Filename)
		c.SaveUploadedFile(file, dst)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})

	})

	// 上传多个文件
	r.POST("/multiUpload", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["file"]

		for index, file := range files {

			log.Print(file.Filename)
			dst := fmt.Sprintf("./upload/%s%d", file.Filename, index)
			c.SaveUploadedFile(file, dst)

		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%d' uploaded!", len(files)),
		})

	})

	// ShouldBindQuery 只返回参数，没有返回参数的值 any表示接收任何请求
	r.Any("/json3", func(c *gin.Context) {
		var u UserInfo

		err := c.ShouldBindQuery(&u)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":   "ok",
				"username": u.Username,
				"password": u.Password,
			})

		}

	})

	// 路由组
	NewRouter := r.Group("/router")

	NewRouter.GET("info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"path": "info",
		})
	})
	NewRouter.GET("news", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"path": "news",
		})
	})

	// 请求重定向到百度

	r.GET("/redirect/baidu", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})

	// 路由重定向
	r.GET("/luyou", func(c *gin.Context) {
		//跳转到/luyou2对应的路由处理函数
		c.Request.URL.Path = "/router/info" //把请求的URL修改
		r.HandleContext(c)                  //继续后续处理
	})

	// 优化404请求
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "404",
		})
	})

	r.Run("127.0.0.1:8001")
}
