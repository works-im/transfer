package mysql

import (
	"fmt"
	"net/url"
	"strings"

	"database/sql"

	"github.com/araddon/dateparse"
	_ "github.com/go-sql-driver/mysql" // mysql driver

	"transfer/database"
)

// MySQL database transfer
type MySQL struct {
	db     *sql.DB
	schema database.Schema

	Driver    database.Driver
	TableName string
	Mapping   database.Mapping
}

// NewMySQL return MySQL transfer
func NewMySQL(args *database.Options) (*MySQL, error) {

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

	if err := mysql.TableSchema(); err != nil {
		return nil, err
	}

	return mysql, nil
}

// Reader database
func (mysql *MySQL) Reader(query database.Query) (packet database.Packet, err error) {

	return nil, nil
}

// Writer data
func (mysql *MySQL) Writer(packet database.Packet) error {

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
			value := row[field.Target]

			// 忽略 nil
			if nil == value {
				continue
			}

			switch value.(type) {
			case int64, float64, bool:
				values = append(values, fmt.Sprint(value))
				continue
			}

			switch field.TargetType {
			case "timestamp":
				time, err := dateparse.ParseAny(fmt.Sprint(value))
				if err != nil {
					return err
				}
				values = append(values, "FROM_UNIXTIME("+fmt.Sprint(time.Unix())+")")
			default:
				values = append(values, "'"+fmt.Sprint(value)+"'")
			}
		}

		sql := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
			mysql.TableName,
			strings.Join(fields, ", "),
			strings.Join(values, ", "),
			strings.Join(updateFields, ", "),
		)

		fmt.Printf("sql: %s\n", sql)

		result, err := mysql.db.Exec(sql)
		if err != nil {
			return err
		}

		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		fmt.Println("affected rows: ", affectedRows)
	}
	return nil
}

// TableSchema return table schema
func (mysql *MySQL) TableSchema() error {
	rows, err := mysql.db.Query("DESCRIBE " + mysql.TableName)
	if err != nil {
		return err
	}

	defer rows.Close()

	var schema []database.FieldMeta

	for rows.Next() {
		var Field, Type string
		var Null, Key, Default, Extra interface{}
		if err := rows.Scan(&Field, &Type, &Null, &Key, &Default, &Extra); err != nil {
			return err
		}

		schema = append(schema, database.FieldMeta{
			Field:   Field,
			Type:    Type,
			Null:    Null,
			Key:     Key,
			Default: Default,
			Extra:   Extra,
		})

		mysql.schema = schema
	}

	return rows.Err()
}

// Fields func
func (mysql *MySQL) Fields() map[string]Field {

	var fields []string
	dbFields := mysql.schema.FieldMap()

	var fs = map[string]Field{}

	for _, field := range mysql.Mapping {

		if v, e := dbFields[field.Target]; e {

			fields = append(fields, v.Field)

			if ts := strings.Split(v.Type, "("); len(ts) > 0 {
				fs[v.Field] = Field{
					Name:      v.Field,
					Type:      DateType(ts[0]),
					Converter: field.Converter,
				}
			}

		}
	}

	return fs
}
