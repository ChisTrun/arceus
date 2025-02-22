package ent

import (
	"context"
	"fmt"
	"time"

	"database/sql/driver"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	database "arceus/pkg/database/api"
	"arceus/pkg/logger/pkg/logging"
)

func NewDriver(config *database.Database) driver.Driver {
	drv := &Driver{config: config}
	if config.GetAuthMethod() == database.Database_AUTH_METHOD_AWS_IAM {
		drv.startRotation()
	}

	return drv
}

type Driver struct {
	drv    mysql.MySQLDriver
	config *database.Database
	token  string
}

func (r *Driver) Open(_ string) (driver.Conn, error) {
	var dbEndpoint = fmt.Sprintf("%s:%d", r.config.GetHost(), r.config.GetPort())

	mysqlConfig := &mysql.Config{
		Addr:                    dbEndpoint,
		DBName:                  r.config.GetName(),
		Net:                     "tcp",
		AllowCleartextPasswords: true,
		AllowNativePasswords:    true,
		ParseTime:               true,
		User:                    r.config.GetUsername(),
	}

	if r.config.GetAuthMethod() == database.Database_AUTH_METHOD_AWS_IAM {
		mysqlConfig.Passwd = r.token
		mysqlConfig.TLSConfig = "rds"
	} else if r.config.GetAuthMethod() == database.Database_AUTH_METHOD_USERNAME_PASSWORD {
		mysqlConfig.Passwd = r.config.GetPassword()
	}
	return r.drv.Open(mysqlConfig.FormatDSN())
}

func (r *Driver) startRotation() {
	// build token for the first time
	if err := r.buildTokenWithRetry(5); err != nil {
		logging.Logger(context.TODO()).Fatal("could not build auth token", zap.Error(err))
	}

	// start token rotation
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			if err := r.buildTokenWithRetry(5); err != nil {
				// exit if token cannot be refresh token
				logging.Logger(context.TODO()).Fatal("could not build auth token", zap.Error(err))
			}

			// rotate after 10 minutes (before current token expired)
		}
	}()
}

func (r *Driver) buildTokenWithRetry(numRetries int) error {
	for {
		if err := r.buildToken(); err != nil {
			if numRetries == 0 {
				return err

			}
			numRetries--
			// retry after 10 seconds
			time.Sleep(10 * time.Second)
			continue
		}

		break
	}
	return nil
}

func (r *Driver) buildToken() error {
	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logging.Logger(context.TODO()).Error("could not load default config", zap.Error(err))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := auth.BuildAuthToken(
		ctx,
		fmt.Sprintf("%s:%d", r.config.GetHost(), r.config.GetPort()),
		r.config.GetAwsRegion(),
		r.config.GetUsername(),
		defaultConfig.Credentials)
	if err != nil {
		logging.Logger(context.TODO()).Error("could not build auth token", zap.Error(err))
		return err
	}

	r.token = token

	return nil
}
