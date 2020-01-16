package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

var flagListenAddr = flag.String("listenaddr", "127.0.0.1:8080", "http service address, default: localhost")
var flagDirPrefix = flag.String("prefix", "/", "prefix")

//var flagDirPrefixShort = flag.String("p", "/", "prefix")

var flagServerDir = flag.String("filedir", "", "directory to make avaliable via HTTP protocol!")
var flagDirStartupList = flag.Bool("sl", false, "list files in directory on start")

// TODO:
var flagDirList = flag.Bool("l", false, "allow index of folder to be displayed to client")
var flagShutdownControl = flag.String("s", "none", "shutdown control type, param: or url:")

func main() {
	fmt.Printf("[init] QuickAndDirty http server 0.1\n[init] server starting...\n")
	flag.Parse()
	rmux := mux.NewRouter()

	ok, err := exists(*flagServerDir)
	if !ok {
		log.Println("directory does not exist! ", err)
	}

	files, err := ioutil.ReadDir(*flagServerDir)
	if err != nil {
		log.Fatal(err)
	}

	if *flagDirStartupList {
		fmt.Printf("[init] listing dir ...\n")
		for _, f := range files {
			fmt.Println("\t" + f.Name())
		}
	}

	fmt.Printf("[init] %d files in directory\n", len(files))

	rmux.PathPrefix(*flagDirPrefix).Handler(
		http.StripPrefix(*flagDirPrefix,
			http.FileServer(http.Dir(*flagServerDir))))

	fmt.Printf("[init] Starting HTTP server on %s\n", *flagListenAddr)
	srv := &http.Server{
		Handler:      rmux,
		Addr:         *flagListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	interrupt := true
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
			interrupt = false
		}
	}()
	if interrupt {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
	}
	fmt.Printf("[init] Server stopped\n")
	fmt.Printf("[init] Server shutdown\n")
}
