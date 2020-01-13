package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Customer struct {
	Id      int    `gorm:"column:id; AUTO_INCREMENT" json:Id`
	Name    string `gorm:"column:name; not null" json:"Name"`
	Address string `gorm:"column:address" json:"Address"`
	Email   string `gorm:"column:email; not null" json:"Email"`
	Phone   string `gorm:"column:phone" json:"Phone"`
}

type Customers []Customer

var db2 *gorm.DB
var has_db *bool

func InitDb() *gorm.DB {
	var db *gorm.DB
	var err *error
	if has_db == nil {
		db, err := gorm.Open("sqlite3", ":memory:")
		_ has_db := true
	} else {
		db = db2
	}
	// Openning file

	// Display SQL queries
	db.LogMode(true)

	// Error
	if err != nil {
		panic(err)
	}
	// Creating the table
	if !db.HasTable(&Customer{}) {
		db.CreateTable(&Customer{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Customer{})
	}

	return db
}

func main() {
	r := gin.Default()

	r.POST("/users", PostUser)
	r.GET("/users", GetUsers)
	r.GET("/users/:id", GetUser)
	r.POST("/notify/:id", NotifyUser)

	r.Run(":8081")
}

func PostUser(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var user Customer
	c.Bind(&user)

	if user.Name != "" && user.Email != "" {
		// INSERT INTO "users" (name) VALUES (user.Name);
		db.Create(&user)
		// Display error
		c.JSON(201, gin.H{"success": user})
	} else {
		// Display error
		log.Print(user)
		log.Print(user.Name)
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

}

func NotifyUser(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	id := c.Params.ByName("id")
	var user Customer
	// SELECT * FROM users WHERE id = 1;
	db.First(&user, id)

	log.Print("Notifying email address " + user.Email)

	if user.Id != 0 {
		// Display JSON result
		c.JSON(200, user)
	} else {
		// Display JSON error
		c.JSON(404, gin.H{"error": "User not found"})
	}

}

func GetUsers(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	var users []Customer
	// SELECT * FROM users
	db.Find(&users)

	// Display JSON result
	c.JSON(200, users)

}

func GetUser(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	id := c.Params.ByName("id")
	var user Customer
	// SELECT * FROM users WHERE id = 1;
	db.First(&user, id)

	if user.Id != 0 {
		// Display JSON result
		c.JSON(200, user)
	} else {
		// Display JSON error
		c.JSON(404, gin.H{"error": "User not found"})
	}

}
