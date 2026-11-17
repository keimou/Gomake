package ninex

import (
	"gmkdb/global"
)

type NexExpert struct {
	UserId   uint   `json:"userId" gorm:"primarykey"`
	NickName string `json:"nickname" gorm:"type:varchar(50);comment:用户昵称"`

	FromForumsName string `json:"fromForumsName" gorm:"type:varchar(50);comment:来自论坛名称"`
	HeadPic        string `json:"headPic" gorm:"comment:头像"`

	FansCount    int `json:"fansCount" gorm:"comment:粉丝数"`
	ShareCount   int `json:"shareCount" gorm:"comment:分享数"`
	ThumbUpCount int `json:"thumbUpCount" gorm:"comment:点赞数"`

	// Records []NexRecord `json:"-" gorm:"foreignKey:UserId;references:UserId;constraint:false"`

	Enable int `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 0冻结"`
	global.TS_MODEL
}

func (NexExpert) TableName() string {
	return "nex_expert"
}
