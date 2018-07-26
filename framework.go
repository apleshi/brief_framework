package brief_framework

import (
	"config"
	"logger"
	"server_plugin"
	"tools"
	_ "schedule"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
)


func Serve() {
	//check older log, to delete
	tools.DoDeleteWork()
	logger.Instance().Info("Init complete, Server start...go")

	serve_addr, err := config.Instance().GetValue("online", "serve_addr")
	if err != nil {
		logger.Instance().Warn("GetValue serve_addr error, err = %v", err)
		serve_addr = ":8081"
	}

	logger.Instance().Info("Server address is %s.", serve_addr)
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(tools.LoggerM())
	router.Use(gin.Recovery()) //for Recovery

	server_plugin.InitHandle(router)

	srv := &http.Server{
		Addr:    serve_addr,
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
