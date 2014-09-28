package main

import (
  "log"
  "github.com/bradbeam/fun/config"
  "github.com/bradbeam/fun/server"
)

func main() {
  // Will most likely want to add a flag for the config file location
  configfile := "./config/test/config_test.cfg"
  log.Println("Loading config file " + configfile )
  configuration := make(config.Config)
  err := configuration.LoadConfigFromFile(configfile)
  if err != nil {
    log.Fatal(err)
  }
  
  log.Println("Starting Server")
  server.Serve(configuration)
}
