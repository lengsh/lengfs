package main

import (
	"context"
	//"net"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lengsh/findme/utils"
	"github.com/lengsh/lengfs/lfs"
 "github.com/lengsh/lengfs/web"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	//	"strings"
	"syscall"
	"time"
)

var ( // main operation modes
	confFile = flag.String("c", "./local.json", "server configure file.")
	port     = flag.String("p", "8080", "server listen port.")
	node     = flag.String("i", "0", "lengfs global node (iNode=0)")
	sec      = flag.Bool("u", true, "url secret(false)")
	queues   = flag.String("s", "", "lengfs server queues, such as localhost:8080;localhost:8081")
	// viewRoot   = flag.String("v", "./view/", "view template root")
)

func init() {

}
func usage() {
	// Fprintf allows us to print to a specifed file handle or stream
	fmt.Fprintf(os.Stderr, "Usage: %s [-flag xyz]\n", os.Args[0])
	flag.PrintDefaults()
}

func runInit() {
	flag.Usage = usage
	flag.Parse()
	// There is also a mandatory non-flag arguments
	if len(os.Args) < 2 {
		usage()
	}
	cdir, _ := getCurrentPath()
	lfs.LNode.Parent = cdir + "static"
	lfs.LNode.Pnode = "lengfs"
	lfs.LNode.Inode = *node
	lfs.LNode.UrlSecret = *sec
	lfs.LNode.Domain = "lengsh"
	lfs.LNode.Queues = *queues
	lfs.LNode.Port, _ = strconv.Atoi(*port)
	utils.ServerConfig.WebDir = cdir // "./"
	utils.ServerConfig.PassSalt = "xrjoN1qR"

	logs.Debug(*queues)
	os.MkdirAll(lfs.LNode.Parent, 0755)
	fmt.Println(lfs.LNode)
	lfs.StatInit()
        lfs.LfsSetFrequency(600)
}

func getCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	d, _ := filepath.Split(path)
	return d, nil
}

func main() {

	runInit()
	logs.SetLogFuncCall(true)
	web.Router()
	//wait := time.Second * 2
	ctx, cancel := context.WithCancel(context.Background()) // context.WithTimeout(context.Background(), wait)
	// defer cancel()
	sp := fmt.Sprintf(":%s", *port)
	fmt.Println("server listen to ", *port)
	s := &http.Server{
		Addr:         sp, //port,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		//      MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	go lfs.JobWatch(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	//    signal.Notify(ch, os.Interrupt)

	// Block until we receive our signal.
	// Handle SIGINT and SIGTERM.
	<-ch
	//  log.Println(<-ch)
	// Stop the service gracefully.
	// log.Println(s.Shutdown(nil))
	// Wait gorotine print shutdown message

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.

	//	rpcs.MicroRpcServerStop()
	// rpcs.MicNopRpcServerStop()

	s.Shutdown(ctx)
	cancel()
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("graceful shutdown -->done!")
	time.Sleep(10 * time.Second)

	os.Exit(0)

}
