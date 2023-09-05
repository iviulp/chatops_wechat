package dao

import (
	"errors"
	"go.uber.org/zap"
	"strings"
	"yuguosheng/int/mychatops/middleware"
	"yuguosheng/int/mychatops/xconst"
)

// SetContext
// @Description: 设置用户上下文对象中上下文
// @receiver c
// @param context
// @return bool
//
//	func (c *Context) SetContext(context model.MessageContext) bool {
//		data, err := json.Marshal(context)
//		if err != nil {
//			middleware.MyLogger.Error("存储上下文序列化错误", zap.Any("用户", c.Name))
//			return false
//		}
//		c.ContextMsg = string(data)
//		return true
//	}
//
// InsertUserContext
// @Description: 插入用户上下文
// @receiver c
// @param db
// @return error

func InsertUserContext(name, cmdtring, contextMsg string) error {
	var hiscontext Context
	hiscontext.Name = name
	hiscontext.Command = cmdtring
	hiscontext.Context_msg = contextMsg
	result := Db.Create(&hiscontext)

	// 检查插入是否成功
	if result.Error != nil {
		middleware.MyLogger.Error("Error inserting history:", zap.Any("Error", result.Error))
		return result.Error
	} else {
		middleware.MyLogger.Info("插入用户上下文成功")
		return nil
	}
}

func RegixMyCommand(cmdstring, name string) (mycomm string, err error) {
	var ccmd []ContextCommand
	result := Db.Where("text like ?", "%"+cmdstring+"%").Find(&ccmd)

	if result.Error != nil {
		middleware.MyLogger.Error("Error Getting user:", zap.Any("Error", result.Error))
		return "", errors.New(result.Error.Error())
	} else {
		middleware.MyLogger.Info("查询成功")
	}

	if len(ccmd) == 0 {
		return "", errors.New("还没有这个能力哦~~~")
	} else if len(ccmd) == 1 {
		allauthpeople := strings.Split(ccmd[0].Auth, ";")
		if IsInList(allauthpeople, name) {
			return ccmd[0].Command, nil
		} else if IsInList(allauthpeople, "all") {
			return ccmd[0].Command, nil
		} else {
			return "", errors.New("你没有这个执行的权限哦~~~")
		}

	} else {
		return "", errors.New("输入的有歧义哦~~，我不知道该执行哪一个了~~~")
	}

}

func IsInList(list []string, element string) bool {
	for _, v := range list {
		if v == element {
			return true
		}
	}
	return false
}

//func MyHoutaiZhixing(sqlstring string) (result string) {
//	query := fmt.Sprintf(`select command from contextcommand where text like ?;`)
//	stmt, err := dao.Db.Prepare(query)
//	if err != nil {
//		return "prepare数据库查询语句"
//	}
//	defer stmt.Close()
//	rows, err := stmt.Query("%" + sqlstring + "%")
//	if err != nil {
//		return "prepare数据库语句执行失败"
//	}
//	rowsCount := 0
//	for rows.Next() {
//		rowsCount++
//	}
//	if err = rows.Err(); err != nil {
//		return "rows统计失败"
//	}
//	if rowsCount == 1 {
//		var command string
//		if err = dao.Db.QueryRow("select command from contextcommand where text like ?;", "%"+sqlstring+"%").Scan(&command); err != nil {
//			return "查询失败"
//		}
//		return command
//	} else if rowsCount == 0 {
//		return "还没有这个能力哦~~~"
//	} else {
//		return "输入的有歧义哦~~，我不知道该执行哪一个了"
//	}
//	return "" + time.Now().String()
//}

