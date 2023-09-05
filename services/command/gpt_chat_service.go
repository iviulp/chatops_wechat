package command

import (
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"yuguosheng/int/mychatops/config"
	"yuguosheng/int/mychatops/dao"
	"yuguosheng/int/mychatops/middleware"
	"yuguosheng/int/mychatops/services/wx"
	"yuguosheng/int/mychatops/to"
	"yuguosheng/int/mychatops/utils/wecom"
)

type MyOPSCommand struct {
	Command string
}

func NewCommand() *MyOPSCommand {
	return &MyOPSCommand{}
}

func (c MyOPSCommand) MyExec(userData to.MsgContent) bool {
	// 检查是否请求过相同内容 不存在调用openai

	respOpenAI := MyCompareCacheAndGetFromApi(userData)

	// 发送到微信
	mode := config.GetSystemConf().MsgMode
	switch mode {
	case "markdown":
		return wx.SendToWxByMarkdown(userData, respOpenAI)
	case "text":
		return wx.SendToWxByText(userData, respOpenAI)
	default:
		return false
	}

}

// 我的新增
func MyCompareCacheAndGetFromApi(data to.MsgContent) (respOpenAI string) {
	fmt.Println(data.Content)

	pattern := `^shell: `
	re := regexp.MustCompile(pattern)
	if re.MatchString(data.Content) {
		if data.FromUsername == "YuGuoSheng" {
			middleware.MyLogger.Info("已经进入shell:命令", zap.Any("命令", data.Content))
			cmdtring := strings.TrimPrefix(data.Content, `shell: `)
			respOpenAI = "> 执行的命令是:\n" + cmdtring + "\n" + "> 执行的结果是:\n" + CommandRun(cmdtring) + "\n>>>>> " + time.Now().Format("2006-01-02 15:04:05")
			dao.InsertUserContext(data.FromUsername, cmdtring, respOpenAI)
			return respOpenAI
		}
	}

	pattern = `^file: `
	re = regexp.MustCompile(pattern)
	if re.MatchString(data.Content) {
		if data.FromUsername == "YuGuoSheng" {
			middleware.MyLogger.Info("已经进入file:命令", zap.Any("命令", data.Content))
			cmdtring := strings.TrimPrefix(data.Content, `file: `)
			file, _ := os.Open(cmdtring)
			defer file.Close()
			filecontents, _ := ioutil.ReadAll(file)
			wecom.SendFileToUser(filecontents, "temp.txt", data.FromUsername)
			respOpenAI = "请接收！"
			return respOpenAI
		}

	}

	cmdtring, err := dao.RegixMyCommand(data.Content, data.FromUsername)
	if err != nil {
		dao.InsertUserContext(data.FromUsername, cmdtring, err.Error())
		respOpenAI := err.Error() + "\n>> " + time.Now().Format("2006-01-02 15:04:05")
		return respOpenAI
	}
	middleware.MyLogger.Info("已经进入常规正则命令", zap.Any("命令", cmdtring))
	respOpenAI = "> 执行的命令是:\n" + cmdtring + "\n" + "> 执行的结果是:\n" + CommandRun(cmdtring) + "\n>>>>> " + time.Now().Format("2006-01-02 15:04:05")
	dao.InsertUserContext(data.FromUsername, cmdtring, respOpenAI)

	return respOpenAI
}

func CommandRun(cmdstring string) string {

	command := exec.Command("/bin/bash", "-c", cmdstring)
	output, err := command.CombinedOutput()
	if err != nil {
		return string(output)
	} else {
		return string(output)
	}

}
