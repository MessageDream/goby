package model

import (
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var (
	ErrDataNotExist     = errors.New("Data does not exist")
	ErrDataAlreadyExist = errors.New("Data already exist")
)

// Engine represents a xorm engine or session.
type Engine interface {
	Delete(interface{}) (int64, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Find(interface{}, ...interface{}) error
	Get(interface{}) (bool, error)
	Id(interface{}) *xorm.Session
	In(string, ...interface{}) *xorm.Session
	Insert(...interface{}) (int64, error)
	InsertOne(interface{}) (int64, error)
	Iterate(interface{}, xorm.IterFunc) error
	Sql(string, ...interface{}) *xorm.Session
	Table(interface{}) *xorm.Session
	Where(interface{}, ...interface{}) *xorm.Session
}

type JsonObjectConverter interface {
	Convert() interface{}
}

type ModelProtocol interface {
	Update(engine Engine, specFields ...string) error

	Create(engine Engine) error

	Get() (bool, error)

	Delete(engine Engine) error

	Exist() (bool, error)
}

type DbCfg struct {
	Type, Host, Name, User, Passwd, Path, SSLMode, LogPath string
}

var (
	x         *xorm.Engine
	tables    []interface{}
	HasEngine bool

	dbCfg *DbCfg

	EnableSQLite3 bool
	EnableTiDB    bool
)

func init() {

	//bug fix for:session/cache(start): gob: name not registered for interface
	gob.Register(&User{})

	tables = append(tables,
		new(User),
		new(UserToken),
		new(App),
		new(Collaborator),
		new(Deployment),
		new(DeploymentVersion),
		new(DeploymentHistory),
		new(Package),
		new(PackageDiff),
		new(PackageMetrics))

	gonicNames := []string{"SSL"}
	for _, name := range gonicNames {
		core.LintGonicMapper[name] = true
	}
}

func InitDBConfig(config *DbCfg) {
	dbCfg = config
}

// parsePostgreSQLHostPort parses given input in various forms defined in
// https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING
// and returns proper host and port number.
func parsePostgreSQLHostPort(info string) (string, string) {
	host, port := "127.0.0.1", "5432"
	if strings.Contains(info, ":") && !strings.HasSuffix(info, "]") {
		idx := strings.LastIndex(info, ":")
		host = info[:idx]
		port = info[idx+1:]
	} else if len(info) > 0 {
		host = info
	}
	return host, port
}

func getEngine() (*xorm.Engine, error) {
	connStr := ""
	param := "?"
	if strings.Contains(dbCfg.Name, param) {
		param = "&"
	}
	switch dbCfg.Type {
	case "mysql":
		if dbCfg.Host[0] == '/' { // looks like a unix socket
			connStr = fmt.Sprintf("%s:%s@unix(%s)/%s%scharset=utf8&parseTime=true",
				dbCfg.User, dbCfg.Passwd, dbCfg.Host, dbCfg.Name, param)
		} else {
			connStr = fmt.Sprintf("%s:%s@tcp(%s)/%s%scharset=utf8&parseTime=true",
				dbCfg.User, dbCfg.Passwd, dbCfg.Host, dbCfg.Name, param)
		}
	case "postgres":
		host, port := parsePostgreSQLHostPort(dbCfg.Host)
		if host[0] == '/' { // looks like a unix socket
			connStr = fmt.Sprintf("postgres://%s:%s@:%s/%s%ssslmode=%s&host=%s",
				url.QueryEscape(dbCfg.User), url.QueryEscape(dbCfg.Passwd), port, dbCfg.Name, param, dbCfg.SSLMode, host)
		} else {
			connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s%ssslmode=%s",
				url.QueryEscape(dbCfg.User), url.QueryEscape(dbCfg.Passwd), host, port, dbCfg.Name, param, dbCfg.SSLMode)
		}
	case "sqlite3":
		if !EnableSQLite3 {
			return nil, errors.New("This binary version does not build support for SQLite3.")
		}
		if err := os.MkdirAll(path.Dir(dbCfg.Path), os.ModePerm); err != nil {
			return nil, fmt.Errorf("Fail to create directories: %v", err)
		}
		connStr = "file:" + dbCfg.Path + "?cache=shared&mode=rwc"
	case "tidb":
		if !EnableTiDB {
			return nil, errors.New("This binary version does not build support for TiDB.")
		}
		if err := os.MkdirAll(path.Dir(dbCfg.Path), os.ModePerm); err != nil {
			return nil, fmt.Errorf("Fail to create directories: %v", err)
		}
		connStr = "goleveldb://" + dbCfg.Path
	default:
		return nil, fmt.Errorf("Unknown database type: %s", dbCfg.Type)
	}
	return xorm.NewEngine(dbCfg.Type, connStr)
}

func NewTestEngine(x *xorm.Engine) (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("Connect to database: %v", err)
	}

	x.SetMapper(core.GonicMapper{})
	return x.StoreEngine("InnoDB").Sync2(tables...)
}

func SetEngine() (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("Fail to connect to database: %v", err)
	}

	x.SetMapper(core.GonicMapper{})

	// WARNING: for serv command, MUST remove the output to os.stdout,
	// so use log file to instead print to stdout.
	logPath := path.Join(dbCfg.LogPath, "xorm.log")
	os.MkdirAll(path.Dir(logPath), os.ModePerm)

	f, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("Fail to create xorm.log: %v", err)
	}
	x.SetLogger(xorm.NewSimpleLogger(f))
	x.ShowSQL(true)
	return nil
}

func NewEngine() (err error) {
	if err = SetEngine(); err != nil {
		return err
	}

	// if err = migrations.Migrate(x); err != nil {
	// 	return fmt.Errorf("migrate: %v", err)
	// }

	if err = x.StoreEngine("InnoDB").Sync2(tables...); err != nil {
		return fmt.Errorf("sync database struct error: %v\n", err)
	}

	return nil
}

func Transaction(action func(sess *xorm.Session) (interface{}, error)) (interface{}, error) {
	sess := x.NewSession()

	defer func() {
		if !sess.IsCommitedOrRollbacked {
			sess.Rollback()
		}
		sess.Close()
	}()

	if err := sess.Begin(); err != nil {
		return nil, err
	}
	if action != nil {
		ret, err := action(sess)
		if err != nil {
			sess.Rollback()
			return nil, err
		}
		return ret, sess.Commit()
	}
	return nil, errors.New("No action to excute")
}

func EngineGenerate(engine Engine) Engine {
	if engine != nil {
		return engine
	}
	return x
}

func ModelCreate(engine Engine, model ModelProtocol) error {

	engine = EngineGenerate(engine)
	if _, err := engine.Insert(model); err != nil {
		return err
	}
	return nil
}

func ModelGet(id uint64, model ModelProtocol) (bool, error) {
	if id > 0 {
		return x.Id(id).Get(model)
	}
	return x.Where("id != ?", 0).Get(model)
}

func ModelUpdate(engine Engine, id uint64, model ModelProtocol, specFields ...string) error {
	engine = EngineGenerate(engine)
	se := engine.Id(id)
	if specFields != nil && len(specFields) > 0 {
		se = se.Cols(specFields...)
	} else {
		se = se.AllCols()
	}
	if _, err := se.Update(model); err != nil {
		return err
	}
	return nil
}

func Ping() error {
	return x.Ping()
}

// DumpDatabase dumps all data from database to file system.
func DumpDatabase(filePath string) error {
	return x.DumpAllToFile(filePath)
}
