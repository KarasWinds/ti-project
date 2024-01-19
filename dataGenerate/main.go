package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

const (
	numMembers    = 1000 // 總會員數
	numBorrowFees = 5000 // 總費用記錄數
)

func main() {
	// 資料庫連接
	db, err := sql.Open("mysql", "test:test@/ti")
	if err != nil {
		panic(err)
	}
	defer db.Close() // 確保在函數退出時關閉資料庫連接

	// 產生會員和費用資料
	memberCreateTimeMap := generateMembers(db)
	generateBorrowFees(db, memberCreateTimeMap)
}

// generateMembers 產生隨機會員資料
func generateMembers(db *sql.DB) map[int]string {
	stmt, err := db.Prepare("INSERT INTO member (username, create_time) VALUES (?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	memberCreateTimeMap := make(map[int]string)
	for i := 0; i < numMembers; i++ {
		username := randomString(16) // 產生隨機用戶名
		createTime := randomPastDate(5).Format("2006-01-02 15:04:05")
		memberCreateTimeMap[i+1] = createTime

		_, err := stmt.Exec(username, createTime)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Members generated")
	return memberCreateTimeMap
}

// generateBorrowFees 產生費用記錄
func generateBorrowFees(db *sql.DB, memberCreateTimeMap map[int]string) {
	stmt, err := db.Prepare("INSERT INTO borrow_fee (member_fk, type, borrow_fee, create_time) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	lastCreateTime := time.Now().AddDate(0, -18, 0) // 從 18 個月前開始
	memberBorrowCount := make(map[int]int)

	for i := 0; i < numBorrowFees; i++ {
		memberID := rand.Intn(numMembers) + 1
		memberBorrowCount[memberID]++
		borrowType := memberBorrowCount[memberID]
		memberCreateTime, _ := time.Parse("2006-01-02 15:04:05", memberCreateTimeMap[memberID])
		if memberCreateTime.After(lastCreateTime) {
			lastCreateTime = memberCreateTime
		}

		borrowFee := rand.Float64() * 10000
		increment := time.Duration(rand.Intn(7200)) * time.Second // 隨機增加最多 2 小時
		if newTime := lastCreateTime.Add(increment); newTime.Before(time.Now()) {
			lastCreateTime = newTime
		}

		_, err := stmt.Exec(memberID, borrowType, borrowFee, lastCreateTime.Format("2006-01-02 15:04:05"))
		if err != nil {
			panic(err)
		}
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
