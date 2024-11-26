package conf

const (
	Mysql    = "mysql"
	Pgsql    = "pgsql"
	Sqlite   = "sqlite"
	Tdengine = "tdengine"
)

type Database struct {
	DBType string `json:",default=mysql,env=dbType,options=mysql|pgsql|sqlite"` //
	//IsInitTable bool   `json:",default=false"`
	IsInitTable bool   `json:",env=dbIsInitTable,default=true"`
	DSN         string `json:",env=dbDSN"` //dsn
}

// 时序数据库（Time Series Database）
type TSDB struct {
	DBType string `json:",default=mysql,env=tsDBType,options=mysql|pgsql|sqlite"`            //
	Driver string `json:",default=taosWS,env=tsDBDriver,options=taosRestful|taosWS|taosSql"` //
	DSN    string `json:",env=tsDBDSN"`                                                      //dsn
}
