package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-openapi/loads"
	"github.com/vkajtazov/notes/notes-api/gen/restapi"
	"github.com/vkajtazov/notes/notes-api/gen/restapi/operations"
)

var (
	port int
)

func initFlags() {
	flag.IntVar(&port, "port", 8080, "api port number")
	flag.Parse()
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	initFlags()

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		fmt.Println(err)
	}
	api := operations.NewNotesAPI(swaggerSpec)

	server := restapi.NewServer(api)
	server.Port = port

	server.EnabledListeners = []string{"http"}
	server.ConfigureAPI()
	defer server.Shutdown()
	go func() {
		// serve API
		if err := server.Serve(); err != nil {
			panic(fmt.Sprintf("ERROR: %v\r\n", err))
		}
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("Running...")
	<-done
}
