package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"study/model"
)

var Loguser *model.User

// 상품 상세 페이지 (대댓글용)
func ShowPDetail(c *gin.Context) {
	db, _ := model.ConnectDatabase()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	id := c.Param("id")

	var prod model.Product
	err1 := db.First(&prod, id).Error
	if err1 != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if c != nil {
		tmpl, err := template.ParseFiles("view/templates/detail.html")
		if err != nil {
			fmt.Println("error while parsing detail.html", err)
			log.Fatalf("%v", err)
		}
		err = tmpl.Execute(c.Writer, gin.H{
			"LID":         Loguser.ID,
			"NowUserName": Loguser.Name,
			"ID":          prod.ID,
			"Pname":       prod.Pname,
			"Price":       prod.Price,
			"PCID":        prod.PCID,
			"Ptext":       prod.Ptext,
			"UID":         prod.UID,
		})
		if err != nil {
			fmt.Println("error while executing detail.html", err)
			log.Fatalf("%v", err)
		}
	}
}

// CreateHandler - Chat 버튼 누르면 채팅룸 만들고 해당 페이지로 이동시키는 메소드
func CreateChatRoom(c *gin.Context) {
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	if err != nil {
		log.Fatalf("Error connecting database %v", err)
	}

	uidStr := c.PostForm("uid")
	uid, err1 := strconv.ParseUint(uidStr, 10, 32)
	pidStr := c.PostForm("pid")
	pid, err2 := strconv.ParseUint(pidStr, 10, 32)
	fmt.Println("pid", pid, "uid", uid)

	if err1 != nil || err2 != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Invalid ID"})
		fmt.Println("err1 and err2 is nil")
		return
	}

	// Verify that the product exists
	var product model.Product
	if err := db.Where("id = ?", pid).First(&product).Error; err != nil {
		fmt.Println("Product not found")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// chatting room duplication check
	var existingData model.ChRoom
	result := db.Where("s_id = ? AND b_id = ? AND p_id = ?", uint(uid), uint(Loguser.ID), uint(pid)).First(&existingData)
	if result.RowsAffected > 0 {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/data/%d", existingData.ID))
		log.Printf("Redirecting to /data/%d/n", existingData.ID)
		return
	}

	// if not duplicated, then create new data
	newChRoom := &model.ChRoom{SID: uint(uid), BID: Loguser.ID, PID: uint(pid)}

	// Set the p_id value to the valid product id
	newChRoom.PID = product.ID

	tx := db.Begin()
	if err := tx.Create(newChRoom).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	tx.Commit()

	// 데이터 생성 후 해당 데이터 출력 페이지로 redirect
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/data/%d", newChRoom.ID))
	log.Printf("Redirecting to /data/%d\n", newChRoom.ID)
}

// 채팅 페이지 GET
func ShowData(c *gin.Context) {
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	//채팅방 정보
	idStr := c.Param("id")

	fmt.Println("id = ", idStr)

	var chR model.ChRoom
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Invalid ID"})
		return
	}

	err = db.First(&chR, id).Error
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"message": "Data not found"})
		return
	}

	//채팅 메세지들
	var chats []model.Chat
	if err := db.Where("r_id = ?", id).Find(&chats).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if c != nil {
		tmpl, err := template.ParseFiles("view/templates/chroom.html")
		if err != nil {
			fmt.Println("error while parsing chroom.html", err)
			log.Fatalf("%v", err)
		}
		err = tmpl.Execute(c.Writer, gin.H{
			"ID":    chR.ID,
			"BID":   chR.BID,
			"SID":   chR.SID,
			"PID":   chR.PID,
			"chats": chats,
		})
		if err != nil {
			fmt.Println("error while executing chroom.html", err)
			log.Fatalf("%v", err)
		}
	}
}

func ShowPost(c *gin.Context) {
	db, _ := model.ConnectDatabase()

	id := c.Param("id")

	var post model.Post
	err1 := db.First(&post, id).Error
	if err1 != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var comments []model.Comment
	result := db.Where("post_id=?", post.ID).Find(&comments)
	if result != nil {
		fmt.Println("댓글없음")
	}
	cid := make([]uint, len(comments))
	for i := range comments {
		cid[i] = comments[i].ID
	}

	if c != nil {
		tmpl, err := template.ParseFiles("view/templates/postdetails.html")
		if err != nil {
			fmt.Println("error while parsing detail.html", err)
			log.Fatalf("%v", err)
		}
		var user model.User
		db.First(&user, post.UserID)
		err = tmpl.Execute(c.Writer, gin.H{
			"ID":        post.ID,
			"Ptitle":    post.Title,
			"Pcontent":  post.Content,
			"Pcomments": post.Comments,
			"UName":     user.Name,
			"comments":  comments,
		})
		if err != nil {
			fmt.Println("error while executing detail.html", err)
			log.Fatalf("%v", err)
		}
	}
}

