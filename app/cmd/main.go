package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/golang/glog"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/Elderly-AI/scrawler/internal/app/crawler"
	"github.com/Elderly-AI/scrawler/internal/app/tasks"
	"github.com/Elderly-AI/scrawler/internal/pkg/config"
	crawlerfacade "github.com/Elderly-AI/scrawler/internal/pkg/crawler"
	crawlerdb "github.com/Elderly-AI/scrawler/internal/pkg/database/crawler"
	tasksdb "github.com/Elderly-AI/scrawler/internal/pkg/database/tasks"
	common "github.com/Elderly-AI/scrawler/internal/pkg/middleware"
	"github.com/Elderly-AI/scrawler/internal/pkg/model"
	tasksfacade "github.com/Elderly-AI/scrawler/internal/pkg/tasks"
	crawlerdesc "github.com/Elderly-AI/scrawler/pkg/proto/crawler"
)

const taskRunnerDuration = time.Minute * 10

func registerServices(opts Options, s *grpc.Server) Impls {
	crawlerDB := crawlerdb.New(opts.PosgtresConnection)
	crawlerFacade := crawlerfacade.New(crawlerDB)
	crawlerImplementation := crawler.New(crawlerFacade)
	crawlerdesc.RegisterCrawlerServer(s, crawlerImplementation)

	tasksDB := tasksdb.New(opts.PosgtresConnection)
	tasksFacade := tasksfacade.New(tasksDB)
	tasksImplementation := tasks.New(tasksFacade, crawlerImplementation)

	return Impls{
		crawlerImpl: crawlerImplementation,
		tasksImpl:   tasksImplementation,
	}
}

func newGateway(ctx context.Context, conn *grpc.ClientConn, opts []gwruntime.ServeMuxOption) (http.Handler, error) {
	mux := gwruntime.NewServeMux(opts...)

	for _, f := range []func(ctx context.Context, mux *gwruntime.ServeMux, conn *grpc.ClientConn) error{
		crawlerdesc.RegisterCrawlerHandler,
	} {
		if err := f(ctx, mux, conn); err != nil {
			return nil, err
		}
	}
	return mux, nil
}

type Options struct {
	Addr               string
	Mux                []gwruntime.ServeMuxOption
	PosgtresConnection *sqlx.DB
}

type Impls struct {
	crawlerImpl crawler.Implementation
	tasksImpl   tasks.Implementation
}

func createInitialOptions(conf config.Config) Options {
	opts := Options{}
	database, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s user=postgres password=postgres dbname=postgres sslmode=disable", conf.PostgresHost))
	if err != nil {
		glog.Fatal(err)
	}
	opts.PosgtresConnection = database
	opts.Addr = "0.0.0.0:8080"
	return opts
}

func addGRPCMiddlewares(opts Options) Options {
	opts.Mux = []gwruntime.ServeMuxOption{}
	return opts
}

func main() {
	conf := config.GetConfig()
	opts := createInitialOptions(conf)
	opts = addGRPCMiddlewares(opts)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer()
	impls := registerServices(opts, s)
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()
	// register services

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		opts.Addr,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gw, err := newGateway(context.Background(), conn, opts.Mux)
	if err != nil {
		log.Fatalln(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", gw)

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: common.AllowCORS(mux), // TODO add panic middleware
	}

	log.Println("starting background tasks")
	go func() {
		impls.tasksImpl.RunTasks(context.Background(), []model.TaskID{model.TaskIDCheat}, time.Duration(taskRunnerDuration))
	}()

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())

}
