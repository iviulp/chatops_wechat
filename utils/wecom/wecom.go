package wecom

import (
	"context"
	"go.uber.org/zap"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/response"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
	"yuguosheng/int/mychatops/config"
	"yuguosheng/int/mychatops/middleware"
	"yuguosheng/int/mychatops/to"
	"yuguosheng/int/mychatops/utils/xstring"
	"yuguosheng/int/mychatops/xconst"
)

var wxCrypt *wxbizmsgcrypt.WXBizMsgCrypt

var WeComApp *work.Work

func init() {
	LoadWeComAppConf()
	LoadWxUtils()
}

func LoadWeComAppConf() {
	middleware.MyLogger.Info("初始化企业微信助手......")
	app, err := work.NewWork(&work.UserConfig{
		CorpID:  config.GetWechatConf().Corpid,     // 企业微信的app id，所有企业微信共用一个。
		AgentID: config.GetWechatConf().AgentId,    // 内部应用的app id
		Secret:  config.GetWechatConf().CorpSecret, // 内部应用的app secret
		OAuth: work.OAuth{
			Callback: config.GetSystemConf().CallBackUrl, //
			Scopes:   nil,
		},
		HttpDebug: false,
	})
	if err != nil {
		middleware.MyLogger.Fatal("初始化企业微信助手失败！", zap.Any("Error", err))
	}
	WeComApp = app
}

func LoadWxUtils() {
	middleware.MyLogger.Info("初始化微信工具包......")
	wxCrypt = wxbizmsgcrypt.NewWXBizMsgCrypt(config.GetWechatConf().WeApiRCallToken, config.GetWechatConf().WeApiEncodingKey, config.GetWechatConf().Corpid, wxbizmsgcrypt.XmlType)
}

// GetReVerifyCallBack 从微信回调解析请求数据
func GetReVerifyCallBack(q to.CallBackParams) []byte {
	msg, cryptErr := wxCrypt.VerifyURL(q.MsgSignature, q.TimeStamp, q.Nonce, q.Echostr)
	if cryptErr != nil {
		middleware.MyLogger.Error("验证Url出错（回调消息解密错误", zap.Any("cryptErr", cryptErr))
		return []byte("")
	}
	middleware.MyLogger.Info("解析的回调字符为：", zap.Any("解析的回调字符", string(msg)))
	return msg
}

// DeCryptMsg 解密消息
func DeCryptMsg(cryptMsg []byte, msgSignature, timeStamp, nonce string) []byte {
	msg, cryptErr := wxCrypt.DecryptMsg(msgSignature, timeStamp, nonce, cryptMsg)
	if cryptErr != nil {
		middleware.MyLogger.Error("回调消息解密错误", zap.Any("cryptErr", cryptErr))
		return nil
	}
	return msg
}

// CryptMessage 加密消息
func CryptMessage(respData, reqTimestamp, reqNonce string) string {
	encryptMsg, cryptErr := wxCrypt.EncryptMsg(respData, reqTimestamp, reqNonce)
	if cryptErr != nil {
		middleware.MyLogger.Error("消息加密错误", zap.Any("cryptErr", cryptErr))
		return ""
	}
	return string(encryptMsg)
}

// SendMarkdownToUSer 发送Markdown消息
func SendMarkdownToUSer(userName string, respMsg string) *response.ResponseMessageSend {
	middleware.MyLogger.Info("4444444收到的信息", zap.Any("信息", respMsg))
	// 封装微信消息体
	messages := &request.RequestMessageSendMarkdown{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:                 userName,
			ToParty:                "",
			ToTag:                  "",
			MsgType:                "markdown",
			AgentID:                config.GetWechatConf().AgentId,
			EnableDuplicateCheck:   1,
			DuplicateCheckInterval: 1800,
		},
		Markdown: &request.RequestMarkdown{
			Content: respMsg,
		},
	}
	// 发送微信消息
	resp, err := WeComApp.Message.SendMarkdown(context.Background(), messages)
	if err != nil {
		middleware.MyLogger.Error("创建微信发送消息内容失败", zap.Any("Error", err))
		return nil
	}
	return resp
}

