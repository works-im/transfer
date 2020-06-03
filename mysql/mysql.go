package mysql

import (
	"fmt"
	"net/url"
	"strings"

	"database/sql"

	"github.com/araddon/dateparse"
	_ "github.com/go-sql-driver/mysql" // mysql driver

	"github.com/works-im/transfer/database"
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

// MappingFields return map[source]Field
func (mysql *MySQL) MappingFields() []Field {

	dbFields := mysql.schema.FieldMap()

	var fields = []Field{}

	for _, mappingField := range mysql.Mapping {

		if meta, exists := dbFields[mappingField.Target]; exists {

			if ts := strings.Split(strings.ToLower(meta.Type), "("); len(ts) > 0 {
				field := Field{
					Name:      meta.Field,
					Type:      DateType(ts[0]),
					Converter: mappingField.Converter,
				}
				fields = append(fields, field)
			}
		}
	}

	return fields
}

// Reader database
func (mysql *MySQL) Reader(query database.Query) (packet database.Packet, err error) {

	return nil, nil
}

// Writer data
func (mysql *MySQL) Writer(packet database.Packet) error {

	var (
		mappingFields = mysql.MappingFields()
	)

	for _, row := range packet {

		var (
			selectFields []string
			values       []string
			updateFields []string
		)

		for _, field := range mappingFields {

			selectFields = append(selectFields, field.Name)
			updateFields = append(updateFields, fmt.Sprintf("%s=VALUES(%s)", field.Name, field.Name))

			value := row[field.Name]
			if nil == value {
				values = append(values, "NULL")
				continue
			}

			switch field.Type {
			case Tinyint, SmallInt, MediumInt, Int, Integer, BigInt, Float, Double, Decimal:
				values = append(values, fmt.Sprint(value))

			case Char, Varchar, TinyText, Text, MediumText, LongText:
				values = append(values, "'"+fmt.Sprint(value)+"'")

			case TinyBlob, Blob, MediumBlob, LongBlob:
				values = append(values, "'"+fmt.Sprint(value)+"'")

			case Date, Time, Year, Datetime, Timestamp:
				time, err := dateparse.ParseAny(fmt.Sprint(value))
				if err != nil {
					return err
				}
				switch field.Type {
				case Date, Time, Year:
					values = append(values, field.Type.String()+"(FROM_UNIXTIME("+fmt.Sprint(time.Unix())+"))")
				case Datetime, Timestamp:
					values = append(values, "FROM_UNIXTIME("+fmt.Sprint(time.Unix())+")")
				}
			}
		}

		// INSERT INTO table_name (field1, field2) VALUES (value1, value2) ON DUPLICATE KEY UPDATE field1=VALUES(field1), field2=VALUES(field2)
		sql := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
			mysql.TableName,
			strings.Join(selectFields, ", "),
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
