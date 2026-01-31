package dao

import (
	"bookstore/config"
	"bookstore/model"
	"bookstore/utils"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// CheckUserNameAndPassword 根据用户名和密码从数据库中查询一条记录
// 参数：
//
//	username - 用户输入的用户名
//	password - 用户输入的密码
//
// 返回值：
//
//	*model.User - 查询到的用户对象指针
//	error      - 错误信息（如果有错误）
func CheckUserNameAndPassword(username string, password string) (*model.User, error) {
	// 定义 SQL 语句，占位符 ? 会被后面的参数替换
	sqlStr := "select id, username,password,email from users where username = ?"

	// 执行查询（返回一行记录）
	// QueryRow 会把结果放在 row 中，如果没有匹配的记录，row.Scan 会返回错误
	row := config.DB.QueryRow(sqlStr, username)
	user := &model.User{}

	// 将查询结果扫描（赋值）到 user 的各个字段
	// 注意：如果没有记录，Scan 会返回 sql.ErrNoRows 错误
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			// 用户不存在
			return &model.User{}, nil
		}
		// 其他数据库错误
		return nil, err
	}

	// 使用 bcrypt 比对哈希密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return &model.User{}, nil
	}
	// 返回 user 对象和 nil 错误（这里没做错误处理）
	return user, nil
}

// CheckUserName 根据用户名和密码从数据库中查询一条记录
func CheckUserName(username string) (*model.User, error) {
	sqlStr := "select id,username,password,email from users where username = ?"
	row := config.DB.QueryRow(sqlStr, username)
	user := &model.User{}
	row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	return user, nil
}

// CheckEmail 检查邮箱是否已存在
func CheckEmail(email string) (*model.User, error) {
	var user model.User
	err := config.DB.QueryRow("SELECT id,username,email FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil //没有找到
		}
		return nil, err
	}
	return &user, nil
}

// SaveUser 向数据库中插入用户信息
func SaveUser(username string, password string, email string) error {
	// 对密码进行 bcrypt 哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败:%w", err)
	}

	//开启事务
	tx, err := config.DB.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败:%w", err)
	}

	//确保在函数返回前正确回滚或提交
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// 执行插入操作
	sqlStr := "insert into users(username,password,email) values(?,?,?)"
	_, err = tx.Exec(sqlStr, username, string(hashedPassword), email)
	if err != nil {
		return fmt.Errorf("插入用户失败:%w", err)
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败:%w", err)
	}
	return nil
}

// GetUserByID 根据用户ID 获取用户信息
func GetUserByID(id int) (*model.User, error) {
	user := &model.User{}
	err := config.DB.QueryRow("SELECT id,username,email FROM users WHERE id= ? ", id).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// UpdateProfileByID 根据用户ID 更新个人资料
func UpdateProfileByID(userID int, username, email string) (*model.User, error) {
	fields := []string{}
	args := []interface{}{}

	if username != "" {
		fields = append(fields, "username = ?")
		args = append(args, username)
	}
	if email != "" {
		fields = append(fields, "email = ?")
		args = append(args, email)
	}

	if len(fields) == 0 {
		return nil, errors.New("用户名或邮箱不能为空")
	}

	// 开启事务
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}

	// 确保在出错时回滚事务
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	// 执行更新
	query := "UPDATE users SET " + strings.Join(fields, ", ") + "WHERE id = ?"
	args = append(args, userID)

	result, err := tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("用户不存在或资料未修改")
	}

	//在同一事务中读取最新数据
	user := &model.User{}
	err = tx.QueryRow("SELECT id, username,email FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	//提交事务
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	committed = true
	return user, nil
}

// UpdatePasswordByID 根据用户ID修改密码
func UpdatePasswordByID(userID int, oldPassword, newPassword string) error {
	if oldPassword == "" || newPassword == "" {
		return errors.New("旧密码和新密码都必须提供")
	}

	// 开启事务
	tx, err := config.DB.Begin()
	if err != nil {
		return errors.New("开启数据库事务失败")
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	// 查询用户当前密码
	var hashedPassword string
	err = tx.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("用户不存在")
		}
		return errors.New("查询当前密码失败")
	}

	// 验证旧密码是否正确
	if !utils.CheckPassword(oldPassword, hashedPassword) {
		return errors.New("旧密码错误")
	}

	// 对新密码进行哈希
	newHashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	result, err := tx.Exec("UPDATE users SET password = ? WHERE id = ?", newHashed, userID)
	if err != nil {
		return errors.New("更新密码失败")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("获取受影响的行数失败")
	}

	if rowsAffected == 0 {
		return errors.New("密码未修改")
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return errors.New("提交事务失败")
	}
	committed = true
	return nil
}

// UpdateUserPasswordByEmail 根据邮箱更新普通用户密码
func UpdateUserPasswordByEmail(email, hashedPassword string) error {
	query := "UPDATE users SET password = ? WHERE email = ?"
	_, err := config.DB.Exec(query, hashedPassword, email)
	return err
}
