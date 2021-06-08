package generator

import (
	"database/sql"
	"fmt"
	"html/template"
	"os"
	"prof-filer/db"
	"prof-filer/util"
	"strings"
)

type DaoSpec struct {
	Table string
}

type AuthorSpec struct {
	Author string
}

type ColumnSchema struct {
	COLUMN_NAME    string
	COLUMN_TYPE    string
	COLUMN_TYPE_TS string
	IS_NULLABLE    string
	COLUMN_DEFAULT string
}

type DaoConfig struct {
	ClassName       string
	Author          string
	Date            string
	InterfaceName   string
	InterfaceSchema string
	ListQuery       string
	InsertQuery     string
}

// Dictionary to convert MYSQL data type to Typescript equivalant
var MYSQL_TO_TS_DICT = map[string]string{
	"tinyint": "0 | 1",
	"int":     "number",
	"bigint":  "number",
	"float":   "number",
	"varchar": "string",
	"char":    "string",
	"date":    "string",
	"text":    "string",
	"json":    "string",
	"enum":    "string",
}

// Generates a DAO file based on the table name
func GenerateDao(spec *DaoSpec, authorSpec *AuthorSpec) {
	conn, err := db.Connection()

	if err != nil {
		panic(err)
	}

	fmt.Println("===========================================================================================")
	fmt.Println("Generating DAO specifications for " + spec.Table + " for " + authorSpec.Author)
	fmt.Println("===========================================================================================")

	columns, err := listColumns(conn, spec)

	if err != nil {
		panic(err)
	}

	columnNames := make([]string, 0)

	for _, column := range columns {
		columnNames = append(columnNames, column.COLUMN_NAME)
	}

	listQuery := getListQuery(spec, authorSpec, columnNames)
	insertQuery := getInsertMethod(spec, authorSpec, columnNames)
	daoInterface := getDaoInterface(spec, authorSpec, columns)

	config := DaoConfig{
		Author:          authorSpec.Author,
		Date:            util.GetCurrDate(),
		ClassName:       spec.Table,
		InterfaceName:   spec.Table + "DAO",
		InterfaceSchema: daoInterface,
		ListQuery:       listQuery,
		InsertQuery:     insertQuery,
	}

	WriteDaoUsingTemplate(&config)

	fmt.Println("Your output has been saved")
	fmt.Println("Happy Hacking!!")

	defer conn.Close()
}

// List table columns
func listColumns(conn *sql.DB, spec *DaoSpec) ([]ColumnSchema, error) {
	columns := make([]ColumnSchema, 0)

	const sqlStatement = "SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_DEFAULT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"
	data, err := conn.Query(sqlStatement, os.Getenv("MYSQL_DATABASE_NAME"), spec.Table)

	if err != nil {
		return columns, err
	}

	for data.Next() {
		var cm ColumnSchema
		data.Scan(&cm.COLUMN_NAME, &cm.COLUMN_TYPE, &cm.IS_NULLABLE, &cm.COLUMN_DEFAULT)
		cm.COLUMN_NAME = util.ToCamelCase(cm.COLUMN_NAME)
		cm.COLUMN_TYPE_TS = MysqlTypesToTypescript(cm.COLUMN_TYPE)
		columns = append(columns, cm)
	}

	err = data.Err()

	return columns, err
}

// insert method generator
func getInsertMethod(spec *DaoSpec, authorSpec *AuthorSpec, columns []string) string {
	insertStmt := "INSERT INTO " + spec.Table + "(" + strings.Join(columns, ", ") + ") VALUES ?"
	return insertStmt
}

// returns list query
func getListQuery(spec *DaoSpec, authorSpec *AuthorSpec, columns []string) string {
	tableAlias := strings.ToLower(spec.Table[:1])

	aliasedColumns := make([]string, 0)
	for _, column := range columns {
		aliasedColumns = append(aliasedColumns, tableAlias+"."+column)
	}

	listStmt := "SELECT " + strings.Join(aliasedColumns, ", ") + " FROM " + spec.Table + " " + tableAlias
	return listStmt
}

// returns Typescript compatible interfaces
func getDaoInterface(spec *DaoSpec, authorSpec *AuthorSpec, columns []ColumnSchema) string {
	str := "{ "

	for _, elem := range columns {
		str += "" + elem.COLUMN_NAME + ": " + elem.COLUMN_TYPE_TS + "; "
	}

	str += " }"

	return str
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// MySQL
func MysqlTypesToTypescript(dataType string) string {
	dt := strings.Split(dataType, "(")

	if dt[0] == "enum" {
		dt[1] = strings.ReplaceAll(dt[1], ")", "")
		byteStr := strings.Split(dt[1], ",")
		dt[1] = strings.Join(byteStr, " | ")
		dt[0] = dt[1]
	} else {
		dt[0] = MYSQL_TO_TS_DICT[dt[0]]
	}

	return dt[0]
}

// Uses a template to write to a file
func WriteDaoUsingTemplate(config *DaoConfig) {
	cwd, err := os.Getwd()
	check(err)
	templatePath := cwd + "/templates/dao.template.txt"
	t, err := template.ParseFiles(templatePath)
	check(err)

	outputFilePath := cwd + "/output/dao.ts"

	f, err := os.Create(outputFilePath)
	check(err)
	defer f.Close()

	err = t.Execute(f, config)

	check(err)
}
