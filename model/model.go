package model

import "github.com/jinzhu/gorm"

// 지역 기반 테이블 모델링
// cf) 보통 필드 뒤에 태그로 포멧이나 외래키 관계성을 명시해주는것이 더 유용하기 때문에 관습적으로 쓰이고 있음

// 중요! gorm.Model의 ID 필드의 type은 uint64임. uint (X)

type City struct {
	gorm.Model
	Name string `json:"name"`
}

type District struct {
	gorm.Model
	CityID uint64 `json:"city_id" gorm:"foreignkey:CityID"`
	City   City   `json:"city" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name   string `json:"name"`
}

type Dong struct {
	gorm.Model
	DistrictID uint64   `json:"district_id" gorm:"foreignkey:DistrictID"`
	District   District `json:"district" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name       string   `json:"name"`
	UserIDs    []int64  `json:"user_ids" gorm:"type:longtext"`
	Users      []User   `json:"users" gorm:"many2many:user_dongs;"`
}

type User struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null;unique"`
	Password    string  `json:"password" gorm:"not null"`
	Nickname    string  `json:"nickname" gorm:"not null;unique"`
	Email       string  `json:"email"`
	Temperature float64 `json:"temperature" gorm:"not null;default:36.5"`

	// Foreign keys
	DongIDs string `json:"dong_ids" gorm:"type:longtext"`

	// Foreign keys deliver
	Dongs    []Dong    `json:"dongs" gorm:"many2many:user_dongs"`
	Products []Product `gorm:"foreignKey:UID"`
	ChRooms  []ChRoom  `gorm:"foreignKey:BID"`
	Chats    []Chat    `gorm:"foreignKey:UID"`
	Posts    []Post
	Comments []Comment
	Cocs     []CoC
}

// 다대다 관계로 생성되는 Table
type UserDong struct {
	UserID uint
	DongID uint
}

////////////////////////////////////////////////////////////////////////
//////////////////////////// SEPARATOR /////////////////////////////////
////////////////////////////////////////////////////////////////////////

// 물품
type Product struct {
	gorm.Model
	Pname string `json:"pname" gorm:"not null"`
	Price int    `json:"price" gorm:"not null;default:0"`
	Ptext string `json:"ptext"`

	//외래키
	UID  uint `json:"user_id" gorm:"not null"`
	PCID uint `json:"product_category_id" gorm:"not null"`
	PSID uint `json:"product_status_id" gorm:"not null"`

	//외래키 전달
	ChRooms []ChRoom `gorm:"foreignKey:PID"`
	ChRoomU []ChRoom `gorm:"foreignKey:SID;references:UID"`
}

// 물품 분류
type ProductCategory struct {
	gorm.Model
	Cname string `json:"cname" gorm:"not null;unique"`

	//외래키 전달
	Products []Product `gorm:"foreignKey:PCID"`
}

// 물품 거래 상태
type ProductStatus struct {
	gorm.Model
	Sname string `json:"sname" gorm:"not null;unique"`

	//외래키 전달
	Products []Product `gorm:"foreignKey:PSID"`
}

// 채팅룸
type ChRoom struct {
	gorm.Model
	//외래키
	PID uint `json:"product_id" gorm:"not null"`
	BID uint `json:"buyer_id" gorm:"not null"`
	SID uint `json:"seller_id" gorm:"not null"`
	//외래키 전달
	Chats []Chat `gorm:"foreignKey:RID"`
}

// 채팅
type Chat struct {
	gorm.Model
	Text string `json:"text" gorm:"not null"`
	//외래키
	RID uint `json:"room_id" gorm:"not null"`
	UID uint `json:"user_id" gorm:"not null"`
}
type Board struct {
	gorm.Model
	Title string
	Posts []Post
}

type Post struct {
	gorm.Model
	BoardID  uint
	UserID   uint
	Title    string
	Content  string
	Comments []Comment
}

type Comment struct {
	gorm.Model
	UserID    uint
	Content   string
	PostID    uint
	Commments []CoC
}

type CoC struct {
	gorm.Model
	UserID    uint
	CommentID uint
	Content   string
}
