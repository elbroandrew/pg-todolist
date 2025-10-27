package server

//данный файл нужен для Graceful shutdown, если происходит остановка сервиса через ctrl+C, либо по неизвестно причине.
//Неизвестная причина может кинуть панику, что перехватит gin.Recovery. Который создан через gin.Default.

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//приложение, зависимости которого (redis cache, MySQL DB) надо закрыть
type App struct {
	HttpServer *http.Server
	Closers []func() error  //слайс функций для закрытия ресурсов (БД, redis)
}

func NewApp(handler http.Handler, port string) *App {
	return &App{
		HttpServer: &http.Server{
			Addr: ":"+port,
			Handler: handler,
		},
		Closers: []func() error{},
	}
}

// метод закрывающий ресурс
func (a *App) AddCloser(closer func() error) {
	a.Closers = append(a.Closers, closer)
}

// Run запускает приложение и делает graceful shutdown, когда необходимо
func (a *App) Run() {
	//запускаю сервер в горутине, чтобы не блокировать код в main, это позволит коду в main дойти до 
	// той части, где он ждёт сигнала о завершении (<-quit) и выолнить логику Graceful Shutdown.
	go func(){
		log.Printf("Server is running on port %s", a.HttpServer.Addr)
		if err := a.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//ожидаю сигнал остановки Sigint (ctrl+C), код 2
	quit := make(chan os.Signal, 1) // 1 - размер буфера, рекоммендация разработчика
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	//Даю время на завершение запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.HttpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	//Закрываю ресурсы
	for _, closer := range a.Closers {
		if err := closer(); err != nil {
			log.Printf("Error closing resource: %v", err)
		}
	}
	log.Println("Server exiting")
}