// SendTextToUSer 发送text消息
func SendTextToUSer(userName string, respMsg string) *response.ResponseMessageSend {
	// 封装微信消息体
	messages := &request.RequestMessageSendText{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:                 userName,
			ToParty:                "",
			ToTag:                  "",
			MsgType:                "text",
			AgentID:                config.GetWechatConf().AgentId,
			Safe:                   0,
			EnableIDTrans:          0,
			EnableDuplicateCheck:   0,
			DuplicateCheckInterval: 1800,
		},
		Text: &request.RequestText{
			Content: respMsg,
		},
	}
	// 发送微信消息
	resp, err := WeComApp.Message.SendText(context.Background(), messages)
	if err != nil {
		middleware.MyLogger.Error("创建微信发送消息内容失败", zap.Any("Error", err))
		return nil
	}
	return resp
}

// SendImageToUser
// @Description: 发送制定二进制图片数据给用户
// @param data
// @param imageExt
// @param userName
// @return *response.ResponseMessageSend
func SendImageToUser(data []byte, imageExt string, userName string) *response.ResponseMessageSend {
	ctx := context.Background()
	dataFrom := &power.HashMap{
		"name":  userName + "_" + xstring.GenerateRandomStr() + imageExt,
		"value": data,
	}
	middleware.MyLogger.Info("上传微信图片素材中....", zap.Any("用户", userName))
	tempImageResp, err := WeComApp.Media.UploadTempImage(ctx, "", dataFrom)
	if err != nil {
		middleware.MyLogger.Error("上传临时图片素材失败", zap.Any("用户", userName), zap.Any("微信临时素材响应", tempImageResp))
		return nil
	}
	middleware.MyLogger.Info("上传微信素材成功，正在发送消息...", zap.Any("用户", userName))
	// 发送图片消息
	message := &request.RequestMessageSendImage{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:                 userName,
			ToParty:                "",
			ToTag:                  "",
			MsgType:                xconst.MSG_TYPE_IMAGE,
			AgentID:                config.GetWechatConf().AgentId,
			Safe:                   0,
			EnableIDTrans:          0,
			EnableDuplicateCheck:   0,
			DuplicateCheckInterval: 1800,
		},
		Image: &request.RequestImage{MediaID: tempImageResp.MediaID},
	}
	resp, err := WeComApp.Message.SendImage(ctx, message)
	if err != nil {
		middleware.MyLogger.Error("发送已上传图片素材失败", zap.Any("用户", userName), zap.Any("微信发送图片消息响应", resp))
		return nil
	}
	middleware.MyLogger.Info("微信图片消息推送成功！", zap.Any("用户", userName))
	return resp
}

// SendFileToUser
//
//	@Description: 发送文件给用户
//	@param data
//	@param fileExt
//	@param userName
//	@return *response.ResponseMessageSend
func SendFileToUser(data []byte, fileExt, userName string) *response.ResponseMessageSend {
	ctx := context.Background()
	// 上传临时文件获取mediaId
	dataFrom := &power.HashMap{
		"name":  userName + "_" + xstring.GenerateRandomStr() + "_contextHistory" + fileExt,
		"value": data,
	}
	middleware.MyLogger.Info("上传微信临时文件中....", zap.Any("用户", userName))
	tempImageResp, err := WeComApp.Media.UploadTempFile(ctx, "", dataFrom)
	if err != nil {
		middleware.MyLogger.Error("上传临时图片素材失败", zap.Any("用户", userName), zap.Any("微信临时素材响应", tempImageResp))
		return nil
	}
	// 发送文件消息
	messages := &request.RequestMessageSendFile{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:                 userName,
			ToParty:                "",
			ToTag:                  "",
			MsgType:                xconst.MSG_TYPE_FILE,
			AgentID:                config.GetWechatConf().AgentId,
			Safe:                   0,
			EnableDuplicateCheck:   0,
			DuplicateCheckInterval: 1800,
		},
		File: &request.RequestFile{
			MediaID: tempImageResp.MediaID,
		},
	}
	resp, err := WeComApp.Message.SendFile(ctx, messages)
	if err != nil {
		middleware.MyLogger.Error("发送已上传文件失败", zap.Any("用户", userName), zap.Any("微信发送文件消息响应", resp))
		return nil
	}
	middleware.MyLogger.Info("微信文件消息推送成功！", zap.Any("用户", userName))
	return resp
}
