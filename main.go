package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/savvy-bit/gin-react-postgres/config"
	"github.com/savvy-bit/gin-react-postgres/database"
	"github.com/savvy-bit/gin-react-postgres/router"
)

func main() {
	// Initialize database
	database.Init()

	addr := flag.String("addr", config.Global.Server.Port, "Address to listen and serve")
	flag.Parse()

	gin.SetMode(config.Global.Server.Mode)
	app := gin.Default()

	app.Static("/images", filepath.Join(config.Global.Server.StaticDir, "img"))
	app.StaticFile("/favicon.ico", filepath.Join(config.Global.Server.StaticDir, "img/favicon.ico"))
	app.MaxMultipartMemory = config.Global.Server.MaxMultipartMemory << 20

	router.Route(app)

	// Listen and Serve
	if err := app.Run(*addr); err != nil {
		log.Fatal(err.Error())
	}
}
