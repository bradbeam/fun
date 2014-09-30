package characters

type Character struct {
  Name string
  Level uint
  Health int
  Mana int
  Type string
  AttackPower int
  DefensePower int
  Status string
}

func Attack(attacker *Character, defender *Character) {
  damage := 0
  if ( defender.DefensePower >= attacker.AttackPower ) {
    damage = 1
  } else {
    damage = attacker.AttackPower - defender.DefensePower
  }

  defender.Health -= damage

  if ( defender.Health <= 0 ) {
    defender.Health = 0
    defender.Status = "dead"
  }
}

func Experience(character *Character, killedLevel int) {

}
