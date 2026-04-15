package Database

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"quay-go-api/Services/Logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	dbType := os.Getenv("DB_TYPE")
	dsn := os.Getenv("DB_DSN")

	// Checks
	if dbType == "" {
		Logger.Error("DB_TYPE environment variable is not set")
		panic("DB_TYPE environment variable is not set")
	}
	if dbType != "postgres" && dbType != "mysql" {
		Logger.Error("Unsupported DB_TYPE: '" + dbType + "'. Supported types are: 'postgres', 'mysql'")
		panic("Unsupported DB_TYPE: '" + dbType + "'. Supported types are: 'postgres', 'mysql'")
	}
	if dsn == "" {
		Logger.Error("DB_DSN environment variable is not set")
		Logger.Info("Ensure that DSN are in format: postgres://user:password@localhost:5432/dbname or  user:password@tcp(localhost:3306)/dbname")
		panic("DB_DSN environment variable is not set")
	}

	// Define Gorm configuration (naming strategy)
	namingStrategy := schema.NamingStrategy{
		TablePrefix:   "",   // no prefix for table names
		SingularTable: true, // use singular table names (e.g. "user" instead of "users")
		NameReplacer:  nil,  // no name replacer
		// By default Gorm uses snake_case for column names, so we don't need to set it explicitly
	}

	switch dbType {
	case "postgres":
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: namingStrategy})
	case "mysql":
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: namingStrategy})
	default:
		Logger.Error("Unsupported DB_TYPE: " + dbType + ". Supported types are: 'postgres', 'mysql'")
	}

	if err != nil {
		Logger.Error("Failed to connect to database: " + err.Error())
		panic("Failed to connect to database: " + err.Error())
	}
}
