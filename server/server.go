package server

import (
  "log"
  "net/http"
  "github.com/bradbeam/fun/config"
  "github.com/bradbeam/fun/server/actions"
  "github.com/bradbeam/fun/fundb"
  "database/sql"
  // "strconv"
)

// type Monsters []characters.Character
// type Players  []characters.Character

func Serve(configuration config.Config) {
    mydb, err := fundb.Connect(configuration)
    if err != nil {
      log.Fatal(err)
    }

    defer mydb.Close()

    // Set up endpoints
    http.Handle("/action/attack", makeHandler(actions.Attack, mydb))

    log.Fatal(http.ListenAndServe(":2000", nil))
}

// This is the fugliest function ever.
func makeHandler(fn func(http.ResponseWriter, *http.Request, *sql.DB), mydb *sql.DB ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, mydb)
	}
}
