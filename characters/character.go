package characters

type Character struct {
  Level uint
  Health int
  Mana int
  Type string
  AttackPower int
  DefensePower int
  Status string
}

func (c *Character) Defend(attacker *Character) {
  damage := 0
  if ( c.DefensePower >= attacker.AttackPower ) {
    damage = 1
  } else {
    damage = attacker.AttackPower - c.DefensePower
  }

  c.Health -= damage

  if ( c.Health <= 0 ) {
    c.Health = 0
    c.Status = "dead"
  }
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
