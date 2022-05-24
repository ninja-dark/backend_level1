package main

import (
	"context"
	"goback1/lesson4/reguser/internal/infrastructure/api/defmux"
	"goback1/lesson4/reguser/internal/infrastructure/db/mem/usermemstore"
	"goback1/lesson4/reguser/internal/infrastructure/server"
	"goback1/lesson4/reguser/internal/usecases/app/repos/userrepo"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	ust := usermemstore.NewUsers()
	us := userrepo.NewUsers(ust)
	h := defmux.NewRouter(us)
	srv := server.NewServer(":8000", h)

	srv.Start(us)
	log.Print("Start")
	<-ctx.Done()
	srv.Stop()
	cancel()
	log.Print("Exit")
}
