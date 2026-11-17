package daosql

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// # INSERT DELETE SELECT UPDATE

type sqlKeyValue struct {
	key   string
	value string
}

// 封装结构体的字段集 (*struct{})
func fieldColumns(data TABBASE) (columns []fieldInfo) {
	val := reflect.ValueOf(data)
	// fmt.Printf("fieldColumns .Kind(): %v\n", val.Kind())
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	parseField(val, &columns)
	return
}
func parseField(val reflect.Value, columns *[]fieldInfo) { // 解析各字段信息, 传入需指针结构体, 否则Value去不到值会抛出
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Kind() == reflect.Struct && val.Field(i).Type().String() != "time.Time" {
			parseField(val.Field(i), columns)
		} else {
			// GetFields(val.Field(i))
			field := val.Type().Field(i).Tag.Get("db")       // 数据库字段名
			fieldType := val.Type().Field(i).Tag.Get("taos") // 数据库类型
			kind := val.Field(i).Type().Kind()               // 原结构体类型 // val.Field(i).Type() = val.Type().Field(i).Type
			// if kind == reflect.Struct {
			// 	kind = val.Field(i).Type().String()
			// }
			switch field {
			case "-":
			default:
				// fmt.Printf("val.Kind(): %v\n", val.Type())
				// fmt.Printf("Inte(): %v\n", val.Field(i).Interface())
				// fmt.Printf("CanAddr(): %v\n", val.Field(i).CanAddr())

				newFv := fieldInfo{Field: field, Type: fieldType, Kind: kind}
				// fmt.Printf("val.Field(%d): %v %v %v\n", i, val.Field(i).Type(), val.Field(i).CanAddr(), val.Field(i))
				if val.Field(i).CanAddr() { // 是否可以寻址
					newFv.Value = val.Field(i).Addr().Interface()
				} else {
					newFv.Value = val.Field(i).Interface()
				}
				// fmt.Printf("val.Value(%d): %v \n", i, newFv.Value)
				if len(fieldType) > 4 && strings.ToUpper(fieldType[:4]) == "TAG;" {
					newFv.Type, newFv.Note = fieldType[4:], "TAG"
				} else if len(fieldType) > 5 && strings.ToUpper(fieldType[:5]) == "NULL;" {
					newFv.Type, newFv.Note = fieldType[5:], "NULL"
				}
				// // fieldInfo: {ts TIMESTAMP 0     struct 2025-05-19 05:12:30 +0800 CST}
				// if newFv.Type == "TIMESTAMP" && newFv.Kind == reflect.Struct {
				// 	// newFv.Value =
				// }
				// fmt.Printf("fieldInfo: %v\n", newFv)
				*columns = append(*columns, newFv)
			}
		}
	}
}

// *

// -

// =

// =

// 已经实现功能, 待清理函数
func (s *SqlDao) Select2(dest TABBASE, args ...string) (err error) {

	fmt.Printf("s.dbName: %v\n", s.dbName)
	query := "SELECT"

	query += fmt.Sprintf(` %s`, "ts,end_time,open_price,close_price,volume,trade_num,is_final,s,i")
	query += fmt.Sprintf(` FROM %s`, s.dbSTable(dest.SubTableName()))

	// query += fmt.Sprintf(` WHERE start_time = %d`, 123)
	query += " ORDER BY ts DESC"
	query += fmt.Sprintf(` LIMIT %d`, 3)

	// 查询区间, ts -> startTime -> endTime
	// 查询数量 limit = ? 最大3000?
	// 倒序 ORDER BY ts DESC

	fmt.Printf("query: %v\n", query)
	rows, err := s.taos.Query(query)
	if err != nil {
		fmt.Printf("Query err: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ts time.Time
		var startTime, endTime int64
		var open, close, volume string
		var num int64
		var isFinal bool
		var symbol, interval string

		var as []any
		as = append(as, &ts)
		as = append(as, &endTime)
		as = append(as, &open)
		as = append(as, &close)
		as = append(as, &volume)
		as = append(as, &num)
		as = append(as, &isFinal)
		as = append(as, &symbol)
		as = append(as, &interval)
		// as := make([]any, 7)
		err = rows.Scan(as...)

		// err = rows.Scan(&ts, &endTime, &open, &close, &volume, &num, &isFinal)
		if err != nil {
			fmt.Printf("Scan err: %v\n", err)
		}
		startTime = ts.UnixMilli()
		// fmt.Printf("dest: %v\n", dest)
		fmt.Printf("ts: %v - ", ts)
		fmt.Println(startTime, endTime, open, close, volume, num, isFinal, symbol, interval)
		fmt.Println("any: ", as)

	}

	return nil
}

func (s *SqlDao) Delete(datas any, args ...string) error {
	return nil
}

// Pass insert 包含更新操作
func (s *SqlDao) Update(data any, args ...map[string]any) error {
	return nil
}
