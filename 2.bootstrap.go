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
	dbPingTries          = 5
	dbPingPause          = 200 * time.Millisecond
)

func Bootstrap() {
	var portFlag, dbNameFlag, dbUserFlag, dbPassFlag string
	flag.StringVar(&portFlag, `port`, ``, Str(`Web port number (or $`, envPort, `)`))
	flag.StringVar(&dbNameFlag, `db`, ``, Str(`Database name (or $`, envDBName, `)`))
	flag.StringVar(&dbUserFlag, `user`, ``, Str(`Database user (or $`, envDBUser, `)`))
	flag.StringVar(&dbPassFlag, `pass`, ``, Str(`Database password (or $`, envDBPass, `)`))
	flag.Parse()

	pick := func(raw ...string) string {
		for _, x := range raw {
			x = strings.TrimSpace(x)
			if x != `` { return x }
		}
		return ``
	}

	port := pick(portFlag, os.Getenv(envPort), `4444`)
	if n, e := strconv.Atoi(port); e != nil || n < 1 || n > 65535 {
		panic(Str(`Invalid port configuration: `, port))
	}

	dbName := pick(dbNameFlag, os.Getenv(envDBName))
	dbUser := pick(dbUserFlag, os.Getenv(envDBUser))
	dbPass := pick(dbPassFlag, os.Getenv(envDBPass))

	missing := make([]string, 0, 3)
	if dbName == `` { missing = append(missing, envDBName) }
	if dbUser == `` { missing = append(missing, envDBUser) }
	if dbPass == `` { missing = append(missing, envDBPass) }
	if len(missing) > 0 {
		panic(Str(`Missing DB configuration: `, strings.Join(missing, `, `)))
	}

	App = App_t{
		DB: OpenDB(dbUser, dbPass, dbName),
		port: port,
		session: Session_t{
			name:     `session`,
			path:     `/`,
			maxAge:   60 * 60 * 24 * 365,
			httpOnly: true,
			secure:   false,
			sameSite: http.SameSiteLaxMode,
		},
	}

	LoadStaticData()
	App.staticVersion = ComputeStaticVersion()
}

func OpenDB(user, pass, name string) *DB_t {
	var lastErr error
	for try := 1; try <= dbPingTries; try++ {
		dbx := Mysql(user, pass, name).
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

func ComputeStaticVersion() string {
	root := `./static`
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
