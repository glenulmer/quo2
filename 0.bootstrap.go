package main

import (
	"flag"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	. "pm/lib/output"
	. "pm/lib/wrapdb"
)

const (
	envPort              = `PM_PORT`
	envDBName            = `PM_DBNAME`
	envDBUser            = `PM_DBUSER`
	envDBPass            = `PM_DBPASS`
	defaultPort          = `4444`
	defaultSessionName   = `session`
	defaultSessionPath   = `/`
	defaultSessionMaxAge = 60 * 60 * 24 * 365
	dbPingTries          = 5
	dbPingPause          = 200 * time.Millisecond
)

type RuntimeDBConfig_t struct {
	name string
	user string
	pass string
}

type RuntimeConfig_t struct {
	port    string
	db      RuntimeDBConfig_t
	session Session_t
}

func Bootstrap() {
	var db struct{ user, pass, name string }
	var port string
	flag.StringVar(&port, `port`, ``, Str(`Web port number (or $`, envPort, `)`))
	flag.StringVar(&db.name, `db`, ``, Str(`Database name (or $`, envDBName, `)`))
	flag.StringVar(&db.user, `user`, ``, Str(`Database user (or $`, envDBUser, `)`))
	flag.StringVar(&db.pass, `pass`, ``, Str(`Database password (or $`, envDBPass, `)`))
	flag.Parse()

	cfg := LoadRuntimeConfig(port, db.name, db.user, db.pass, nil)
	dbx := OpenRuntimeDB(cfg.db)
	App = App_t{
		DB:            dbx,
		sessionState:  make(map[string]SimpleState_t),
		sessionFilters: make(map[string]FilterState_t),
		sessionCustomers: make(map[string]CustomerState_t),
		sessionEpoch: make(map[string]int),
	}
	LoadStaticData()
	LoadSelectElements()
	App.Port = cfg.port
	App.StaticVersion = ComputeStaticVersion()
	App.Session = cfg.session
}

func OpenRuntimeDB(db RuntimeDBConfig_t) *DB_t {
	var lastErr error
	for try := 1; try <= dbPingTries; try++ {
		dbx := Mysql(db.user, db.pass, db.name).
			Option(`charset`, `utf8`).
			Option(`collation`, `utf8mb4_german2_ci`).
			Open()
		if e := dbx.Ping(); e == nil {
			return dbx
		} else {
			lastErr = e
			_ = dbx.Close()
		}
		if try < dbPingTries { time.Sleep(dbPingPause) }
	}
	panic(Error(`database ping failed: `, lastErr))
}

func LoadRuntimeConfig(portFlag, dbNameFlag, dbUserFlag, dbPassFlag string, getenv func(string) string) (cfg RuntimeConfig_t) {
	if getenv == nil { getenv = os.Getenv }
	cfg.port = PickRuntimeConfig(portFlag, getenv(envPort), defaultPort)
	cfg.db.name = PickRuntimeConfig(dbNameFlag, getenv(envDBName))
	cfg.db.user = PickRuntimeConfig(dbUserFlag, getenv(envDBUser))
	cfg.db.pass = PickRuntimeConfig(dbPassFlag, getenv(envDBPass))
	cfg.session = RuntimeSessionConfig()
	ValidateRuntimePort(cfg.port)
	ValidateDBRuntimeConfig(cfg.db)
	ValidateSessionRuntimeConfig(cfg.session)
	return
}

func PickRuntimeConfig(raw ...string) string {
	for _, x := range raw {
		x = strings.TrimSpace(x)
		if x != `` { return x }
	}
	return ``
}

func ValidateRuntimePort(port string) {
	n, e := strconv.Atoi(port)
	if e != nil || n < 1 || n > 65535 {
		panic(Str(`Invalid port configuration: `, port))
	}
}

func ValidateDBRuntimeConfig(db RuntimeDBConfig_t) {
	missing := make([]string, 0, 3)
	if db.name == `` { missing = append(missing, envDBName) }
	if db.user == `` { missing = append(missing, envDBUser) }
	if db.pass == `` { missing = append(missing, envDBPass) }
	if len(missing) > 0 {
		panic(Str(`Missing DB configuration: `, strings.Join(missing, `, `)))
	}
}

func RuntimeSessionConfig() Session_t {
	return Session_t{
		Name:     defaultSessionName,
		Path:     defaultSessionPath,
		MaxAge:   defaultSessionMaxAge,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

func ValidateSessionRuntimeConfig(x Session_t) {
	if Trim(x.Name) == `` { panic(`Invalid session configuration: name is required`) }
	if Trim(x.Path) == `` { panic(`Invalid session configuration: path is required`) }
	if x.MaxAge <= 0 { panic(`Invalid session configuration: max age must be > 0`) }
	if x.SameSite <= 0 { panic(`Invalid session configuration: same-site is required`) }
}

func ComputeStaticVersion() string {
	return ComputeStaticVersionFromRoot(`./static`)
}

func ComputeStaticVersionFromRoot(root string) string {
	var count, newest, sumSize int64
	e := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil { return err }
		if d.IsDir() { return nil }
		info, e := d.Info()
		if e != nil { return e }
		count++
		sumSize += info.Size()
		mod := info.ModTime().UnixNano()
		if mod > newest { newest = mod }
		return nil
	})
	if e != nil { panic(Error(`could not compute static version from `, root, `: `, e)) }
	if count == 0 { panic(Error(`could not compute static version from `, root, `: no files`)) }
	return Str(count, `-`, sumSize, `-`, newest)
}
