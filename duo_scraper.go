package main

import (
  "fmt"
  // "net/http"
  "io/ioutil"
  "encoding/json"
  // "reflect"
)

type User struct {
  Name string
  UserName string
  Languages []LingoData
}

type LingoData struct {
  Language string
  Level int
  Points int
}

func (u *User) save() error {
  filename := "test.txt"
  return ioutil.WriteFile(filename, []byte("test"), 0600)
}

func (u *User) Unmarshal(data []byte) error {
  var f interface{}

  err := json.Unmarshal(data, &f)

  languages := f.(map[string]interface{})["languages"].([]interface{})

  var ld LingoData

  for i := range languages {
    lingo_map := languages[i].(map[string]interface{})
    ld.Language = lingo_map["language_string"].(string)
    ld.Level = int(lingo_map["level"].(float64))
    ld.Points = int(lingo_map["points"].(float64))
    u.Languages = append(u.Languages, ld)
  }

  return err
}

func SeedUsers() []User {
  users := []User{}

  users = append(users,
    User{ Name: "Kevin K", UserName: "KevinKelle6"},
    User{ Name: "Max W", UserName: "MaxWofford"},
    User{ Name: "Casey R", UserName: "nocashvaluedrumz"},
    User{ Name: "Rob C", UserName: "robcole42"},
    User{ Name: "Marc G", UserName: "ogoog"},
    User{ Name: "Wei L", UserName: "puffpuffpuff"},
    User{ Name: "Matt S", UserName: "Stringerbell1"},
    User{ Name: "Alexey K", UserName: "alexeymk"},
    User{ Name: "Nicola H", UserName: "nicolarhughes"},
    )

  return users
}

func main() {
  users := SeedUsers()

  for i := range users {
    fmt.Println(users[i])
  }

  /*

    fmt.Println("In main.")
  resp, _ := http.Get("https://www.duolingo.com/users/nocashvaluedrumz")
  defer resp.Body.Close()

  data, _ := ioutil.ReadAll(resp.Body)

  u.Unmarshal(data)
  */
}
