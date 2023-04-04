package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"study/controller"
	"study/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skratchdot/open-golang/open"
)

var router *gin.Engine

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("view/templates/*")

	// GET
	/*
		// HEALTHCHECK
		r.GET("/health", controller.HealthCheck)
		// Retrieving data
		r.GET("/cities", controller.GetCities)
		r.GET("/districts", controller.GetDistricts)
		r.GET("/dongs", controller.GetDongs)

	*/

	// web framework
	r.GET("/clientLogin", controller.GetLoginForm)
	r.GET("/users", controller.GetUsers)
	r.GET("/products", controller.GetProducts)
	r.GET("/register", controller.GetRegister)
	r.GET("/login", controller.GetLoginForm)
	r.GET("/unregistered-user", controller.GetUnregisteredForm)
	r.GET("/chrooms", controller.GetChRooms)
	r.GET("/mainpage", func(c *gin.Context) {
		c.HTML(http.StatusOK, "mainpage.html", nil)
		c.String(http.StatusOK, controller.Loguser.Name+"님 환영합니다")
	})

	r.POST("/chat", controller.AddChat)

	//버튼 추가 페이지 GET
	r.GET("/products/:id", controller.ShowPDetail)
	// 데이터 보여주는 페이지
	r.GET("/data", controller.ShowData)
	r.GET("/data/:id", controller.ShowData)
	// 게시물 리스트 보여주는 페이지
	r.GET("/board", func(c *gin.Context) {
		controller.ShowBoardList(c)
	})
	// 게시물
	r.GET("/postings", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post.html", nil)
	})

	// POST
	// 채팅방으로 이동시키는 POST 핸들러
	r.POST("/create-data", controller.CreateChatRoom)
	// ???
	r.POST("/board", func(c *gin.Context) {
		val := c.PostForm("action")
		if val == "post" {
			c.Redirect(http.StatusFound, "./postings")
		}
		if val == "search" {
			// 추가 예정
		}
	})

	// 게시물 올리기
	// FIXME : main.go 파일 말고 다른곳에서 함수 선언하고 여기서 호출만 하고 싶음...
	r.POST("/postings", func(c *gin.Context) {
		title := c.PostForm("title")
		content := c.PostForm("content")
		if controller.Loguser != nil {
			if err := controller.Posting(title, content); err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				c.Redirect(http.StatusFound, "/board")
			}
		} else {
			c.Redirect(http.StatusFound, "/")
		}
	})
	// ???
	r.POST("/postings/:id", func(c *gin.Context) {
		if controller.Loguser != nil {
			pid, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				fmt.Println(err)
			}
			content := c.PostForm("content")
			if err := controller.Comment(uint(pid), content); err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				c.Redirect(http.StatusFound, "postings/:id")
			}
		} else {
			c.Redirect(http.StatusFound, "/")
		}
	})
	r.GET("/posts/:id", controller.ShowPost)
	r.POST("/posts/:id", func(c *gin.Context) {
		action := c.PostForm("action")
		if controller.Loguser != nil && action == "comment" {
			pid, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				fmt.Println(err)
			}
			content := c.PostForm("content")
			if err := controller.Comment(uint(pid), content); err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				c.Redirect(http.StatusFound, "/posts/"+c.Param("id"))
			}
		} else if controller.Loguser == nil {
			c.Redirect(http.StatusFound, "/")
		}
		//else if controller.Loguser!=nil && action=="coc"{
		//	//cid,
		//}

	})

	// Putting data directly into database (using postman)

	r.POST("/cities", controller.PostCities)
	r.POST("/districts", controller.PostDistricts)
	r.POST("/dongs", controller.PostDongs)
	r.POST("/users", controller.PostUsers)
	r.POST("/productsStatus", controller.PostProductsStatus)
	r.POST("/productsCategory", controller.PostProductsCategory)
	r.POST("/products", controller.PostProducts)

	// Register , Login

	r.POST("/register", controller.PostRegister)
	r.POST("/login", func(c *gin.Context) {
		controller.LoginHandler(c)
	})

	go func() {
		if err := r.Run("127.0.0.1:1212"); err != nil {
			log.Fatalf("Error running Gin: %v", err)
		}
	}()

	// wait for the server to start
	time.Sleep(100 * time.Millisecond)

	url := "http://localhost:1212/login"
	err := open.Start(url)
	if err != nil {
		log.Fatalf("Error opening website: %v", err)
	}

	return r
}

func main() {
	// Set up the router
	router = setupRouter()

	// connect to the database
	model.ConnectDatabase()

	// Start the server
	err := router.Run()
	if err != nil {
		log.Fatal("Error starting Gin server: ", err)
	}
}
