package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  // "reflect"
)

type User struct {
  languages string
}

func (u *User) save() error {
  filename := "test.txt"
  return ioutil.WriteFile(filename, []byte(u.languages), 0600)
}

func main() {
  resp, err := http.Get("https://www.duolingo.com/users/nocashvaluedrumz")
  defer resp.Body.Close()

  data, err := ioutil.ReadAll(resp.Body)

  fmt.Println("data:")
  // fmt.Println(string(data))

  var f interface{}

  _ = json.Unmarshal(data, &f)

  m := f.(map[string]interface{})

  var points map[string]int

  languages := m["languages"].([]interface{})

  for i := range languages {
    lingo_map := languages[i].(map[string]interface{})
    ls := lingo_map["language_string"]
    pts := lingo_map["points"]
    lvl := lingo_map["level"]

    // fmt.Printf("Lingo map is: %s\n", lingo_map)
    // fmt.Printf("Type is: %s\n", reflect.TypeOf(lingo_map))
    fmt.Printf("Language string is: %s\n", ls)
    fmt.Printf("Level is: %v\n", lvl)
    fmt.Printf("Total points: %v\n\n", pts)
    // levels[lingo_map["language_string"]] = lingo_map["level"]
  }

  for k, v := range points {
    fmt.Println(k + ": " + string(v))
  }

  if err != nil {
    panic(err)
  }
}
