package transfer

import (
	"fmt"
	"net/url"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// MySQL database transfer
type MySQL struct {
	db *sql.DB

	Driver    Driver
	TableName string
	Mapping   Mapping
}

// NewMySQL return MySQL transfer
func NewMySQL(args *DatabaseOptions) (*MySQL, error) {

	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=%s&parseTime=true",
		args.Driver.Username,
		args.Driver.Password,
		args.Driver.Host,
		args.Driver.Port,
		args.Driver.Database,
		url.QueryEscape("Asia/Shanghai"),
	)

	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	mysql := &MySQL{
		db:        db,
		TableName: args.Driver.Table,
		Mapping:   args.Mapping,
	}

	return mysql, nil
}

// Reader database
func (mysql *MySQL) Reader(query Query) (packet Packet, err error) {

	return nil, nil
}

// Writer data
func (mysql *MySQL) Writer(packet Packet) error {

	var (
		fields       = mysql.Mapping.Fields()
		updateFields []string
	)

	for _, field := range mysql.Mapping {
		updateFields = append(updateFields, fmt.Sprintf("%s=VALUES(%s)", field.Target, field.Target))
	}

	for _, row := range packet {

		var values []string
		for _, field := range mysql.Mapping {
			values = append(values, "'"+fmt.Sprint(row[field.Target])+"'")
		}

		sql := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
			mysql.TableName,
			strings.Join(fields, ", "),
			strings.Join(values, ", "),
			strings.Join(updateFields, ", "),
		)

		fmt.Printf("sql: %s\n", sql)

		if _, err := mysql.db.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}
