package sqlx

import (
	"github.com/go-sql-driver/mysql"
)

type XTime struct {
	mysql.NullTime
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
