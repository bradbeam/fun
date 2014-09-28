package fundb

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/bradbeam/fun/config"
  "strconv"
  "log"
)

func Connect(c config.Config) (*sql.DB, error) {
  log.Println("Creating DB connect string")
  var connectionstring string
  connectionstring += c["DatabaseUsername"]
  // Test to see if we actually have a database password
  if _, ok := c["DatabasePassword"]; ok && c["DatabasePassword"] != "" {
    connectionstring += ":" + c["DatabasePassword"]
  }
  connectionstring += "@tcp(" + c["DatabaseHost"] + ":" + c["DatabasePort"] + ")"
  // Test to see if we have a database defined
  if _, ok := c["Database"]; ok {
    connectionstring += "/" + c["Database"]
  } else {
    connectionstring += "/"
  }

  log.Println("Connecting to database")
  db, err := sql.Open(c["DatabaseDriver"], connectionstring)
  if err != nil {
    return nil, err
  }

  if _, ok := c["DatabaseMaxIdleConnections"]; ok {
    maxidlecon, err := strconv.Atoi(c["DatabaseMaxIdleConnections"])
    if err != nil {
      return nil, err
    }
    db.SetMaxIdleConns(maxidlecon)
  }

  if _, ok := c["DatabaseMaxOpenConnections"]; ok {
    maxcon, err := strconv.Atoi(c["DatabaseMaxOpenConnections"])
    if err != nil {
      return nil, err
    }
    db.SetMaxOpenConns(maxcon)
  }

  log.Println("Returning db")
  return db, err
}
