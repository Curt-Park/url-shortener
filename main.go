package main

import (
	"flag"
	"log"
	"net/http"

	internal "github.com/Curt-Park/url-shortener/internal"

	_ "github.com/Curt-Park/url-shortener/docs"

	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "go.uber.org/automaxprocs"
)

var (
	port    string
	profile bool
)

// @title         URL Shortener.
// @description   profiling - http://localhost:8080/debug/pprof/
// @contact.name  Curt-Park
// @contact.email www.jwpark.co.kr@gmail.com
func main() {
	// Parse the args.
	flag.StringVar(&port, "port", "8080", "Service Port. Default: 10000")
	flag.BoolVar(&profile, "profile", false, "Enable profliling.")
	flag.Parse()

	// Create a server with echo.
	e := echo.New()
	if profile {
		pprof.Register(e)
		log.Println("Profiler On")
	}

	// Enable metrics middleware.
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// Middlewares.
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// APIs
	e.GET("/", healthcheck)
	e.POST("/shorten", internal.ShortenURL)
	e.GET("/:key", internal.OriginalURL)

	// Start the server
	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":" + port))
}

// @Summary     Healthcheck
// @Description It returns true if the api server is alive.
// @Accept      json
// @Produce     json
// @Success     200 {object} bool "API server's liveness"
// @Router      / [get].
func healthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, true)
}
