package server

import (
  "log"
  "net/http"
  "github.com/bradbeam/fun/config"
  "github.com/bradbeam/fun/server/actions"
  "github.com/bradbeam/fun/server/info"
  "github.com/bradbeam/fun/fundb"
  "database/sql"
  "encoding/base64"
  "strings"
  // "strconv"
  "expvar"
)

// type Monsters []characters.Character
// type Players  []characters.Character

var (
  FailedAuthCount = expvar.NewInt("FailedAuthCount")
  SuccessfulAuthCount = expvar.NewInt("SuccessfulAuthCount")
  TotalConnections = expvar.NewInt("TotalConnections")
)

func Serve(configuration config.Config) {
    mydb, err := fundb.Connect(configuration)
    if err != nil {
      log.Fatal(err)
    }

    defer mydb.Close()

    // Set up endpoints
    http.Handle("/action/attack", makeHandler(actions.Attack, mydb))
    http.Handle("/info/account", makeHandler(info.Account, mydb))

    log.Fatal(http.ListenAndServe(":2000", nil))
}

// This is the fugliest function ever.
func makeHandler(fn func(http.ResponseWriter, *http.Request, *sql.DB), mydb *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    TotalConnections.Add(1)

    auth := r.Header.Get("Authorization")
    if auth == "" {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      FailedAuthCount.Add(1)
      return
    }

    username, password, ok := parseBasicAuth(auth)
    if !ok {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      FailedAuthCount.Add(1)
      return
    }

    if !checkAuth(username,password,mydb) {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      FailedAuthCount.Add(1)
      return
    }

    SuccessfulAuthCount.Add(1)
		fn(w, r, mydb)
	}
}

// Shoutout + highfive kelsey hightower
// https://code.google.com/p/go/source/detail?r=5e03333d2dcf
// This should get added into the next release of Go, so we'll be able to
// just use it natively.
// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
  if !strings.HasPrefix(auth, "Basic ") {
    return
  }
  c, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
  if err != nil {
    return
  }
  cs := string(c)
  s := strings.IndexByte(cs, ':')
  if s < 0 {
    return
  }
  return cs[:s], cs[s+1:], true
}


func checkAuth(username string, password string, mydb *sql.DB) bool {
  log.Println("Authenticating " + username)
  stmt, err := mydb.Prepare("SELECT IF(password = ?, true, false) AS authenticated " +
                              "FROM fun.account " +
                              "WHERE username = ?;")
  if err != nil {
    return false
  }
  defer stmt.Close()

  // boolean in mysql is backwards from like everything else
  // 0 = false
  // 1 = true
  authenticated := 0
  err = stmt.QueryRow(password, username).Scan(&authenticated)
  if err != nil {
    return false
  }

  if authenticated == 0 {
    return false
  } else {
    return true
  }
}
