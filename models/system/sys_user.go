package system

import (
	"gmkdb/global"

	uuid "github.com/satori/go.uuid"
)

type SysUser struct {
	global.GIA_MODEL
	UUID        uuid.UUID      `json:"uuid" gorm:"comment:用户UUID"`                                                           // 用户UUID
	Username    string         `json:"userName" gorm:"comment:用户登录名"`                                                        // 用户登录名
	Password    string         `json:"-"  gorm:"comment:用户登录密码"`                                                             // 用户登录密码
	NickName    string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                            // 用户昵称
	SideMode    string         `json:"sideMode" gorm:"default:dark;comment:用户侧边主题"`                                          // 用户侧边主题
	HeaderImg   string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	BaseColor   string         `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`                                           // 基础颜色
	ActiveColor string         `json:"activeColor" gorm:"default:#1890ff;comment:活跃颜色"`                                      // 活跃颜色
	AuthorityId uint           `json:"authorityId" gorm:"default:888;comment:用户角色ID"`                                        // 用户角色ID
	Authority   SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
	Phone       string         `json:"phone"  gorm:"comment:用户手机号"`                     // 用户手机号
	Email       string         `json:"email"  gorm:"comment:用户邮箱"`                      // 用户邮箱
	Enable      int            `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
}

func (SysUser) TableName() string {
	return "sys_users"
}

// type SysUser struct {
// 	// global.GIA_MODEL
// 	// UUID     uuid.UUID `json:"uuid" gorm:"comment:用户UUID"`                  // 用户UUID

// 	User_id  uint   `gorm:"primarykey"`       // 主键ID
// 	Username string `json:"userName" gorm:""` // 账号
// 	Password string `json:"-"  gorm:""`       // 登录密码

// 	Firstname string `json:"firstname"  gorm:""`
// 	Lastname  string `json:"lastname"  gorm:""`

// 	AuthorityId   uint           `json:"authorityId" gorm:"default:888"`                                 // 用户角色ID
// 	Authority     SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId"` // 用户角色
// 	Authorities   []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
// 	User_group_id uint           // 用户组别

// 	Phone string `json:"phone"  gorm:"comment:用户手机号"` // 用户手机号
// 	Email string `json:"email"  gorm:"comment:用户邮箱"`  // 用户邮箱

// 	NickName    string `json:"nickName" gorm:"default:系统用户"`                                            // 用户昵称
// 	SideMode    string `json:"sideMode" gorm:"default:dark"`                                            // 用户侧边主题
// 	HeaderImg   string `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg"` // 用户头像
// 	BaseColor   string `json:"baseColor" gorm:"default:#fff"`                                           // 基础颜色
// 	ActiveColor string `json:"activeColor" gorm:"default:#1890ff"`                                      // 活跃颜色

// 	Status int `json:"enable" gorm:"default:1"` // Enable 用户是否被冻结 1正常 2冻结

// 	CreatedAt time.Time // 创建时间
// 	UpdatedAt time.Time // 更新时间

// 	// DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
// }
