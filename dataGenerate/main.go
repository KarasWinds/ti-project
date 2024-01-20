package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
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

const (
	numMembers    = 1000 // 總會員數
	numBorrowFees = 5000 // 總費用記錄數
)

func main() {
	username := getEnv("DB_USERNAME", "root")
	password := getEnv("DB_PASSWORD", "root")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_NAME", "default")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)

	// 資料庫連接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		os.Exit(1)
	}

	db.AutoMigrate(&Member{}, &BorrowFee{})
	// 產生會員和費用資料
	memberCreateTimeMap := generateMembers(db)
	generateBorrowFees(db, memberCreateTimeMap)

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database")
	}

	err = sqlDB.Close()
	if err != nil {
		panic("failed to close database")
	}
}

// getEnv 用於獲取環境變數，如果該環境變數不存在或為空，則返回預設值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// generateMembers 產生隨機會員資料
func generateMembers(db *gorm.DB) map[int]string {
	memberCreateTimeMap := make(map[int]string)

	for i := 0; i < numMembers; i++ {
		member := Member{
			Username:   randomString(16),
			CreateTime: randomPastDate(5),
		}

		db.Create(&member)

		memberCreateTimeMap[i+1] = member.CreateTime.Format("2006-01-02 15:04:05")
	}

	fmt.Println("Members generated")
	return memberCreateTimeMap
}

// generateBorrowFees 產生費用記錄
func generateBorrowFees(db *gorm.DB, memberCreateTimeMap map[int]string) {
	lastCreateTime := time.Now().AddDate(0, -18, 0) // 從 18 個月前開始
	memberBorrowCount := make(map[int]int)

	for i := 0; i < numBorrowFees; i++ {
		memberID := rand.Intn(numMembers) + 1
		memberBorrowCount[memberID]++
		memberCreateTime, _ := time.Parse("2006-01-02 15:04:05", memberCreateTimeMap[memberID])
		if memberCreateTime.After(lastCreateTime) {
			lastCreateTime = memberCreateTime
		}

		increment := time.Duration(rand.Intn(7200)) * time.Second // 隨機增加最多 2 小時
		if newTime := lastCreateTime.Add(increment); newTime.Before(time.Now()) {
			lastCreateTime = newTime
		}

		borrowFee := BorrowFee{
			MemberFK:   uint(memberID),
			Type:       memberBorrowCount[memberID],
			BorrowFee:  rand.Float64() * 10000,
			CreateTime: lastCreateTime,
		}

		db.Create(&borrowFee)
	}
	fmt.Println("Borrow fees generated")
}

// randomString 產生隨機字串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// randomPastDate 產生過去日期
func randomPastDate(yearsBack int) time.Time {
	now := time.Now()
	years := rand.Intn(yearsBack) + 1
	months := rand.Intn(12)
	days := rand.Intn(30)
	return now.AddDate(-years, -months, -days)
}
