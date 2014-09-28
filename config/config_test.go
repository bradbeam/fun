package config

import "testing"

func TestConfig(t *testing.T) {
  t.Log("Setting configuration")
  configuration := Config{
    "dbname": "l1g",
    "dbhost": "127.0.0.1",
    "dbusername": "root",
    "dbpassword": "password",
  }

  datatests(configuration, t)
}

func TestLoadConfig(t *testing.T) {
  configfile := "./test/config_test.cfg"
  t.Log("Loading config file " + configfile )
  configuration := make(Config)
  err := configuration.LoadConfigFromFile(configfile)
  if err != nil {
    t.Error(err)
  }
  datatests(configuration, t)
}

func TestFailFileLoad(t *testing.T) {
  configfile := "./test/config_test_that_doesnt_exist.cfg"
  t.Log("Loading config file " + configfile )
  configuration := make(Config)
  err := configuration.LoadConfigFromFile(configfile)
  if err == nil {
    t.Error("Loaded a configuration file that shouldn't exist. WTH")
    t.Log(err)
  }
}

func TestInvalidConfigItem(t *testing.T) {
  configfile := "./test/bad_config_test.cfg"
  t.Log("Loading config file " + configfile )
  configuration := make(Config)
  err := configuration.LoadConfigFromFile(configfile)
  if err == nil {
    t.Error("Successfully loaded an invalid configuration file.")
    t.Log(err)
  }
}

func datatests(configuration Config, t *testing.T) {
  if configuration["dbname"] != "l1g" {
    t.Error("config.dbname is not correct")
    t.Log("config.dbname = " + configuration["dbname"])
  }

  if configuration["dbhost"] != "127.0.0.1" {
    t.Error("config.dbhost is not correct")
    t.Log("config.dbhost = " + configuration["dbhost"])
  }

  if configuration["dbusername"] != "root" {
    t.Error("config.dbusername is not correct")
    t.Log("config.dbusername = " + configuration["dbusername"])
  }

  if configuration["dbpassword"] != "password" {
    t.Error("config.dbpassword is not correct")
    t.Log("config.dbpassword = " + configuration["dbpassword"])
  }
}
