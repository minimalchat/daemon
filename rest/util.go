package rest

import (
  "fmt"
  "net/http"
  // "encoding/json"

  "github.com/julienschmidt/httprouter"

  // "github.com/minimalchat/mnml-daemon/store"
)


// Default Route

// GET /api/
func CatchAll (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
  resp.WriteHeader(http.StatusOK)
  fmt.Fprint(resp, "{\"clients\": \"/api/clients\", \"client\": \"/api/client/:id\", \"chats\":\"/api/chats\", \"chat\":\"/api/chat/:id\", \"messages\":\"/api/chat/:id/messages\", \"message\":\"/api/chat/:id/message/:mid\", \"operators\":\"/api/operators\", \"operators\":\"/api/operators\", \"operator\":\"/api/operator/:id\"}")
}

func NotImplemented (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  resp.Header().Set("Content-Type", "text/plain; charset=UTF-8")
  resp.WriteHeader(http.StatusNotImplemented)

  fmt.Fprintf(resp, "Not Implemented")
}
