package email

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type BodyData struct {
	SendTo        interface{} // 接受人
	EmailCcTo     []string    // 抄送人邮箱列表
	Subject       string      // 标题
	Classify      []int       // 通知类型
	Id            int         // 工单ID
	Title         string      // 工单标题
	Creator       string      // 工单创建人
	Priority      int         // 工单优先级
	PriorityValue string      // 工单优先级
	CreatedAt     string      // 工单创建时间
	Content       string      // 通知的内容
	Description   string      // 表格上面的描述信息
	ProcessId     int         // 流程ID
	Domain        string      // 域名地址
}

func (b *BodyData) ParsingTemplate() (err error) {
	// 读取模版数据
	var (
		buf bytes.Buffer
	)

	tmpl, err := template.ParseFiles("email.html")
	if err != nil {
		return
	}

	b.Domain = viper.GetString("settings.domain.url")
	err = tmpl.Execute(&buf, b)
	if err != nil {
		return
	}

	b.Content = buf.String()

	return
}

func server(mailTo []string, ccTo []string, subject, body string) error {
	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user":  viper.GetString("settings.email.user"),
		"pass":  viper.GetString("settings.email.pass"),
		"host":  viper.GetString("settings.email.host"),
		"port":  viper.GetString("settings.email.port"),
		"alias": viper.GetString("settings.email.alias"),
	}

	mailConn["user"] = "xxxxxx"
	mailConn["pass"] = "xxxx"
	mailConn["host"] = "xxxx"
	mailConn["port"] = "587"
	mailConn["alias"] = "xxxx"

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mailConn["user"], mailConn["alias"]))
	m.SetHeader("To", mailTo...)
	m.SetHeader("Cc", ccTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	fmt.Println(mailTo)
	println("-----------------------------------")
	fmt.Println(ccTo)
	println("-----------------------------------")
	fmt.Println(subject)
	println("-----------------------------------")
	fmt.Println(body)
	println("-----------------------------------")
	fmt.Println(mailConn)
	println("-----------------------------------")
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	err := d.DialAndSend(m)
	fmt.Println(err)
	println("-----------------------------------")
	return err

}

func ConfigSetup(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		fmt.Println(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}
}

func TestSendEmail(t *testing.T) {
	bodyData := &BodyData{Description: "测试邮件发送EOF", Title: "Test Sending Email",
		Creator: "xxxxxx", PriorityValue: "紧急", CreatedAt: "2021-12-31 22:22:22",
		Domain: "http://10.0.2.98", Id: 1, ProcessId: 1,
	}
	ConfigSetup("settings.yml")
	bodyData.ParsingTemplate()
	server([]string{"xxxxxx"}, []string{"xxxxxxx"}, "TEST EMIAL EOF", bodyData.Content)
}
