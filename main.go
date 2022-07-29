package main

import (
	"context"
	"net/http"
	"log"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"
	chart "Sneakers/chart"
	//"Sneakers/scraper"
	//"io/ioutil"

	"github.com/julienschmidt/httprouter"
	//charts "github.com/wcharczuk/go-chart/v2"
)



func newRouter() *httprouter.Router {
	mux := httprouter.New()

	mux.GET("/sneaker", getChannelStats())

	return mux
}

func getChannelStats() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		
		chart.ChartIt(w, r)
	}
}

func main() {
	//http.HandleFunc(makeChart())
	
	srv := &http.Server{
		Addr:    ":10101",
		Handler: newRouter(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		log.Println("service interrupt received")

		log.Println("http server shutting down")
		time.Sleep(5 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error: %v", err)
		}

		log.Println("shutdown complete")

		close(idleConnsClosed)

	}()

	log.Printf("Starting server on port 10101")
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fatal http server failed to start: %v", err)
		}
	}

	<-idleConnsClosed
	log.Println("Service Stop")

}
