package actions

import (
  "log"
  "net/http"
  "encoding/json"
  "github.com/bradbeam/fun/client"
  // "github.com/bradbeam/fun/config"
  "github.com/bradbeam/fun/characters"
  // "github.com/bradbeam/fun/fundb"
  "database/sql"
  "strconv"
)

type AttackRequest struct {
  Source client.ClientData
  Target client.ClientData
}

type AttackResponse struct {
  Result string "json:Result"
}

func Attack(rw http.ResponseWriter, req *http.Request, mydb *sql.DB) {
    decoder := json.NewDecoder(req.Body)
    var ar AttackRequest
    err := decoder.Decode(&ar)
    if err != nil {
        panic(err)
    }
    log.Println(ar)

    // Load up Account+Character
    sourceCharacter, err := loadData(ar.Source, mydb)
    if err != nil {
      log.Println(err)
      http.NotFound(rw, req)
      return
    }
    log.Println(sourceCharacter)

    targetCharacter, err := loadData(ar.Target, mydb)
    if err != nil {
      log.Println(err)
      http.NotFound(rw, req)
      return
    }
    log.Println(targetCharacter)

    // Simulate attack
    i := 1;
    for ( sourceCharacter.Status != "dead" && targetCharacter.Status != "dead" ) {
      log.Println("Turn " + strconv.Itoa(i) + " Attack")
      characters.Attack(&sourceCharacter, &targetCharacter)
      log.Println(sourceCharacter)
      log.Println(targetCharacter)
      log.Println("Turn " + strconv.Itoa(i) + " Defend")
      characters.Attack(&targetCharacter, &sourceCharacter)
      log.Println(sourceCharacter)
      log.Println(targetCharacter)
      i++
    }

    // Update DB
    if ( ar.Source.Type == "player" ) {
      log.Println("Saving player data " + ar.Source.Character)
      saveData(ar.Source, mydb)

    }
    if ( ar.Target.Type == "player" ) {
      log.Println("Saving player data " + ar.Target.Character)
      saveData(ar.Target, mydb)

    }

    // return results
    if ( ar.Source.Type == "player" ) {
        response := AttackResponse{}
        if ( targetCharacter.Status == "dead" ) {
          response = AttackResponse{"You won!"}
        } else {
          response = AttackResponse{"You lost!"}
        }
        results, err := json.Marshal(response)
        if err != nil {
          http.Error(rw, err.Error(), http.StatusInternalServerError)
          log.Println(err.Error())
          return
        }

        rw.Header().Set("Content-Type", "application/json")
        rw.Write(results)
    }
}


func loadData(character client.ClientData, mydb *sql.DB) (characters.Character, error) {
  char := characters.Character{}

  var err error
  if ( character.Type == "monster" ) {
    log.Println("Loading monster " + character.Character)
    stmt, err := mydb.Prepare("SELECT level,health,mana,attackpower,defensepower,status " +
                              "FROM fun.monsters " +
                              "WHERE name = ?");
    if err != nil {
      return char, err
    }
    defer stmt.Close()
    err = stmt.QueryRow(character.Character).Scan(&char.Level, &char.Health, &char.Mana, &char.AttackPower, &char.DefensePower, &char.Status)
  } else {
    log.Println("Loading player " + character.Character)
    stmt, err := mydb.Prepare("SELECT level,health,mana,attackpower,defensepower,status " +
                                "FROM fun.characters " +
                                "JOIN fun.account ON fun.characters.account_id = fun.account.id " +
                                "WHERE fun.account.username  = ? AND fun.characters.name = ?;")
    if err != nil {
      return char, err
    }
    defer stmt.Close()
    err = stmt.QueryRow(character.Account, character.Character).Scan(&char.Level, &char.Health, &char.Mana, &char.AttackPower, &char.DefensePower, &char.Status)
  }


  return char, err
}

func saveData(character client.ClientData, mydb *sql.DB) (error) {

  log.Println("Saving player " + character.Character)
  stmt, err := mydb.Prepare("UPDATE fun.characters JOIN fun.account ON fun.characters.account_id = fun.account.id " +
                              "SET kills=kills+1 " +
                              "WHERE fun.account.username  = ? AND fun.characters.name = ?;")
  if err != nil {
    return err
  }
  defer stmt.Close()

  _,err = stmt.Exec(character.Account, character.Character)

  return err
}
