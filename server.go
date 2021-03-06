package main

import (
	"github.com/UndeadBigUnicorn/CompanyStatistics/config"
	"github.com/UndeadBigUnicorn/CompanyStatistics/dbworker"
	"github.com/UndeadBigUnicorn/CompanyStatistics/handlers"
	. "github.com/UndeadBigUnicorn/CompanyStatistics/infrastructure/logging"
	"github.com/UndeadBigUnicorn/CompanyStatistics/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
)

func main() {

	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if config.GetSetting("mode").(string) == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	route := gin.Default()

	route.Use(middleware.Auth())

	go func() {

		// company routes
		company := route.Group("/company")
		{
			company.POST("/add", handlers.AddCompany)
			company.GET("/stats", handlers.GetTotalStats)
			company.POST("/update", handlers.UpdateCompany)
		}

		// stats routes
		stats := route.Group("/statistic")
		{
			stats.POST("/add", handlers.AddStats)
			stats.POST("/stats", handlers.GetDetailStats)
		}

		// start server
		route.Run(config.GetSetting("port").(string))

	}()


	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	Info.Println("Stopping the server")
	dbworker.CloseConnection()
	Info.Println("Closing connection to database")
	Info.Println("End of a program")

}
