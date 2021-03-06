package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/256dpi/fire/coal"
	"github.com/256dpi/fire/flame"
)

var debug = getEnv("DEBUG", "no") == "yes"
var mongo = getEnv("MONGODB_URI", "mongodb://localhost/iad-archive")
var secret = getEnv("SECRET", "abcd1234abcd1234")
var port = getEnv("PORT", "8000")

// TODO: Resize images on the fly.

func main() {
	// create store
	store := coal.MustCreateStore(mongo)

	// prepare database
	err := prepareDatabase(store)
	if err != nil {
		panic(err)
	}

	// create main mux
	mux := http.NewServeMux()

	// build v1 api handler
	mux.Handle("/", handler(store, secret, debug))

	// get port
	port, _ := strconv.Atoi(port)

	// run plain server
	fmt.Printf("Running on http://0.0.0.0:%d\n", port)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), mux)
	if err != nil {
		panic(err)
	}
}

func prepareDatabase(store *coal.Store) error {
	// ensure indexes
	err := indexer.Ensure(store)
	if err != nil {
		return err
	}

	// ensure first user
	err = flame.EnsureFirstUser(store, "Root", "root@archive.iad.zhdk.ch", "root")
	if err != nil {
		return err
	}

	// ensure admin application
	adminAppKey, err := flame.EnsureApplication(store, "Admin", "admin", "abcd1234")
	if err != nil {
		return err
	}

	// print main application keys
	fmt.Printf("Admin Application Key: %s\n", adminAppKey)

	return nil
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return def
}
