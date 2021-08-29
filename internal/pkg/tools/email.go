package tools

import (
	"fmt"

	"github.com/astaxie/beego/utils"
	"github.com/tuplz/tuplz-be/config"
)

func SendEmail(target []string, data string) {
	serv := `{"username":"` + config.EmailUser + `", "password":"` + config.EmailAuth +
		`", "host":"` + config.EmailHost + `", "port":` + config.EmailPort + `}`
	email := utils.NewEMail(serv)
	email.To = target
	email.From = config.EmailUser
	email.Subject = "Verification"
	email.Text = "Welcome to our tuplz system!"
	email.HTML = "<h1>" + data + "</h1>"
	// email.AttachFile("1.jpg") // 附件
	// email.AttachFile("1.jpg", "1") // 内嵌资源
	err := email.Send()
	if err != nil {
		fmt.Println(err)
		return
	}
}
