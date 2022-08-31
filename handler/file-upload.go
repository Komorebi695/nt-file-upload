package handler

import (
	"fmt"
	"log"
	"net/http"
	"ntfileupload/db"
	"ntfileupload/service"
	"ntfileupload/service/oss"
	"ntfileupload/util"
	"strings"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// 参数
	name := r.FormValue("nt_name")
	studentId := r.FormValue("nt_studentId")
	email := r.FormValue("nt_email")

	fmt.Println(name, studentId, email, "参数")

	// 检查是否有该用户
	userOk := db.CheckUser(studentId)
	// 没有该用户
	if !userOk {
		resp := util.GenSimpleRespStream(-1, "抱歉!你没有报名比赛")
		_, _ = w.Write(resp)
		return
	}

	// 检查用户是否上传过文件
	uploadOk := db.CheckUpload(studentId)
	if uploadOk {
		resp := util.GenSimpleRespStream(-1, "你已经上传过文件!")
		_, _ = w.Write(resp)
		return
	}

	// 获取文件
	file, head, err := r.FormFile("nt_file")
	if nil != err {
		fmt.Println("上传文件错误:", err.Error())
		return
	}
	defer file.Close()

	fmt.Println(head.Filename)

	temp := strings.Split(head.Filename, ".")
	fileFormat := temp[1]

	ossPath := "nt/" + name + studentId + "." + fileFormat

	fd, err := head.Open()

	// 上传到 oss
	err = oss.Bucket().PutObject(ossPath, fd)
	if err != nil {
		log.Println("oss Put失败!")
		resp := util.GenSimpleRespStream(-1, "上传失败!")
		_, _ = w.Write(resp)
		return
	}

	// 上传成功,改变文件上传状态
	db.UpdateStatus(studentId)

	// 上传成功，发送邮件通知。
	service.SendEmail(email, "2022 NT挑战赛", name+"同学,你好! 你的作品已经提交成功啦~~~")

	resp := util.RespMsg{
		Code: 0,
		Msg:  "上传成功!",
		Data: "",
	}
	_, _ = w.Write(resp.JSONBytes())
	return
}
