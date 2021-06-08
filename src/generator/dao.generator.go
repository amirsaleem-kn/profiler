package generator

import (
	"database/sql"
	"fmt"
	"os"
	"prof-filer/db"
	"prof-filer/util"
	"strings"
	"text/template"
)

type DaoSpec struct {
	Table string
}

type AuthorSpec struct {
	Author string
}

type ColumnSchema struct {
	COLUMN_NAME       string
	COLUMN_NAME_ALIAS string
	COLUMN_TYPE       string
	COLUMN_TYPE_TS    string
	IS_NULLABLE       string
	COLUMN_DEFAULT    string
}

type DaoConfig struct {
	ClassName              string
	Author                 string
	Date                   string
	InterfaceName          string
	InterfaceSchema        string
	MetaInterfaceName      string
	MetaInterfaceSchema    string
	FiltersInterfaceName   string
	FiltersInterfaceSchema string
	FILTERS_MODEL          string
	ListQuery              string
	InsertQuery            string
	InsertValues           string
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

	util.HandleException(err)

	fmt.Println("===========================================================================================")
	fmt.Println("Generating DAO specifications for " + spec.Table + " for " + authorSpec.Author)
	fmt.Println("===========================================================================================")

	columns, err := listColumns(conn, spec)

	util.HandleException(err)

	columnNames := make([]string, 0)
	aliasedColumnNames := make([]string, 0)

	for _, column := range columns {
		columnNames = append(columnNames, column.COLUMN_NAME)
		aliasedColumnNames = append(aliasedColumnNames, column.COLUMN_NAME_ALIAS)
	}

	// Generates a SELECT query
	listQuery := getListQuery(spec, authorSpec, aliasedColumnNames)

	// Generates an INSERT query
	insertQuery := getInsertQuery(spec, authorSpec, columnNames)
	// Generates the payload map method
	insertValues := getInsertValues(spec, authorSpec, aliasedColumnNames)

	// generates Typescript compatible interface for the DAO insert
	daoInterface := getDaoInterface(spec, authorSpec, columns)
	// generates Typescript compatible interface for Filters
	filtersInterface := getFiltersInterface(spec, authorSpec, columns)

	// filters model
	filtersModel := getFiltersModel(spec, authorSpec, columns)

	// Template Configuration
	config := DaoConfig{
		Author:                 authorSpec.Author,
		Date:                   util.GetCurrDate(),
		ClassName:              spec.Table,
		InterfaceName:          spec.Table + "DAO",
		InterfaceSchema:        daoInterface,
		MetaInterfaceName:      spec.Table + "DAOWithMeta",
		FiltersInterfaceName:   spec.Table + "ListFilters",
		FiltersInterfaceSchema: filtersInterface,
		FILTERS_MODEL:          filtersModel,
		ListQuery:              listQuery,
		InsertQuery:            insertQuery,
		InsertValues:           insertValues,
	}

	// Write DAO in output dir using a template file
	WriteDaoUsingTemplate(&config)

	fmt.Println("Your output has been saved")
	fmt.Println("Happy Hacking!!")

	// close the database connection
	defer conn.Close()
}

// List table columns
func listColumns(conn *sql.DB, spec *DaoSpec) ([]ColumnSchema, error) {
	columns := make([]ColumnSchema, 0)
	tableAlias := strings.ToLower(spec.Table[:1])

	const sqlStatement = "SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_DEFAULT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"
	data, err := conn.Query(sqlStatement, os.Getenv("MYSQL_DATABASE_NAME"), spec.Table)

	util.HandleException(err)

	for data.Next() {
		var cm ColumnSchema
		data.Scan(&cm.COLUMN_NAME, &cm.COLUMN_TYPE, &cm.IS_NULLABLE, &cm.COLUMN_DEFAULT)
		cm.COLUMN_NAME = util.ToCamelCase(cm.COLUMN_NAME)
		cm.COLUMN_NAME_ALIAS = tableAlias + "." + cm.COLUMN_NAME
		cm.COLUMN_TYPE_TS = MysqlTypesToTypescript(cm.COLUMN_TYPE)
		columns = append(columns, cm)
	}

	err = data.Err()

	return columns, err
}

// insert method generator
func getInsertQuery(spec *DaoSpec, authorSpec *AuthorSpec, columns []string) string {
	insertStmt := "INSERT INTO " + spec.Table + "(" + strings.Join(columns, ", ") + ") VALUES ?"
	return insertStmt
}

// returns list query
func getListQuery(spec *DaoSpec, authorSpec *AuthorSpec, aliasedColumnNames []string) string {
	tableAlias := strings.ToLower(spec.Table[:1])
	listStmt := "SELECT " + strings.Join(aliasedColumnNames, ", ") + " FROM " + spec.Table + " " + tableAlias
	return listStmt
}

// get Typescript compatible insert values
func getInsertValues(spec *DaoSpec, authorSpec *AuthorSpec, aliasedColumnNames []string) string {
	str := "[payload.map((p) => [" + strings.Join(aliasedColumnNames, ", ") + "])]"
	return str
}

// returns Typescript compatible interfaces
func getDaoInterface(spec *DaoSpec, authorSpec *AuthorSpec, columns []ColumnSchema) string {
	str := "{ "

	for _, elem := range columns {
		str += "" + elem.COLUMN_NAME + ": " + elem.COLUMN_TYPE_TS + ";\n "
	}

	str += " }"

	return str
}

// returns interface for Filters
func getFiltersInterface(spec *DaoSpec, authorSpec *AuthorSpec, columns []ColumnSchema) string {
	str := "{}"
	return str
}

// returns model for filters
func getFiltersModel(spec *DaoSpec, authorSpec *AuthorSpec, columns []ColumnSchema) string {
	str := ""

	for _, elem := range columns {
		str += elem.COLUMN_NAME + ": { value: " + "\"" + elem.COLUMN_NAME_ALIAS + "" + " = ?" + "\"},\n "
	}

	return str
}

// Converts MySQL datatypes to typescript types
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
	util.HandleException(err)
	templatePath := cwd + "/templates/dao.template.txt"
	t, err := template.ParseFiles(templatePath)
	util.HandleException(err)

	outputFilePath := cwd + "/output/" + config.ClassName + ".dao.ts"

	f, err := os.Create(outputFilePath)
	util.HandleException(err)
	defer f.Close()

	err = t.Execute(f, config)

	util.HandleException(err)
}
