package main

import (
	"os"
	"prof-filer/generator"

	"github.com/joho/godotenv"
)

// Program Entry
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	authorSpec := generator.AuthorSpec{
		Author: os.Getenv("AUTHOR"),
	}

	daoSpec := generator.DaoSpec{
		Table: os.Getenv("DAO_TABLE_NAME"),
	}

	generator.GenerateDao(&daoSpec, &authorSpec)
}
