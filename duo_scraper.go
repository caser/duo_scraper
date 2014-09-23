package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  // "reflect"
)

type User struct {
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

func (u *User) Unmarshal(data []byte) {
  var f interface{}

  err := json.Unmarshal(data, &f)

  if err != nil {
    fmt.Println(err)
  }

  languages := f.(map[string]interface{})["languages"].([]interface{})

  var ld LingoData

  for i := range languages {
    lingo_map := languages[i].(map[string]interface{})
    ld.Language = lingo_map["language_string"].(string)
    ld.Level = int(lingo_map["level"].(float64))
    ld.Points = int(lingo_map["points"].(float64))
    u.Languages = append(u.Languages, ld)
  }
}

func main() {
  users := map[string]string
  users["kevin"] = "KevinKelle6"
  users["max"] = "MaxWofford"
  users["casey"] = "nocashvaluedrumz"
  users["rob"] = "robcole42"

  fmt.Println("In main.")
  resp, _ := http.Get("https://www.duolingo.com/users/nocashvaluedrumz")
  defer resp.Body.Close()

  data, _ := ioutil.ReadAll(resp.Body)

  u := User{UserName: "rob_cole"}

  u.Unmarshal(data)
}
