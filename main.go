package main

import (
	"flag"
	"log"
	"net/http"

	internal "github.com/Curt-Park/url-shortener/internal"

	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	port    string
	profile bool
)

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
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// APIs
	e.GET("/", healthcheck)
	e.POST("/shorten", internal.ShortenURL)
	e.GET("/:key", internal.OriginalURL)

	// Start the server
	e.Logger.Fatal(e.Start(":" + port))
}

// Healthcheck API.
func healthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, true)
}
