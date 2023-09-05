package services

import (
	"go.uber.org/zap"
	"yuguosheng/int/mychatops/middleware"

	"yuguosheng/int/mychatops/services/command"
	"yuguosheng/int/mychatops/to"
)

// DoTextMsg
// @Description: Text消息逻辑
// @param cmd
func DoTextMsg(userData to.MsgContent) {
	ok := false

	//ok = command.NewGPTCommand().Exec(userData)
	ok = command.NewCommand().MyExec(userData)

	if !ok {
		middleware.MyLogger.Error("执行指令失败！", zap.Any("data", userData))
	}
	middleware.MyLogger.Info("执行指令成功!", zap.Any("data", userData))
	return
}
