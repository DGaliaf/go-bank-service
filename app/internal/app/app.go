package app

import (
	pg "avito-tech/app/adapters/db/postgresql"
	_ "avito-tech/app/docs"
	"avito-tech/app/internal/config"
	v1 "avito-tech/app/internal/controller/http/v1"
	"avito-tech/app/internal/domain/service/user"
	"avito-tech/app/pkg/client/postgresql"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"time"
)

type App struct {
	cfg *config.Config

	router     *httprouter.Router
	httpServer *http.Server

	pgClient *pgx.Conn
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	log.Println("router initializing")
	router := httprouter.New()

	log.Println("swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	pgConfig := postgresql.NewPgConfig(
		config.PostgreSQL.Username, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)

	pgClient, err := postgresql.NewClient(ctx, 5, time.Second*5, pgConfig)
	if err != nil {
		log.Fatalln(err)
	}

	userStorage := pg.NewUserStorage(pgClient)
	userService := user.NewUserService(userStorage)
	userHandler := v1.NewUserHandler(userService)
	userHandler.Register(router)

	return App{
		cfg:      config,
		router:   router,
		pgClient: pgClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, _ := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startHTTP()
	})

	return grp.Wait()
}

func (a *App) startHTTP() error {
	log.Println("HTTP Server initializing")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	if err != nil {
		log.Fatalln("failed to create listener")
	}

	a.httpServer = &http.Server{
		Handler:      a.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Println("server shutdown")
		default:
			log.Fatal(err)
		}
	}

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return err
}
