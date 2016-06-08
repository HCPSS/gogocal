package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"gopkg.in/redis.v3"
)

const (
	// BANNER tells people about this application.
	BANNER = `   _____        _____        _____      _
  / ____|      / ____|      / ____|    | |
 | |  __  ___ | |  __  ___ | |     __ _| |
 | | |_ |/ _ \| | |_ |/ _ \| |    / _` + "`" + ` | |
 | |__| | (_) | |__| | (_) | |___| (_| | |
  \_____|\___/ \_____|\___/ \_____\__,_|_|

A Google Calendar integration.
Version: %s

`
	// VERSION of the application
	VERSION = "v0.1.0"
)

var (
	keyfile   string
	redisAddr string
	redisPass string
	redisDB   int64

	version bool
)

func init() {
	// Parse user input
	flag.StringVar(&keyfile, "k", "key.json", "Specify the Google key file.")
	flag.StringVar(&redisAddr, "a", "redis:6379", "Specify the redis address.")
	flag.StringVar(&redisPass, "p", "", "Specify the redis password.")
	flag.Int64Var(&redisDB, "d", 0, "Specify the redis database index.")

	flag.BoolVar(&version, "version", false, "Print the version.")
	flag.BoolVar(&version, "v", false, "Print the version.")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if version {
		fmt.Printf("%s\n", VERSION)
		os.Exit(0)
	}
}

func main() {
	var app GoGoCal

	// Redis Client
	rc := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       redisDB,
	})
	app.SetRedisClient(rc)

	// Calendar Repo
	srv, err := newCalendarService(keyfile)
	if err != nil {
		panic(err)
	}
	app.SetCalendarService(srv)

	// Set the application logger
	logger := log.New(os.Stdout, "", log.LstdFlags)
	app.SetLogger(logger)

	fmt.Println("GoGoCal is monitoring...")
	app.Run()
}

// newCalendarService returns a new calendar service from a key file.
func newCalendarService(keyfile string) (*calendar.Service, error) {
	// Client
	data, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(
		data,
		"https://www.googleapis.com/auth/calendar",
	)
	if err != nil {
		return nil, err
	}

	client := conf.Client(oauth2.NoContext)

	// Service
	srv, err := calendar.New(client)
	if err != nil {
		return nil, err
	}

	return srv, nil
}
