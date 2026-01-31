package dao

import (
	"bookstore/config"
	"bookstore/model"
	"bookstore/utils"
	"errors"
)

// ErrPasswordIncorrect 自定义错误信息
var ErrPasswordIncorrect = errors.New("密码错误")

// AdminLogin 管理员登录，并生成 session (单点登录)
//func AdminLogin(username, passsword string) (*model.Session, error) {
//	admin, err := GetAdminByUsername(username)
//	if err != nil || admin == nil {
//		return nil, err
//	}
//
//	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(passsword))
//	if err != nil {
//		return nil, ErrPasswordIncorrect
//	}
//
//	DeleteSessionByUserID(admin.ID)
//
//	uuid := utils.GenerateUUID()
//
//	sess := &model.Session{
//		SessionID: uuid,
//		UserName:  admin.Username,
//		UserID:    admin.ID,
//		UserType:  "admin",
//	}
//
//	err = AddSession(sess)
//	if err != nil {
//		return nil, err
//	}
//
//	return sess, nil
//
//}

// CreateAdmin 后台创建管理员
func CreateAdmin(username, password string) error {
	if username == "" || password == "" {
		return errors.New("用户名或密码不能为空")
	}

	//检查用户是否存在
	var exists int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM admins WHERE username = ?", username).
		Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("用户名已存在")
	}

	//utils/hash.go 中方法对密码进行  bcrypt 加密
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	//新增数据加事务
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO admins (username,password) VALUES (?,?)", username, string(hashedPassword))
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetAdminByUsername 根据用户名查询管理员
func GetAdminByUsername(username string) (*model.Admin, error) {
	admin := &model.Admin{}
	err := config.DB.QueryRow("SELECT id,username,password,role FROM admins WHERE username = ?", username).
		Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Role)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// GetAdminByID 根据ID 查询管理员
func GetAdminByID(id int) (*model.Admin, error) {
	admin := &model.Admin{}
	err := config.DB.QueryRow("SELECT id,username,password,role FROM admins WHERE id=?", id).
		Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Role)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// UpdatedAdminPasswordByEmail 根据邮箱更新管理员密码
func UpdatedAdminPasswordByEmail(email, hashedPassword string) error {
	query := "UPDATE admins SET password = ? WHERE email = ?"
	_, err := config.DB.Exec(query, hashedPassword, email)
	return err
}

func GetAdminList(keyword string, page, pageSize int) ([]*model.Admin, int, error) {
	offset := (page - 1) * pageSize

	where := "WHERE 1=1"
	var args []interface{}

	if keyword != "" {
		where += "AND username LIKE ?"
		args = append(args, "%"+keyword+"%")
	}

	countSQL := "SELECT COUNT(*) FROM admins " + where

	var total int
	err := config.DB.QueryRow(countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listSQL := `
       SELECT id,username,password,role
       FROM admins    
       ` + where + `
       ORDER BY id DESC
	   LIMIT ? OFFSET ?	
`

	args = append(args, pageSize, offset)
	rows, err := config.DB.Query(listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	admins := []*model.Admin{}

	for rows.Next() {
		a := &model.Admin{}
		err := rows.Scan(&a.ID, &a.Username, &a.Password, &a.Role)
		if err != nil {
			return nil, 0, err
		}
		admins = append(admins, a)
	}

	return admins, total, nil
}
