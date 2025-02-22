package ent

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"entgo.io/ent/dialect"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/go-sql-driver/mysql"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"

	carbon "arceus/pkg/database/api"
)

func Open(name string, cfg *carbon.Database) (*entsql.Driver, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	if cfg.AuthMethod == carbon.Database_AUTH_METHOD_AWS_IAM {
		pem, err := loadAwsRDSCAPem()
		if err != nil {
			return nil, err
		}
		rootCertPool := x509.NewCertPool()
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, errors.New("failed to append AWS RDS CA")
		}
		err = mysql.RegisterTLSConfig("rds", &tls.Config{
			RootCAs: rootCertPool,
		})
	}
	var db *sql.DB

	driver := NewDriver(cfg)
	sql.Register(name, driver)
	if cfg.TracingEnabled {
		sqltrace.Register(name, driver, sqltrace.WithServiceName(os.Getenv("DD_SERVICE")))
		db, err = sqltrace.Open(name, "", sqltrace.WithServiceName(os.Getenv("DD_SERVICE")))
		if err != nil {
			return nil, err
		}
	} else {
		db, err = sql.Open(name, "")
		if err != nil {
			return nil, err
		}
	}
	drv := entsql.OpenDB(dialect.MySQL, db)
	if cfg.GetMaxIdleConns() > 0 {
		drv.DB().SetMaxIdleConns(int(cfg.GetMaxIdleConns()))
	}
	if cfg.GetMaxOpenConns() > 0 {
		drv.DB().SetMaxOpenConns(int(cfg.GetMaxOpenConns()))
	}
	if cfg.GetConnMaxIdleTime() > 0 {
		drv.DB().SetConnMaxIdleTime(time.Duration(cfg.GetConnMaxIdleTime()) * time.Minute)
	}
	if cfg.GetConnMaxLifeTime() > 0 {
		drv.DB().SetConnMaxLifetime(time.Duration(cfg.GetConnMaxLifeTime()) * time.Minute)
	}
	return drv, nil
}

type DriverOptions struct {
	Cfg    *carbon.Database
	Statsd *statsd.Client
}

func NewOptions(cfg *carbon.Database, statsd *statsd.Client) *DriverOptions {
	return &DriverOptions{
		Cfg:    cfg,
		Statsd: statsd,
	}
}

func OpenWithOption(name string, opts *DriverOptions) (*entsql.Driver, *DBStatsD, error) {
	drv, err := Open(name, opts.Cfg)
	if err != nil {
		return nil, nil, err
	}
	if opts.Statsd != nil {
		statsd := NewDBStatsD(name, opts.Statsd, drv)
		statsd.Start()
		return drv, statsd, nil
	}
	return drv, nil, nil
}

const _pem = "https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem"

func loadAwsRDSCAPem() ([]byte, error) {
	resp, err := http.Get(_pem)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
