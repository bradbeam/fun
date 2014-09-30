package info

import (
  "log"
  "net/http"
  "encoding/json"
  "github.com/bradbeam/fun/client"
  // "github.com/bradbeam/fun/config"
  "github.com/bradbeam/fun/characters"
  // "github.com/bradbeam/fun/fundb"
  "database/sql"
  // "strconv"
  "expvar"
)

type AccountInfo struct {
  Username string `json:Username`
  CharacterList []characters.Character `json:Characters`
}

var InfoRequests = expvar.NewInt("InfoRequests")

func Account(rw http.ResponseWriter, req *http.Request, mydb *sql.DB) {
  InfoRequests.Add(1)

  decoder := json.NewDecoder(req.Body)
  var ir client.ClientInfoRequest
  err := decoder.Decode(&ir)
  if err != nil {
      panic(err)
  }
  log.Println("InfoRequest:")
  log.Println(ir)

  // ai := AccountInfo{}

  myChars, err := loadData(ir.Account, mydb)
  if err != nil {
    log.Println(err)
    http.NotFound(rw, req)
    return
  }

  log.Println(myChars)

  // return results
  response := client.ClientInfoResponse{}
  response = client.ClientInfoResponse{ir.Account, myChars}

  results, err := json.Marshal(response)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    log.Println(err.Error())
    return
  }

  rw.Header().Set("Content-Type", "application/json")
  rw.Write(results)

}

func loadData(account string, mydb *sql.DB) ([]characters.Character, error) {
  chars := []characters.Character{}

  var err error

  log.Println("Loading account " + account)
  stmt, err := mydb.Prepare("SELECT name,level,health,mana,attackpower,defensepower,status " +
                              "FROM fun.characters " +
                              "JOIN fun.account ON fun.characters.account_id = fun.account.id " +
                              "WHERE fun.account.username = ?;")
  if err != nil {
    return chars, err
  }
  defer stmt.Close()

  rows, err := stmt.Query(account)
  for rows.Next() {
    var char characters.Character

    err := rows.Scan(&char.Name, &char.Level, &char.Health, &char.Mana, &char.AttackPower, &char.DefensePower, &char.Status);
    if err != nil {
      log.Fatal(err)
    }
    chars = append(chars, char)
  }


  return chars, err
}
