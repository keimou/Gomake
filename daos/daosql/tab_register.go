package daosql

import (
	"fmt"
	"gmkdb/utils"
	"log"
	"reflect"
	"strings"
)

type TABBASE interface {
	TableName() string
	SubTableName(s ...string) string
}

type fieldInfo struct {
	Field    string
	Type     string
	Length   int    // 某些数据类型需要长度（如BINARY）
	Note     string // 额外说明（TAG/普通列）
	Encode   string
	Compress string
	Level    string

	Kind  reflect.Kind
	Value any
}

// CREATE STABLE IF NOT EXISTS kline_[s] (%v) TAGS (s BINARY(50),i BINARY(5));
// 注册超级表
func (s *SqlDao) AutoMigrate(dst ...TABBASE) (err error) {

	for _, dest := range dst {
		var columns []fieldInfo
		val := reflect.TypeOf(dest) // ValueOf
		dbTableField(val, &columns)
		// fmt.Printf("columns: %v\n", columns)
		// columns := fieldColumns(dest) // 不通用, 后期考虑兼容调试

		// fmt.Printf("columns: %v\n", columns)
		// fmt.Println("DESCRIBE stables sql exec..." + s.dbSTable(dest.TableName()))
		// 获取超级表结构
		if rows, err := s.taos.Query(fmt.Sprintf("DESCRIBE %s", s.dbSTable(dest.TableName()))); err != nil { // [0x2603]Table does not exist
			// fmt.Println("no existing, creatSupTable: ", err)
			var tableFields, tags string
			for _, v := range columns {
				if v.Note == "TAG" {
					utils.StrJoin(&tags, v.Field+" "+v.Type)
				} else {
					utils.StrJoin(&tableFields, v.Field+" "+v.Type)
				}
			}
			superTableSQL := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s (%v) TAGS (%s)", s.dbSTable(dest.TableName()), tableFields, tags)
			// fmt.Printf("superTableSQL: %v\n", superTableSQL)
			if _, err := s.taos.Exec(superTableSQL); err != nil {
				fmt.Printf("superTableSQL err: %v\n'%v'\n", err, superTableSQL)
				return err
				// } else {
				// 	if rowsAffected, err := res.RowsAffected(); err != nil {
				// 		return err
				// 	} else {
				// 		fmt.Println("Create stable successfully, rowsAffected:", rowsAffected)
				// 	}
			}
			// fmt.Printf("superTableSQL end: %v\n", err)

		} else {
			// fmt.Println("existing, add column")
			defer rows.Close()
			existingColumns := make(map[string]bool)
			for rows.Next() {
				var f fieldInfo
				if err := rows.Scan(&f.Field, &f.Type, &f.Length, &f.Note, &f.Encode, &f.Compress, &f.Level); err != nil {
					return err
				}
				existingColumns[f.Field] = true
			}
			for _, col := range columns {
				if !existingColumns[col.Field] && col.Note != "TAG" {
					sql := fmt.Sprintf("ALTER STABLE %s ADD COLUMN %s %s", s.dbSTable(dest.TableName()), col.Field, col.Type)
					_, err := s.taos.Exec(sql)
					if err != nil {
						log.Printf("添加列 %s 失败: %v", col.Field, err)
					} else {
						log.Printf("已添加列: %s; ", col.Field)
					}
				}
			}
		}
	}
	return
}

func dbTableField(val reflect.Type, columns *[]fieldInfo) {
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Type.Kind() == reflect.Struct && val.Field(i).Type.String() != "time.Time" {
			// fmt.Printf(".Type: %v\n", val.Field(i).Type)
			// fmt.Printf(".Kind(): %v\n", val.Field(i).Type.Kind())
			dbTableField(val.Field(i).Type, columns)
		} else {
			// fmt.Printf("t: %v \t %v\n", val.Field(i).Tag.Get("db"), val.Field(i).Type)
			// GetFields(val.Field(i))
			field := val.Field(i).Tag.Get("db")
			kind := val.Field(i).Tag.Get("taos")
			switch field {
			case "-":
			default:
				if len(kind) > 4 && strings.ToUpper(kind[:4]) == "TAG;" {
					kind = kind[4:]
					*columns = append(*columns, fieldInfo{Field: field, Type: kind, Note: "TAG"})
					// *tags = utils.StrJoin(*tags, field+" "+kind)
				} else if len(kind) > 5 && strings.ToUpper(kind[:5]) == "NULL;" {
					kind = kind[5:]
					*columns = append(*columns, fieldInfo{Field: field, Type: kind, Note: "NULL"})
					// newFv.Type, newFv.Note = fieldType[2:], "-"
				} else {
					// fmt.Printf("field: %v %v\n", field, kind)
					*columns = append(*columns, fieldInfo{Field: field, Type: kind})
					// *tableFields = utils.StrJoin(*tableFields, field+" "+kind)
				}
			}
		}
	}
}

// 注册子表