func GetProducts(c *gin.Context) {
	// Connect to the database
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer sqlDB.Close()
	var prods []model.Product

	err2 := db.Find(&prods).Error
	if err2 != nil {
		log.Fatalf("%v", err)
	}

	if c != nil {
		tmpl, err := template.ParseFiles("view/templates/products.html")
		if err != nil {
			fmt.Println("error while parsing products.html", err)
			log.Fatalf("%v", err)
		}
		err = tmpl.Execute(c.Writer, gin.H{
			"prods": prods,
		})
		if err != nil {
			fmt.Println("error while executing products.html", err)
			log.Fatalf("%v", err)
		}
	}
}

// 메인페이지에서 채팅방 목록 불러오기
func GetChRooms(c *gin.Context) {
	// Connect to the database
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}

	var chrooms []model.ChRoom

	if err := db.Where("b_id = ?", Loguser.ID).Or("s_id = ?", Loguser.ID).Find(&chrooms).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if c != nil {
		tmpl, err := template.ParseFiles("view/templates/chrooms.html")
		if err != nil {
			fmt.Println("error while parsing chrooms.html", err)
			log.Fatalf("%v", err)
		}
		err = tmpl.Execute(c.Writer, gin.H{
			"chrooms": chrooms,
		})
		if err != nil {
			fmt.Println("error while executing chrooms.html", err)
			log.Fatalf("%v", err)
		}
	}
}

// 채팅 생성
func AddChat(c *gin.Context) {
	// Connect to the database
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}

	text := c.PostForm("ctext")
	ridStr := c.PostForm("rid")
	rid, err := strconv.ParseUint(ridStr, 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Invalid ID"})
		return
	}

	newChat := &model.Chat{Text: text, UID: uint(Loguser.ID), RID: uint(rid)}

	tx := db.Begin()
	if err := tx.Create(newChat).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	tx.Commit()

	// 데이터 생성 후 redirect
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/data/%d", newChat.RID))
	log.Printf("Redirecting to /data/%d\n", newChat.RID)
}

func GetRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func GetLoginForm(c *gin.Context) {
	Logout()
	c.HTML(200, "clientLogin.html", gin.H{})
}

func GetUnregisteredForm(c *gin.Context) {
	c.HTML(http.StatusOK, "unregistered-users.html", gin.H{})
}

func LoginHandler(c *gin.Context) {
	// Get the user's username and password from the request
	name := c.PostForm("name")
	password := c.PostForm("password")

	user, flag := CheckUser(name, password)
	// 유저 데이터가 없을 시, unregistered-user.html으로 리디렉션 시켜줘라.
	// 여기에는 register.html로 이동할 수 있는 링크와, 다시 login.html으로 이동할 수 있는 링크 삽입. << 완료.
	if flag == false {
		c.HTML(http.StatusOK, "unregistered-user.html", nil)
		return
	}

	// Render the login success template with the user's name
	log.Printf("Rendering loginSuccess.html template for user %s", user.Name)

	c.HTML(http.StatusOK, "loginSuccess.html", gin.H{"message": "Login successful"})
}
func Logout() {
	Loguser = nil
}

func CheckUser(name, password string) (*model.User, bool) {
	// Connect to the database
	DB, _ := model.ConnectDatabase()
	sqlDB, err := DB.DB()
	defer sqlDB.Close()

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Retrieve the user data from the database
	var user model.User
	record := DB.Where(&model.User{Name: name}).Where(&model.User{Password: password}).First(&user)

	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		fmt.Println("Record not found", record.Error)
		return nil, false
	}
	if record.RowsAffected > 0 {
		Loguser = &user
		fmt.Println(Loguser.Name)
		fmt.Println("login success.")
		// Return the authenticated user
		return &user, true
	}
	fmt.Println("login failed.", Loguser.Name)
	return nil, false
}

