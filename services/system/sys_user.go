package system

import (
	"errors"
	"fmt"
	"gmkdb/global"
	"gmkdb/models/system"
	"gmkdb/utils"

	// uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: userInter system.SysUser, err error

type UserService struct{}

func (userService *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.GIA_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	// 否则 附加uuid 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	// u.UUID = uuid.NewV4()
	err = global.GIA_DB.Create(&u).Error
	return u, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) Login(u *system.SysUser) (userInter *system.SysUser, err error) {
	if nil == global.GIA_DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.SysUser
	err = global.GIA_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}

		var SysAuthorityMenus []system.SysAuthorityMenu
		err = global.GIA_DB.Where("sys_authority_authority_id = ?", user.AuthorityId).Find(&SysAuthorityMenus).Error
		if err != nil {
			return
		}

		var MenuIds []string

		for i := range SysAuthorityMenus {
			MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
		}

		var am system.SysBaseMenu
		ferr := global.GIA_DB.First(&am, "name = ? and id in (?)", user.Authority.DefaultRouter, MenuIds).Error
		if errors.Is(ferr, gorm.ErrRecordNotFound) {
			user.Authority.DefaultRouter = "404"
		}
	}

	return &user, err
}
