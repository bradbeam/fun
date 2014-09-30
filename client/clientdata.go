package client

import "github.com/bradbeam/fun/characters"

type ClientActionRequest struct {
  Account string `json:Account`
  Character string `json:Character`
  Type string `json:Type`
}

type ClientAttackResponse struct {
  Result string "json:Result"
}

type ClientInfoRequest struct {
  Account string `json:Account`
  Character string `json:Character`
  Type string `json:Type`
}

type ClientInfoResponse struct {
  Account string `json:Account`
  Character []characters.Character `json:Characters`
}