//
// // UpdateContext
// // @Description: 更新用户上下文
// // @receiver c
// // @param Db
// // @return error
//
//	func UpdateContext(contextMsg, userName string) error {
//		query := `update context set context_msg=?,update_time=? where id in{
//					select id from context where name=? ORDER BY update_time limit 1}`
//		stmt, err := Db.Prepare(query)
//		if err != nil {
//			return errors.Wrap(err, "prepare数据库查询语句")
//		}
//		defer stmt.Close()
//		_, err = stmt.Exec(contextMsg, time.Now().UnixMilli(), userName)
//		if err != nil {
//			return errors.Wrap(err, "prepare数据库语句执行失败")
//		}
//		return nil
//	}
//
// // GetLatestUserContext
// // @Description:  获取最新用户上下文对象中上下文
// // @receiver c
// // @param Db
// // @return *model.MessageContext
// // @return error
//
//	func GetLatestUserContext(userName string) (*Context, error) {
//		sql := `select * from context where name=? Order By update_time DESC limit 1 `
//		stmt, err := Db.Prepare(sql)
//		if err != nil {
//			return nil, errors.Wrap(err, "prepare数据库查询语句")
//		}
//		defer stmt.Close()
//		rows, err := stmt.Query(userName)
//		defer rows.Close()
//		if err != nil {
//			return nil, errors.Wrap(err, "prepare数据库语句执行失败")
//		}
//		if rows.Next() {
//			temp := &Context{}
//			err = rows.Scan(&temp.Id, &temp.Name, &temp.ContextMsg, &temp.UpdateTime)
//			if err != nil {
//				return nil, errors.Wrap(err, "prepare数据库语句执行结果赋值失败")
//			}
//			return temp, nil
//		}
//		return nil, nil
//	}
//
// // DeleteHistoryContext
// // @Description:  删除全部用户上下文
// // @receiver c
// // @param Db
// // @return error
//
//	func DeleteHistoryContext(userName string) error {
//		sql := `delete from context where name=?`
//		stmt, err := Db.Prepare(sql)
//		if err != nil {
//			return errors.Wrap(err, "prepare数据库查询语句")
//		}
//		defer stmt.Close()
//		_, err = stmt.Exec(userName)
//		if err != nil {
//			return errors.Wrap(err, "prepare数据库语句执行失败")
//		}
//		return nil
//	}

// InsertUser
// @Description: 插入用户
// @receiver u
// @return error

func InsertUser(name string) error {

	newUser := User{
		Name: name,
	}

	// 插入数据到数据库
	result := Db.Create(&newUser)

	// 检查插入是否成功
	if result.Error != nil {
		middleware.MyLogger.Error("Error inserting user:", zap.Any("Error", result.Error))
		return result.Error
	} else {
		middleware.MyLogger.Info("插入用户成功")
		return nil
	}

}

// GetUser
// @Description: 查询用户信息
// @receiver u
// @return error

func GetUser(userName string) (*User, error) {

	var users []User
	result := Db.Where("name = ?", userName).Find(&users)

	if result.Error != nil {
		middleware.MyLogger.Error("Error Getting user:", zap.Any("Error", result.Error))
	} else {
		middleware.MyLogger.Info("查询成功")
	}

	if len(users) == 1 {
		temp := users[0]
		return &temp, nil
	} else if len(users) == 0 {
		middleware.MyLogger.Error("没有此用户", zap.Any("user", userName))
	} else {
		middleware.MyLogger.Error("发现多个用户", zap.Any("users", users))
	}

	return nil, nil
}

// CheckUserAndCreate
// @Description: 查询并创建用户数据
// @param userName
// @return bool
func CheckUserAndCreate(userName string) bool {
	// 查询用户是否存在
	user, err := GetUser(userName)
	if err != nil {
		middleware.MyLogger.Error(xconst.USER_DAO_SEARCH_ERR, zap.Any("用户名", userName))
		return false
	}
	// 存在不创建
	if user != nil {
		return true
	}
	// 创建用户
	middleware.MyLogger.Info(xconst.USER_DAO_FIRST_CREATE, zap.Any("用户名", userName))
	err = InsertUser(userName)
	if err != nil {
		middleware.MyLogger.Error(xconst.USER_DAO_INSERT_ERR, zap.Any("用户名", userName))
		return false
	}
	return true
}
