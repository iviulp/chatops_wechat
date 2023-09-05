package controller

import (
	"encoding/xml"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"yuguosheng/int/mychatops/dao"
	"yuguosheng/int/mychatops/middleware"
	"yuguosheng/int/mychatops/services"
	"yuguosheng/int/mychatops/to"
	"yuguosheng/int/mychatops/utils/wecom"
)

// VerifyCallBack 回调验证
func VerifyCallBack(c *gin.Context) {
	var q to.CallBackParams
	if err := c.Bind(&q); err != nil {
		middleware.MyLogger.Error("绑定回调Query错误", zap.Any("Error", err))
	}
	msg := wecom.GetReVerifyCallBack(q)
	_, _ = c.Writer.Write(msg)
}

// WxChatCommand 实际处理用户消息
func WxChatCommand(c *gin.Context) {
	var dataStuc to.CallBackData
	if err := c.ShouldBindQuery(&dataStuc); err != nil {
		middleware.MyLogger.Error("绑定回调Query错误", zap.Any("Error", err))
	}
	// 解析请求体
	raw, err := c.GetRawData()
	if err != nil {
		middleware.MyLogger.Error("解析微信回调参数失败", zap.Any("Error", err))
		return
	}
	userData := to.MsgContent{}
	userDataDecrypt := wecom.DeCryptMsg(raw, dataStuc.MsgSignature, dataStuc.TimeStamp, dataStuc.Nonce)
	// 解密失败返回空
	if userDataDecrypt == nil {
		middleware.MyLogger.Error("解密失败", zap.Any("用户数据", userData))
	}
	// 提前向微信返回成功接受，防止微信多次回调
	c.JSON(http.StatusOK, "")
	// 异步处理用户请求
	go func() {
		err = xml.Unmarshal(userDataDecrypt, &userData)
		if err != nil {
			middleware.MyLogger.Error("反序列化用户数据错误", zap.Any("Error", err))
			return
		}
		// 检查用户是否存在，不存在创建
		if !dao.CheckUserAndCreate(userData.FromUsername) {
			middleware.MyLogger.Error("创建用户失败", zap.Any("用户信息", userData.FromUsername))
			return
		}
		// 分发消息类型 进行处理
		switch userData.MsgType {
		case models.CALLBACK_MSG_TYPE_TEXT:
			// 处理text消息
			services.DoTextMsg(userData)
			//case models.CALLBACK_MSG_TYPE_EVENT:
			//	// 处理事件消息
			//	services.DoEventMsg(userDataDecrypt)
		}
	}()
}
