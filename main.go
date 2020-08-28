package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	klog "gitdev.inno.ktb/coach/thaichanabe/log"
	"gitdev.inno.ktb/coach/thaichanabe/place"
)

func main() {
	flag.Bool("migrate", false, "migrate database")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetDefault("app.port", "8080")
	viper.SetDefault("db.conn", "root:my-secret-pw@/thaichana?charset=utf8&parseTime=True&loc=Local")

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

	r.Handle("/checkin", place.CheckInHandler(db, client))
	r.HandleFunc("/places", place.Handler(db))
	r.HandleFunc("/checkout", place.CheckOutHandler(db))

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("127.0.0.1:%s", viper.GetString("app.port")),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
