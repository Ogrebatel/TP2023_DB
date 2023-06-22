package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	"TP2023_DB/app/cmd/server"

	forumDelivery "TP2023_DB/app/internal/forum/delivery"
	forumRep "TP2023_DB/app/internal/forum/repository"
	forumUsecase "TP2023_DB/app/internal/forum/usecase"
	postDelivery "TP2023_DB/app/internal/post/delivery"
	postRep "TP2023_DB/app/internal/post/repository"
	postUsecase "TP2023_DB/app/internal/post/usecase"
	serviceDelivery "TP2023_DB/app/internal/service/delivery"
	serviceRep "TP2023_DB/app/internal/service/repository"
	serviceUsecase "TP2023_DB/app/internal/service/usecase"
	threadDelivery "TP2023_DB/app/internal/thread/delivery"
	threadRep "TP2023_DB/app/internal/thread/repository"
	threadUsecase "TP2023_DB/app/internal/thread/usecase"
	userDelivery "TP2023_DB/app/internal/user/delivery"
	userRep "TP2023_DB/app/internal/user/repository"
	userUsecase "TP2023_DB/app/internal/user/usecase"
)

var cfgPg = postgres.Config{DSN: "host=localhost user=db_pg password=db_postgres database=db_forum port=5432"}

func main() {
	db, err := gorm.Open(postgres.New(cfgPg),
		&gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	forumDB := forumRep.New(db)
	userDB := userRep.New(db)
	postDB := postRep.New(db)
	threadDB := threadRep.New(db)
	serviceDB := serviceRep.New(db)

	forumUC := forumUsecase.New(forumDB, userDB)
	userUC := userUsecase.New(userDB)
	postUC := postUsecase.New(postDB, userDB, threadDB, forumDB)
	threadUC := threadUsecase.New(threadDB, userDB, forumDB)
	serviceUC := serviceUsecase.New(serviceDB)

	e := echo.New()

	e.Logger.SetHeader(`time=${time_rfc3339} level=${level} prefix=${prefix} ` +
		`file=${short_file} line=${line} message:`)
	e.Logger.SetLevel(elog.INFO)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `time=${time_custom} remote_ip=${remote_ip} ` +
			`host=${host} method=${method} uri=${uri} user_agent=${user_agent} ` +
			`status=${status} error="${error}" ` +
			`bytes_in=${bytes_in} bytes_out=${bytes_out}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	e.Use(middleware.Recover())

	forumDelivery.NewDelivery(e, forumUC)
	userDelivery.NewDelivery(e, userUC)
	postDelivery.NewDelivery(e, postUC)
	threadDelivery.NewDelivery(e, threadUC)
	serviceDelivery.NewDelivery(e, serviceUC)

	s := server.NewServer(e)
	if err := s.Start(); err != nil {
		e.Logger.Fatal(err)
	}
}
