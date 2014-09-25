package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "html/template"
  // "reflect"
)

type User struct {
  Name, UserName string
  // ex => Languages["Spanish"]["Points"]
  Languages map[string]map[string]int
}

type Output struct {
  Users []User
}

func (u *User) UnmarshalFromDL(data []byte) error {
  var f interface{}

  err := json.Unmarshal(data, &f)

  languages := f.(map[string]interface{})["languages"].([]interface{})

  for i := range languages {
    lingo_map := languages[i].(map[string]interface{})
    language := lingo_map["language_string"].(string)

    m := make(map[string]int)

    m["Level"] = int(lingo_map["level"].(float64))
    m["Points"] = int(lingo_map["points"].(float64))

    if (*u).Languages[language] == nil {
      language_map := make(map[string]map[string]int)
      language_map[language] = m
      (*u).Languages = language_map
    } else {
      (*u).Languages[language] = m
    }
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
    
    lingo_map := data_map["Languages"].(map[string]map[string]int)
    
    for language, sub_map := range lingo_map {
      u.Languages[language]["Level"] = sub_map["Level"]
      u.Languages[language]["Points"] = sub_map["Points"]
    }
    *users = append(*users, u)
  }
}

func leaderBoardHandler(w http.ResponseWriter, r *http.Request) {
  title := "DuoLingo Leaderboard:"
  t, _ := template.ParseFiles("leaderboard.html")
  // var u []User
  t.Execute(w, title)
}

func main() {
  users := []User{}

  // data := loadUserData()

  //unmarshalSavedData(&users, data)

  /*
  http.HandleFunc("/leaderboard/", leaderBoardHandler)
  http.ListenAndServe(":8080", nil)
  */

  seedUsers(&users)

  scrapeLanguageData(&users)

  saveUserData(users)

  fmt.Println(users)

}
