package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	muxtrace "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux"
	otelglobal "go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/stdout"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	klog "gitdev.inno.ktb/coach/thaichanabe/log"
	"gitdev.inno.ktb/coach/thaichanabe/place"
)

var (
	buildcommit = "development"
	buildtime   = time.Now().Format(time.RFC3339)

	tracer = otelglobal.Tracer("mux-server")
)

func main() {
	flag.Bool("migrate", false, "migrate database")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetDefault("app.port", "8000")
	viper.SetDefault("db.conn", "ktb@ktbserver:Passw0rd@tcp(ktbserver.mysql.database.azure.com)/thaichana?charset=utf8&parseTime=True&loc=Local")

	viper.SetConfigName("config")         // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		log.Printf("warning error config file: %s \n", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	initTracer()

	db, err := gorm.Open("mysql", viper.GetString("db.conn"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db = db.Debug()

	if viper.GetBool("migrate") {
		db.AutoMigrate(place.CheckIn{})
		return
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	client := &http.Client{}

	r := mux.NewRouter()
	r.Use(klog.Middleware(logger))
	r.Use(muxtrace.Middleware("my-server"))
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Strict-Transport-Security", "max-age=604800; includeSubDomains; preload")
			handler.ServeHTTP(w, r)
		})
	})

	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/ok", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"commit":    buildcommit,
			"timestamp": buildtime,
		})
	}))
	r.Handle("/checkin", place.CheckInHandler(db, client)).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/places", place.Handler(db)).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/checkout", place.CheckOutHandler(db)).Methods(http.MethodPost, http.MethodOptions)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%s", viper.GetString("app.port")),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func initTracer() {
	exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}
	cfg := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}
	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(cfg),
		sdktrace.WithSyncer(exporter),
	)
	if err != nil {
		log.Fatal(err)
	}
	otelglobal.SetTraceProvider(tp)
}
