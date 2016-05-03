package main
import (
  "github.com/gin-gonic/gin"
  "net/http"
  "fmt"
  "encoding/json"
  "bytes"
  "crypto/md5"
  "encoding/hex"
  "io/ioutil"
)


type SessionResponse struct {
    Session   string      `json:"session"`
}

type LoginResponse struct {
    Result   string          `json:"result"`
    Session   string         `json:"session"`
}

func debug_json(resp *http.Response) {

  json, _ := ioutil.ReadAll((resp.Body))
  fmt.Println(string(json))

}

func main() {
  host := ""
  user := ""
  password := ""
  sessionPayload := map[string]string{"cmd": "login"}
  loginUrl := fmt.Sprintf("http://%v:81/json", host)

  r := gin.Default()

  r.GET("/status",  func(c *gin.Context) {
    c.String(http.StatusOK, "")
  })

  r.GET("/", func(c *gin.Context) {
    payload, _ := json.Marshal(sessionPayload)
    req, err := http.NewRequest("POST", loginUrl, bytes.NewBufferString(string(payload)))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    session := SessionResponse{}
    json.NewDecoder(resp.Body).Decode(&session)
    hasher := md5.New()
    hasher.Write([]byte(fmt.Sprintf("%v:%v:%v", user, session.Session, password)))
    hexString := hex.EncodeToString(hasher.Sum(nil))

    loginPayload := map[string]string{ "cmd": "login", "session": session.Session, "response": hexString  }
    payload2, _ := json.Marshal(loginPayload)
    req2, err2 := http.NewRequest("POST", loginUrl, bytes.NewBufferString(string(payload2)))
    req2.Header.Set("Content-Type", "application/json")

    client2 := &http.Client{}
    resp2, err2 := client2.Do(req2)
    if err2 != nil {
        panic(err)
    }
    defer resp2.Body.Close()

    session2 := LoginResponse{}
    json.NewDecoder(resp2.Body).Decode(&session2)
    fmt.Println(session2)

    camListPayload := map[string]string{ "cmd": "camlist", "session": session2.Session }
  })

  r.Run(":8080")
}