func Posting(title, content string) error {
	// Connect to the database
	DB, _ := model.ConnectDatabase()
	sqlDB, err := DB.DB()
	defer sqlDB.Close()
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
	}
	var board model.Board
	if err := DB.Where("id = ?", 1).First(&board).Error; err != nil {
		// user already exists, redirect to the login page
		board1 := model.Board{Title: "게시판"}
		DB.Create(&board1)
	}
	post := model.Post{BoardID: 1, UserID: Loguser.ID, Title: title, Content: content}
	result := DB.Create(&post)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func Comment(pid uint, content string) error {
	// Connect to the database
	DB, _ := model.ConnectDatabase()
	sqlDB, err := DB.DB()
	defer sqlDB.Close()
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
	}

	comment := model.Comment{UserID: Loguser.ID, PostID: pid, Content: content}
	result := DB.Create(&comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 회원가입
func PostRegister(c *gin.Context) {
	// Connect to the database
	db, err := model.ConnectDatabase()
	sqlDB, _ := db.DB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"errorMessage": "Internal server error",
		})
		return
	}
	defer sqlDB.Close()

	// Receive from the website
	name := c.PostForm("name")
	password := c.PostForm("password")
	dongIDs := c.PostForm("address")
	nickname := c.PostForm("nickname")
	email := c.PostForm("email")

	// check if the user already exists
	var user model.User
	if err := db.Where("name = ?", name).First(&user).Error; err == nil {
		// user already exists, redirect to the login page
		c.HTML(http.StatusOK, "registerFail.html", gin.H{})
		return
	}

	// Create an anonymous user
	user = model.User{
		Name:     name,
		Password: password,
		DongIDs:  dongIDs,
		Nickname: nickname,
		Email:    email,
	}

	// create a new user
	if err := db.Create(&user).Error; err != nil {
		log.Println("error creating user :", err)
		c.HTML(http.StatusInternalServerError, "errorCreateUser.html", gin.H{
			"errorMessage": "Internal server error",
		})
		return
	}

	// set up associations if user has more than one dong
	if len(user.DongIDs) > 1 {
		if err := model.SetupAssociations(db, &user); err != nil {
			c.HTML(http.StatusInternalServerError, "errorCreateUser.html", gin.H{
				"errorMessage": "Internal server error",
			})
			return
		}
	}

	// redirect to the success page
	c.HTML(http.StatusOK, "registerSuccess.html", gin.H{})
}

func ShowBoardList(c *gin.Context) {
	fmt.Println("ShowBoardList executed")
	// Connect to the database
	DB, _ := model.ConnectDatabase()
	sqlDB, err := DB.DB()
	defer sqlDB.Close()
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
	}

	var boardList []model.Post
	err2 := DB.Find(&boardList).Error
	if err2 != nil {
		log.Fatalf("%v", err2)
	}

	tmpl, err := template.ParseFiles("view/templates/board.html")
	if err != nil {
		fmt.Println("error while parsing board.html", err)
		log.Fatalf("%v", err)
	}

	err = tmpl.Execute(c.Writer, gin.H{
		"posts": boardList,
	})
	if err != nil {
		fmt.Println("error while executing board.html", err)
		log.Fatalf("%v", err)
	}

}
func GetUsers(c *gin.Context) {
	// Connect to the database
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer sqlDB.Close()

	// Retrieve all user data from the database
	var users []model.User
	errs := db.Find(&users).Error
	if errs != nil {
		fmt.Println("No users found")
		log.Fatalf("%v", errs)
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("view/templates/users.html")
	if err != nil {
		fmt.Println("Failed to parse Files")
		log.Fatalf("%v", err)
	}

	// Execute the template with the user data
	err = tmpl.Execute(c.Writer, gin.H{
		"users": users,
	})
	if err != nil {
		fmt.Println("Failed to execute Template")
		log.Fatalf("%v", err)
	}
}

// 헬스체크

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// GET
// 데이터 조회하기 //

func GetCities(c *gin.Context) {
	// Connect to the database
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	defer sqlDB.Close()

	// Retrieve all cities data from the database
	var city model.City
	cities, err := model.SearchData(db, &city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the City data to the user
	c.JSON(http.StatusOK, gin.H{"cities": cities})
}
func GetDistricts(c *gin.Context) {
	// Connect to the database
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	defer sqlDB.Close()

	// Retrieve all Districts data from the database
	var dis model.District
	diss, err := model.SearchData(db, &dis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the District data to the user
	c.JSON(http.StatusOK, gin.H{"districts": diss})
}
func GetDongs(c *gin.Context) {
	// Connect to the database
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	defer sqlDB.Close()

	// Retrieve all User data from the database
	var dong model.Dong
	dongs, err := model.SearchData(db, &dong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the Dong data to the user
	c.JSON(http.StatusOK, gin.H{"dongs": dongs})
}

// POST
// DB에 데이터 직접 넣기 //

func PostCities(c *gin.Context) {
	// Connect to the database
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	defer sqlDB.Close()

	// Bind the incoming JSON data to a City struct
	var cities model.City
	if err := c.ShouldBindJSON(&cities); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the City data into the database
	errs := model.InsertData(db, &cities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error()})
		return
	}

	// Return a success message to the user
	c.JSON(http.StatusOK, gin.H{"message": "City created successfully"})
}

func PostDistricts(c *gin.Context) {
	// Connect to the database
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	defer sqlDB.Close()

	// Bind the incoming JSON data to a District struct
	var districts model.District
	if err := c.ShouldBindJSON(&districts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the District into the database
	errs := model.InsertData(db, &districts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error()})
		return
	}

	// Return a success message to the user
	c.JSON(http.StatusOK, gin.H{"message": "District created successfully"})
}

func PostDongs(c *gin.Context) {
	// Connect to the database
	db, _ := model.ConnectDatabase()
	sqlDB, err := db.DB()
	defer sqlDB.Close()

	// Bind the incoming JSON data to a Dong struct
	var dongs model.Dong
	if err := c.ShouldBindJSON(&dongs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the Dong data into the database
	errs := model.InsertData(db, &dongs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error()})
		return
	}

	// Return a success message to the user
	c.JSON(http.StatusOK, gin.H{"message": "Dong created successfully"})
}

