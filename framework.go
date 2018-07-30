package brief_framework

import (
	"brief_framework/config"
	"brief_framework/logger"
	"brief_framework/plugin"
	"brief_framework/util"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Serve() {
	//check older log, to delete
	util.Clean()
	logger.Instance().Info("Init complete, start Servant ...")

	serveAddr, err := config.Instance().GetValue(config.RunningMode(), "serve_addr")
	if err != nil {
		logger.Instance().Warn("Serve get serve_addr config error, msg = %v", err)
		serveAddr = ":8089"
	}

	logger.Instance().Info("Server address is %s.", serveAddr)
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(plugin.LoggerM())
	router.Use(gin.Recovery()) //for Recovery

	plugin.InitHandler(router)

	srv := &http.Server{
		Addr:    serveAddr,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			logger.Instance().Error("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Instance().Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Instance().Error("Server Shutdown:", err)
		panic(err)
	}
	time.Sleep(time.Millisecond * 100)
	logger.Instance().Close()
}
