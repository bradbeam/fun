package config

import (
  "bufio"
  "errors"
  "os"
  "strings"
)

type Config map[string]string

func (c *Config) LoadConfigFromFile(filename string) error {
  // open input file
  file, err := os.Open(filename)
  if err != nil {
      return err
  }

  // close file on exit and check for its returned error
  defer func() {
      if err := file.Close(); err != nil {
          panic(err)
      }
  }()

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
      input := strings.Split(scanner.Text(), "=")
      if len(input) != 2 {
        err := errors.New("Invalid configuration item: " + scanner.Text() )
        return err
      }
      (*c)[strings.Trim(input[0], " ")] = strings.Trim(input[1], " ")
  }

  return nil
}
