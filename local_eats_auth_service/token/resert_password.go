package token

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

func SendCode(email string, code string) {
	from := "hatamovjamshid47@gmail.com"
	password := "xnis whmf zjbb dscw"

	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Your verification code \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Passwd string
	}{

		Passwd: code,
	})

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email sended to:", email)

}
