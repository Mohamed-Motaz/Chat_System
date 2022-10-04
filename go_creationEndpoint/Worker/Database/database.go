package Database

import (
	logger "Worker/Logger"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBWrapper struct {
	Db *gorm.DB
}

//return a thread-safe *gorm.DB that can safely be used
//by multiple goroutines
func New() *DBWrapper {
	db := connect()
	logger.LogInfo(logger.DATABASE, logger.ESSENTIAL, "Db setup complete")
	wrapper := &DBWrapper{
		Db: db,
	}

	return wrapper
}

func connect() *gorm.DB {

	dsn := generateDSN(
		DbUser, DbPassword, DbProtocol,
		"", DbHost, DbPort, DbSettings)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.FailOnError(logger.DATABASE, logger.ESSENTIAL, "Unable to connect to db with this error %v", err)
	}
	return db
}

func generateDSN(user, password, protocol, dbName, myHost, myPort, settings string) string {

	return fmt.Sprintf(
		"%v:%v@%v(%v:%v)/%v?%v",
		user, password, protocol, myHost, myPort, dbName, settings)
}
