package daogorm

import (
	"errors"
	"gmkdb/utils"

	"gorm.io/gorm"
)

type gormDB struct {
	db *gorm.DB
}

// # INSERT DELETE SELECT UPDATE
func NewDao(db *gorm.DB) MysqlDao {
	return &gormDB{
		db: db,
	}
}

func checkValue(datas any) (err error) {
	if ok := utils.IsStruct(datas); !ok {
		err = errors.New("sql data err: data must be a pointer to a struct")
	}
	return
}

func checkValues(datas any) (err error) {
	if ok := utils.IsStructs(datas); !ok {
		err = errors.New("sql data err: data must be a pointer to a dbtables")
	}
	return
}
