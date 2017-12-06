package dictx

import (
	"github.com/go-sql-driver/mysql"
)

type XTime struct {
	mysql.NullTime
}
