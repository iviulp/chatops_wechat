package xconst

const (
	USER_DAO_FIRST_CREATE = "用户第一次使用，创建用户中..."
	USER_DAO_SEARCH_ERR   = "查询用户数据失败"
	USER_DAO_INSERT_ERR   = "查询用户数据创建失败"
)

const (
	COMMAND_HELP              = "@help"
	COMMAN_GPT_DELETE_CONTEXT = "@clear"
	COMMAN_GPT_EXPORT         = "@export"
)

const (
	PROMPT_DEFAULT = "请全程使用中文与我对话"
)

// GetDefaultNoticeMenu
// @Description: 默认 提示消息
// @return string
func GetDefaultNoticeMenu() string {
	// 默认 提示消息
	return `这里是帮助菜单（如下是支持的菜单，以下不存在默认不进行处理）：
@help：帮助菜单 -> 例子：@help
@clear：清除聊天上下文 -> 例子：@clear
@export：导出你的本次对话内容 -> 例子：@export
使用建议：每次开启新对话前，可以先执行@clear清除上下文内容，使用@prompt-set设置角色（也可不设置）开启一次对话。`
}