func PostUsers(c *gin.Context) {
	// Connect to the database
	DB, _ := model.ConnectDatabase()
	sqlDB, err1 := DB.DB()
	defer sqlDB.Close()
	if err1 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err1.Error()})
		return
	}

	// Bind the incoming JSON data to a User struct
	var user model.User
	if err2 := c.ShouldBindJSON(&user); err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}

	// Insert the User data into the database
	err3 := model.InsertData(DB, &user)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err3.Error()})
		return
	}

	// Return a success message to the user
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func PostProductsStatus(c *gin.Context) {
	// Connect to the database
	DB, _ := model.ConnectDatabase()
	sqlDB, err1 := DB.DB()
	defer sqlDB.Close()

	if err1 != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "database error"})
		c.JSON(http.StatusNotFound, gin.H{"error": err1.Error()})
		return
	}

	// Bind the incoming JSON data to a User struct
	var status model.ProductStatus
	if err2 := c.ShouldBindJSON(&status); err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Binding error"})
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}

	// Insert the User data into the database
	err3 := model.InsertData(DB, &status)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Inserting error"})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err3.Error()})
		return
	}

	// Return a success message to the user
	c.JSON(http.StatusOK, gin.H{"message": "ProductStatus created successfully"})
}

func PostProductsCategory(c *gin.Context) {
	// Connect to the database
	DB, _ := model.ConnectDatabase()
	sqlDB, err1 := DB.DB()
	defer sqlDB.Close()
	if err1 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err1.Error()})
		return
	}

	// Bind the incoming JSON data to a User struct
	var category model.ProductCategory
	if err2 := c.ShouldBindJSON(&category); err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}

	// Insert the User data into the database
	err3 := model.InsertData(DB, &category)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err3.Error()})
		return
	}

	// Return a success message to the user
	c.JSON(http.StatusOK, gin.H{"message": "ProductCategory created successfully"})
}

func PostProducts(c *gin.Context) {
	// Connect to the database
	DB, _ := model.ConnectDatabase()
	sqlDB, err1 := DB.DB()
	defer sqlDB.Close()
	if err1 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err1.Error()})
		return
	}

	// Bind the incoming JSON data to a User struct
	var products model.Product
	if err2 := c.ShouldBindJSON(&products); err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}

	// Insert the User data into the database
	err3 := model.InsertData(DB, &products)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err3.Error()})
		return
	}

	// Return a success message to the user
	c.JSON(http.StatusOK, gin.H{"message": "Products created successfully"})
}
