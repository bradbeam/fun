package fundb

import (
	"database/sql"
	"github.com/bradbeam/fun/config"
	"testing"
	_ "github.com/go-sql-driver/mysql"
)

func createConfig() config.Config {
	configuration := config.Config{
		"DatabaseHost":     "127.0.0.1",
		"Database":         "thisisatestingl1gdatabase",
		"DatabasePort":     "3306",
		"DatabaseUsername": "root",
		"DatabasePassword": "",
		"DatabaseDriver":   "mysql",
		// Don't retain any idle connections
		"DatabaseMaxIdleConnections": "0",
		// Unlimited connections
		"DatabaseMaxOpenConnections": "0",
	}

	return configuration
}

// We'll use this function to
// Seed our test data
func TestInitialization(t *testing.T) {
	// t.Log("Initializing testing database and schema")
	configuration := createConfig()
	// Override default database name
	// Since we'll want to just connect to mysql
	// and then create the database afterwards
	configuration["Database"] = ""
	db, err := Connect(configuration)
	if err != nil {
		t.Error(err)
	}

	// Reset config values to default
	configuration = createConfig()
	_, err = db.Exec("create database if not exists " + configuration["Database"])
	if err != nil {
		t.Error(err)
	}

	// Seed sample data
	_, err = db.Exec("create table if not exists " + configuration["Database"] + ".accounts (name varchar(255), password varchar(255), lastactive varchar(255) )")
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec("insert into " + configuration["Database"] + ".accounts values ( \"myaccountname\", \"mypassword\", \"2014-09-06 03:59:00\")")
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec("create user 'bob'@'localhost' identified by 'password'")
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec("grant all on " + configuration["Database"] + ".* to 'bob'@'localhost'")
	if err != nil {
		t.Error(err)
	}

  _, err = db.Exec("create user 'bob'@'%' identified by 'password'")
  if err != nil {
    t.Error(err)
  }

  _, err = db.Exec("grant all on " + configuration["Database"] + ".* to 'bob'@'%'")
  if err != nil {
    t.Error(err)
  }

	_, err = db.Exec("flush privileges")
	if err != nil {
		t.Error(err)
	}

	// We're done with the initialization
	// So we can go ahead and close our connection
	db.Close()
}

func setupConnection() (*sql.DB, error) {
	configuration := createConfig()
	db, err := Connect(configuration)
	return db, err
}

func TestConnection(t *testing.T) {
	db, err := setupConnection()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Error(err)
	}
}

func TestForBadDatabase(t *testing.T) {
	configuration := createConfig()
	configuration["Database"] = "imabogusdatabase"
	db, err := Connect(configuration)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	err = db.Ping()
	if err == nil {
		t.Error(err)
	}
}

func TestQuery(t *testing.T) {
	db, err := setupConnection()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	var name, version string
	row := db.QueryRow("show variables where variable_name = 'version';")
	err = row.Scan(&name, &version)
	if err != nil {
		t.Error(err)
	}
	if name != "version" {
		t.Error("Pulled back a weird field named version but is not version?")
	}
}

func TestQuery2(t *testing.T) {
	db, err := setupConnection()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	var name string
	row := db.QueryRow("select name from accounts")
	err = row.Scan(&name)
	if err != nil {
		t.Error(err)
	}
	if name != "myaccountname" {
		t.Error("Didn't get myaccountname as a result")
	}
}

func TestInsert(t *testing.T) {
	db, err := setupConnection()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	_, err = db.Exec("insert into accounts values (\"mysecondaccount\",\"mysecondpassword\",\"someuglydate\")")
	if err != nil {
		t.Error(err)
	}

	var name string
	row := db.QueryRow("select name from accounts where lastactive = 'someuglydate'")
	err = row.Scan(&name)
	if err != nil {
		t.Error(err)
	}
	if name != "mysecondaccount" {
		t.Error("Didn't get mysecondaccount as a result")
	}
}

func TestNoPassword(t *testing.T) {
	configuration := createConfig()
	delete(configuration, "DatabasePassword")
	db, err := Connect(configuration)
	if err != nil {
		t.Error(err)
	}
	db.Close()
}

func TestNoDatabase(t *testing.T) {
	configuration := createConfig()
	delete(configuration, "Database")
	db, err := Connect(configuration)
	if err != nil {
		t.Error(err)
	}
	db.Close()
}

func TestNoDriver(t *testing.T) {
	configuration := createConfig()
	configuration["DatabaseDriver"] = "imnotreallyadriver"
	_, err := Connect(configuration)
	if err == nil {
		t.Error("Connected to a database with a bogus driver")
	}
}

func TestInvalidIdleCon(t *testing.T) {
	configuration := createConfig()
	configuration["DatabaseMaxIdleConnections"] = "lulz"
	_, err := Connect(configuration)
	if err == nil {
		t.Error("lulz")
	}
}

func TestInvalidMaxCon(t *testing.T) {
	configuration := createConfig()
	configuration["DatabaseMaxOpenConnections"] = "lulz"
	_, err := Connect(configuration)
	if err == nil {
		t.Error("lulz")
	}
}

func TestPasswordLogin(t *testing.T) {
	configuration := createConfig()
	configuration["DatabaseUsername"] = "bob"
	configuration["DatabasePassword"] = "password"
	db, err := Connect(configuration)

	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	var name string
	row := db.QueryRow("select name from accounts")
	err = row.Scan(&name)
	if err != nil {
		t.Error(err)
	}
	if name != "myaccountname" {
		t.Error("Didn't get myaccountname as a result")
	}
}

func TestTeardown(t *testing.T) {
	// Teardown
	configuration := createConfig()
	dbnametodelete := configuration["Database"]
	delete(configuration, "Database")
	db, _ := Connect(configuration)

	_, err := db.Exec("drop database " + dbnametodelete)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec("drop user 'bob'@'localhost'")
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec("drop user 'bob'@'%'")
  if err != nil {
    t.Error(err)
  }

  _, err = db.Exec("flush privileges")
	if err != nil {
		t.Error(err)
	}
	db.Close()

}
