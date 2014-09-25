package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "html/template"
  "time"
  "os"
  "io"
  "sort"
  // "reflect"
)

type User struct {
  Name, UserName string
  DailyDiff int
  // ex => Languages["Spanish"]["Points"]
  Languages map[string]map[string]int
}

type scrapeHistory struct {
  Date string
}

type Output struct {
  Users []User
}

// define methods necessary for sorting output
func (o Output) Len() int { return len(o.Users) }
func (o Output) Swap(i, j int) { o.Users[i], o.Users[j] = o.Users[j], o.Users[i] }
func (o Output) Less(i, j int) bool { return o.Users[i].Level() < o.Users[j].Level() }

func (u User) Level() int {
  return u.Languages["Spanish"]["Level"]
}

func (u User) Points() int {
  return u.Languages["Spanish"]["Points"]
}

func (s *scrapeHistory) Save() error {
  filename := "scrape_history.txt"
  // open file to which we append the latest scrape history
  f, _ := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0660);
  output, err := json.Marshal(*s)
  io.WriteString(f, string(output))
  return err
}

func loadScrapeHistory() []scrapeHistory {
  // read in file
  filename := "scrape_history.txt"
  data, _ := ioutil.ReadFile(filename)
  // convert JSON to an interface which we can process
  var f interface{}
  if err := json.Unmarshal(data, &f); err != nil {
    panic(err)
  }
  data_arr := f.([]interface{})
  // fill an array of scrapeHisotry from the data
  scrapeHistories := []scrapeHistory{}
  for i := range data_arr {
    sH := data_arr[i].(scrapeHistory)
    fmt.Println(sH)
    scrapeHistories = append(scrapeHistories, sH)
  }
  return scrapeHistories
}

func lastScrapeHistory() scrapeHistory {
  sh := loadScrapeHistory()
  last := sh[len(sh)]
  return last
}

func (u *User) UnmarshalFromDL(data []byte) error {
  var f interface{}

  err := json.Unmarshal(data, &f)

  // access "languages" part of JSON map through type assertions
  languages := f.(map[string]interface{})["languages"].([]interface{})

  // iterate through languages and add to the user
  for i := range languages {
    // type assert the langauge and get data
    lingo_map := languages[i].(map[string]interface{})
    language := lingo_map["language_string"].(string)

    // intialize map to store JSON data
    m := make(map[string]int)

    m["Level"] = int(lingo_map["level"].(float64))
    m["Points"] = int(lingo_map["points"].(float64))

    // initialize diff variable for calculating point differential
    var diff int

    // check to see if user is already initialized and then add new data
    if (*u).Languages[language] == nil {
      language_map := make(map[string]map[string]int)
      language_map[language] = m
      (*u).Languages = language_map
    } else {
      // assign daily point differential to User
      if language == "Spanish" {
        // find difference between today's and yesterday's points
        priorPoints := (*u).Languages[language]["Points"]
        diff = m["Points"] - priorPoints
      }

      last_scrape := lastScrapeHistory()
      scrape_date, _ := time.Parse(time.RFC1123, last_scrape.Date)

      if time.Since(scrape_date) >= 24*time.Hour {
        (*u).DailyDiff = diff
        (*u).Languages[language] = m
      }
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

    // unmarshal scraped data from DuoLingo into a user
    err = user.UnmarshalFromDL(data); if err != nil {
      fmt.Println(err)
    }

    (*users)[i] = user
  }

  // create a log entry of this scrape
  var s scrapeHistory
  s.Date = time.Now().Format(time.RFC1123)
  s.Save()
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

  // unmarshal data into raw JSON
  _ = json.Unmarshal(data, &f)

  data_arr := f.([]interface{})

  // iterate over JSON data
  for i := range data_arr {
    var u User
    data_map := data_arr[i].(map[string]interface{})
    
    // gather simple primitive data from JSON
    u.Name =  data_map["Name"].(string)
    u.UserName = data_map["UserName"].(string)
    u.DailyDiff = int(data_map["DailyDiff"].(float64))
    u.Languages = make(map[string]map[string]int)
    
    // iterate over languages map from the user and store in the struct
    lingo_map := data_map["Languages"].(map[string]interface{})
    for language, subset := range lingo_map {
      sub_map := subset.(map[string]interface{})
      u.Languages[language] = make(map[string]int)
      u.Languages[language]["Level"] = int(sub_map["Level"].(float64))
      u.Languages[language]["Points"] = int(sub_map["Points"].(float64))
    }
    // append the new user to the users slice
    *users = append(*users, u)
  }
}

func leaderBoardHandler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("leaderboard.html")
  var o Output
  o.Users = loadUsers()
  sort.Sort(sort.Reverse(o))
  t.Execute(w, o)
}

func loadUsers() []User {
  users := []User{}

  data := loadUserData()

  unmarshalSavedData(&users, data)

  return users
}

func onStart() []User {
  // check if scrape_history.txt already exists, if not create it
  if _, err := os.Stat("scrape_history.txt"); os.IsNotExist(err) {
    os.Create("scrape_history.txt")
  }
  // check if user_data.txt already exists, if not create it
  // and make initial scrape
  var users []User
  if _, err := os.Stat("user_data.txt"); os.IsNotExist(err) {
    os.Create("user_data.txt")
    seedUsers(&users)
    scrapeLanguageData(&users)
    saveUserData(users)
  } else {
    users = loadUsers()
  }
  return users
}

func main() {
  _ = onStart()
  
  http.HandleFunc("/leaderboard/", leaderBoardHandler)
  http.ListenAndServe(":8080", nil)

  // fmt.Println(users)
}
