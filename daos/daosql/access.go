package daosql

import (
	"database/sql"
	"strings"
)

type SqlDao struct { // SqlDB SqlDao
	taos   *sql.DB
	dbName string
}

// # INSERT DELETE SELECT UPDATE
func NewDao(taos *sql.DB, dbName string) *SqlDao {
	return &SqlDao{
		taos:   taos,
		dbName: dbName,
	}
}

func (s *SqlDao) dbSTable(stableName string, _dbName ...string) string {
	dbName := s.dbName
	if len(_dbName) > 0 && strings.TrimSpace(_dbName[0]) != "" {
		dbName = _dbName[0]
	}
	return dbName + "." + stableName
}

func (t *SqlDao) Close() (err error) {
	if t.taos != nil {
		err = t.taos.Ping()
		if err == nil {
			err = t.taos.Close()
		}
	}
	return
}
