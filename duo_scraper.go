package main

import (
  "fmt"
  "net/http"
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

func (u *User) UnmarshalFromDL(data []byte) error {
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

func seedUsers(users *[]User) {
  *users = append(*users,
    User{ Name: "Kevin K", UserName: "KevinKelle6"},
    User{ Name: "Max W", UserName: "MaxWofford"},
    User{ Name: "Casey R", UserName: "nocashvaluedrumz"},
    User{ Name: "Rob C", UserName: "robcole42"},
    User{ Name: "Marc G", UserName: "ogoog"},
    User{ Name: "Wei L", UserName: "puffpuffpuff"},
    User{ Name: "Matt S", UserName: "Stringerbell1"},
    User{ Name: "Alexey K", UserName: "AlexeyKomi"},
    User{ Name: "Nicola H", UserName: "nicolarhughes"},
    User{ Name: "Luka K", UserName: "lukakacil"},
    User{ Name: "Charlie G", UserName: "charlierguo"},
    )
}

func scrapeLanguageData(users *[]User) {
  for i := range *users {
    user := (*users)[i]
    url := "https://www.duolingo.com/users/" + user.UserName
    
    resp, err := http.Get(url); if err != nil {
      fmt.Println(err)
    }
    defer resp.Body.Close()

    data, err := ioutil.ReadAll(resp.Body); if err != nil {
      fmt.Println(err)
    }

    err = user.UnmarshalFromDL(data); if err != nil {
      fmt.Println(err)
    }

    (*users)[i] = user
  }
}

func saveUserData(users []User) error {

  filename := "user_data.txt"

  data, err := json.Marshal(users)

  if err != nil {
    return err
  }

  return ioutil.WriteFile(filename, data, 0600)
}

func loadUserData() []byte {
  filename := "user_data.txt"
  data, _ := ioutil.ReadFile(filename)
  return data
}

func unmarshalSavedData(users *[]User, data []byte) {
  var f interface{}

  _ = json.Unmarshal(data, &f)

  data_arr := f.([]interface{})

  for i := range data_arr {
    var u User
    data_map := data_arr[i].(map[string]interface{})
    
    u.Name =  data_map["Name"].(string)
    u.UserName = data_map["UserName"].(string)
    
    language_data := data_map["Languages"].([]interface{})
    var ld LingoData

    for i := range language_data {
      lingo_map := language_data[i].(map[string]interface{})
      ld.Language = lingo_map["Language"].(string)
      ld.Level = int(lingo_map["Level"].(float64))
      ld.Points = int(lingo_map["Points"].(float64))
      u.Languages = append(u.Languages, ld)
    }
    *users = append(*users, u)
  }
}

func main() {
  users := []User{}

  data := loadUserData()

  unmarshalSavedData(&users, data)

  fmt.Println(users[4])

  // seedUsers(&users)

  // scrapeLanguageData(&users)

  // saveUserData(users)
}
