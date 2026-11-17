package ninex

import (
	"gmkdb/global"
	"time"
)

type NexRecord struct {
	// global.GIA_MODEL
	PredictionID string `json:"predictionId" gorm:"primarykey;type:varchar(20);comment:预测ID"` // 20262982313211553 => 2026 298 2 31 3211553 （date period lotteryType  predictionType Uid）

	Values string `json:"value" gorm:"not null;type:varchar(30);comment:预测生肖组"` // 预测的生肖数组

	// Period           int `json:"period" gorm:"comment:期数"`                          // 298
	Period int `json:"period" gorm:"index;not null;comment:期数"` // 2026298

	Success         bool      `json:"success" gorm:"comment:是否命中"`
	ContinuousCount string    `json:"continuousCount" gorm:"type:varchar(12);comment:连中数"`
	LotteryNumber   string    `json:"lotteryNumber" gorm:"type:varchar(12);comment:开奖"`
	CreateTime      time.Time `json:"createTime" gorm:"comment:创建时间"` // "2026年01月12日 15:50:58" 2026014

	LotteryType      int `json:"lotteryType" gorm:"default:2;comment:彩票类型"`         // 2
	PredictionTypeID int `json:"predictionTypeId" gorm:"default:31;comment:预测类型ID"` // 31

	UserId uint `json:"userId" gorm:"index"`
	// Expert NexExpert `json:"expert" gorm:"foreignKey:UserId;references:UserId;constraint:false;comment:关联角色"`

	global.TS_MODEL
}

func (NexRecord) TableName() string {
	return "nex_record"
}
