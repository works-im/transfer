package mysql

// DateType MySQL data type
type DateType string

func (t DateType) String() string {
	return t.String()
}

const (
	Tinyint   DateType = "tinyint"
	SmallInt  DateType = "smallint"
	MediumInt DateType = "mediumint"
	Int       DateType = "int"
	Integer   DateType = "integer"
	BigInt    DateType = "bigint"

	Float   DateType = "float"
	Double  DateType = "double"
	Decimal DateType = "decimal"

	Date      DateType = "date"
	Time      DateType = "time"
	Year      DateType = "year"
	Datetime  DateType = "datetime"
	Timestamp DateType = "timestamp"

	Char       DateType = "char"
	Varchar    DateType = "varchar"
	TinyText   DateType = "tinytext"
	Text       DateType = "text"
	MediumText DateType = "mediumtext"
	LongText   DateType = "longtext"

	TinyBlob   DateType = "tinyblob"
	Blob       DateType = "blob"
	MediumBlob DateType = "mediumblob"
	LongBlob   DateType = "longblob"
)

// Field for database table field schema
type Field struct {
	Name      string
	Type      DateType
	Converter string
}
