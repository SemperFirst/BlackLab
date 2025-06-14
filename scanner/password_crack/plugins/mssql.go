package plugins
import (
	"SafeDp/scanner/password_crack/models"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"database/sql"
)

func ScanMssql(service models.Service) (result models.ScanResult, err error) {
	result.Service = service
	dataSourceName := fmt.Sprintf("server=%v;port=%v;user id=%v;password=%v;database=%v", service.Ip, service.Port, service.Username, service.Password, "master")
	db, err := sql.Open("mssql", dataSourceName)
	if err != nil {
		return result, err
	}
	err = db.Ping()
	if err != nil {
		return result, err
	}
	result.Result = true
	defer func() {
		if db != nil {
			_ = db.Close()
		}
	}()
	return result, nil
}