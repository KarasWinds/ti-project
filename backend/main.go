package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Member 結構，代表客戶
type Member struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Username   string    `json:"username" binding:"required"`
	CreateTime time.Time `json:"create_time"`
}

// BorrowFee 結構，代表交易紀錄
type BorrowFee struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	MemberFK   uint      `json:"member_fk"`
	Type       int       `json:"type"`
	BorrowFee  float64   `json:"borrow_fee"`
	CreateTime time.Time `json:"create_time"`
}

type MemberTotalData struct {
	MemberID   uint      `json:"member_id"`
	Username   string    `json:"username"`
	CreateTime time.Time `json:"create_time"`
	TotalFee   float64   `json:"total_fee"`
}

// getEnv 用於獲取環境變數，如果該環境變數不存在或為空，則返回預設值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// 初始化資料庫連接
func InitDB() *gorm.DB {
	username := getEnv("DB_USERNAME", "root")
	password := getEnv("DB_PASSWORD", "root")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_NAME", "default")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	db.AutoMigrate(&Member{}, &BorrowFee{})

	return db
}

func main() {
	router := gin.Default()
	db := InitDB()

	// 獲取所有客戶資訊及過去一年的交易金額
	router.GET("/api/members", func(c *gin.Context) {
		var memberTotalData []MemberTotalData
		result := db.Table("borrow_fees").
			Select("members.id as member_id, members.username, members.create_time, SUM(borrow_fees.borrow_fee) as total_fee").
			Joins("left join members on members.id = borrow_fees.member_fk").
			Where("borrow_fees.create_time > ?", time.Now().AddDate(-1, 0, 0)).
			Group("members.id, members.username, members.create_time").
			Find(&memberTotalData)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, memberTotalData)
	})

	// 修改指定客戶資料
	router.PUT("/api/member/:id", func(c *gin.Context) {
		var member Member
		id := c.Param("id")

		if err := db.First(&member, id).Error; err != nil {
			log.Println("Member not found:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}

		if err := c.ShouldBindJSON(&member); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db.Save(&member)
		c.JSON(http.StatusOK, gin.H{"success": "Member updated"})
	})

	// 新增客戶
	router.POST("/api/member", func(c *gin.Context) {
		var member Member
		if err := c.ShouldBindJSON(&member); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db.Create(&member)
		c.JSON(http.StatusCreated, gin.H{"success": "Member created", "id": member.ID})
	})

	// 列出指定客戶指定時間內的交易紀錄
	router.GET("/api/member/:id/transactions", func(c *gin.Context) {
		var transactions []BorrowFee
		memberId := c.Param("id")
		startDate := c.DefaultQuery("start", "2023-01-01")
		endDate := c.DefaultQuery("end", "2023-12-31")

		if err := db.Where("member_fk = ? AND create_time BETWEEN ? AND ?", memberId, startDate, endDate).Find(&transactions).Error; err != nil {
			log.Println("Error fetching transactions:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, transactions)
	})

	router.Run(":8080")
}
