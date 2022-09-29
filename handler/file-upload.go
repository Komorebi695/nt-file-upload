package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"ntfileupload/db"
	"ntfileupload/service"
	"ntfileupload/service/oss"
	"ntfileupload/util"
	"strings"
)

func Home(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./view/index.html")
	if err != nil {
		_, _ = io.WriteString(w, "内部服务器错误")
		return
	}
	_, _ = io.WriteString(w, string(data))
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// 参数
	name := r.FormValue("nt_name")
	studentId := r.FormValue("nt_studentId")
	email := r.FormValue("nt_email")

	fmt.Println(name, studentId, email, "参数")

	// 获取文件
	file, head, err := r.FormFile("nt_file")
	if nil != err {
		fmt.Println("上传文件错误:", err.Error())
		resp := util.GenSimpleRespStream(-1, "上传文件错误!")
		_, _ = w.Write(resp)
		return
	}
	defer file.Close()

	//fmt.Println(head.Filename)

	temp := strings.Split(head.Filename, ".")
	fileFormat := temp[1]

	// 判断文件格式是否符合
	if fileFormat != "zip" && fileFormat != "rar" && fileFormat != "7z" {
		resp := util.GenSimpleRespStream(1, "文件格式错误!")
		_, _ = w.Write(resp)
		return
	}

	// 检查是否有该用户
	userOk := db.CheckUser(studentId)
	// 没有该用户
	if !userOk {
		resp := util.GenSimpleRespStream(1, "抱歉!你没有报名比赛")
		_, _ = w.Write(resp)
		return
	}

	// 检查用户是否上传过文件
	uploadOk := db.CheckUpload(studentId)
	if uploadOk {
		resp := util.GenSimpleRespStream(1, "你已经上传过文件!")
		_, _ = w.Write(resp)
		return
	}

	// oos 中的路径
	ossPath := "nt/" + name + studentId + "." + fileFormat

	// 获得上传的文件
	fd, err := head.Open()

	// 上传到 oss
	err = oss.Bucket().PutObject(ossPath, fd)
	if err != nil {
		log.Println("oss Put失败!")
		resp := util.GenSimpleRespStream(1, "上传失败!")
		_, _ = w.Write(resp)
		return
	}

	// 上传成功,改变文件上传状态
	db.UpdateStatus(studentId)

	go func() {
		// 上传成功，发送邮件通知。
		to := []string{email}
		cc := []string{}
		bcc := []string{}

		err1 := service.SendToMail("2022 NT挑战赛", name+"同学，你好! 你的作品已经提交成功啦~~~\n请加入QQ群:615848077 及时获取答辩相关信息哦!", "text", "", to, cc, bcc)
		if err1 != nil {
			fmt.Println(err1.Error())
		}
	}()

	resp := util.GenSimpleRespStream(0, "上传成功!")
	_, _ = w.Write(resp)
	return
}
