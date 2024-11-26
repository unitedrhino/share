package conf

const (
	Mysql    = "mysql"
	Pgsql    = "pgsql"
	Sqlite   = "sqlite"
	Tdengine = "tdengine"
)

type Database struct {
	DBType      string `json:",default=mysql,options=mysql|pgsql|sqlite"`          //
	Driver      string `json:",default=taosWS,options=taosRestful|taosWS|taosSql"` //
	IsInitTable bool   `json:",default=true"`
	//IsInitTable bool   `json:",default=true"`
	DSN string `json:""` //dsn
}

// 时序数据库（Time Series Database）
type TSDB struct {
	DBType string `json:",default=mysql,env=tsDBType,options=mysql|pgsql|sqlite"`            //
	Driver string `json:",default=taosWS,env=tsDBDriver,options=taosRestful|taosWS|taosSql"` //
	DSN    string `json:",env=tsDBDSN"`                                                      //dsn
}
