package db

import (
	"database/sql"
	"log"
	mydb "ntfileupload/db/mysql"
)

// User 用户表
type User struct {
	UserName     string
	StudentId    string
	Email        string
	UploadStatus int
}

// CheckUpload 获取用户文件上传状态
func CheckUpload(studentId string) bool {
	stmt, err := mydb.DBConn().Prepare("select `status` from user where `student_id`=?")
	if nil != err {
		log.Println("db error:", err.Error())
		return false
	}
	defer stmt.Close()

	var status int

	// 执行
	err = stmt.QueryRow(studentId).Scan(&status)
	if nil != err {
		log.Println("select error: ", err.Error())
		return false
	}

	if status == 0 {
		return false
	}
	return true
}

// CheckUser 检查用户是否存在
func CheckUser(studentId string) bool {
	stmt, err := mydb.DBConn().Prepare("select `id` from user where `student_id`=?")
	if nil != err {
		log.Println("db error:", err.Error())
		return false
	}
	defer stmt.Close()

	var id int

	err = stmt.QueryRow(studentId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到数据
			return false
		}
		log.Println("select error: ", err.Error())
		return false
	}

	return true
}

// UpdateStatus 改变用户文件上传状态
func UpdateStatus(studentId string) bool {
	stmt, err := mydb.DBConn().Prepare("update `user` set  `status` =? where `student_id`=?")
	if nil != err {
		log.Println("db error:", err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(1, studentId)
	if nil != err {
		log.Println("update error: ", err.Error())
		return false
	}
	return true
}
