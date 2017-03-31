package socket

import (
  "time"
  "fmt"
  "encoding/json"
  "bytes"
  "log"
  "net/http"

  "github.com/googollee/go-socket.io" // Socket

  "github.com/minimalchat/mnml-daemon/chat"
  "github.com/minimalchat/mnml-daemon/operator"
  "github.com/minimalchat/mnml-daemon/client"
  "github.com/minimalchat/mnml-daemon/store" // InMemory store
)

// Log levels
const (
  DEBUG string = "DEBUG"
  INFO string = "INFO"
  WARNING string = "WARN"
  ERROR string = "ERROR"
  FATAL string = "FATAL"
)

type SocketListener struct {
  Operators map[string]*operator.Operator
  Clients map[string]*client.Client
  Chats map[string]*chat.Chat
  Server *socketio.Server
}

func Listen(ds *store.InMemory) *SocketListener {
  log.Println(DEBUG, "socket:", "Listening for WebSocket clients ...")

  srv, err := socketio.NewServer(nil)
  sck := SocketListener{
    Operators: make(map[string]*operator.Operator),
    Clients: make(map[string]*client.Client),
    Chats: make(map[string]*chat.Chat),
    Server: srv,
  }

  // TODO: Return an error instead
  if err != nil {
      log.Fatal(err)
  }

  srv.On("connection", sck.onConnection(ds))
  srv.On("error", sck.onError)

  return &sck
}

func (this SocketListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  this.Server.ServeHTTP(w, r);
}


func (this SocketListener) emitToOperators(event string, data string) {
  // if (event == nil) {
  //   log.Println(WARNING, "Unknown event to emit")
  //   return
  // }

  // Update Operators of the new messages
  for _, op := range this.Operators {
    log.Println(DEBUG, "socket:", fmt.Sprintf(" Sending %s \"%s\" to %s", event, data, op.Socket.Id()))

    op.Socket.Emit(event, data, nil)
  }
}


func (this SocketListener) onConnection(ds *store.InMemory) func (sock socketio.Socket) {
  return func (sock socketio.Socket) {
    log.Println(DEBUG, "socket:", fmt.Sprintf("Incoming connection %s %s", sock.Id(), sock.Request().URL.Query().Get("type")))

    t := sock.Request().URL.Query().Get("type")

    // TODO: Verify that the socket connection is real
    if (t == "operator") {

      // Create Operator

      // TODO: Pull Operator from DB if we already know them
      this.Operators[sock.Id()] = operator.Create(operator.Operator{
          // FirstName: "Operator",
          // LastName: "Steve",
          UserName: "steve",
        }, sock)

      // Save Operator to DB
      ds.Put(*this.Operators[sock.Id()])
    } else if (t == "client") {

      // Create Client

      // TODO: See if we can "recall" if this is a returning client?
      this.Clients[sock.Id()] = client.Create(client.Client{
          Name: "Site Visitor",
        }, sock)

      // Save Client to DB
      ds.Put(*this.Clients[sock.Id()])

      // Create Chat

      // TODO: See if we can "recall" the returning chat?
      this.Chats[sock.Id()] = chat.Create(chat.Chat{
        Client: this.Clients[sock.Id()],
        Operator: nil,
        Open: true,
        CreationTime: time.Now(),
        UpdatedTime: time.Now(),
      })

      // Save Chat to DB
      ds.Put(*this.Chats[sock.Id()]);

      jsonChat, _ := json.Marshal(this.Chats[sock.Id()])
      var buffer bytes.Buffer
      buffer.Write(jsonChat)
      buffer.WriteString("\n")

      // Emit to Operators
      this.emitToOperators("chat:new", buffer.String())
    } else {

      // TODO: Write some proper error handling here, do we close the connection?
      log.Println(ERROR, "socket:", "Unknown chat type specified")
    }


    sock.On("client:message", this.onClientMessage(ds, sock))
    sock.On("operator:message", this.onOperatorMessage(ds, sock))

    // Disconnection event
    sock.On("disconnection", func () {
      log.Println(DEBUG, "socket:", fmt.Sprintf("%s disconnected", sock.Id()))

      // TODO: Save chat?

      if (t == "operator") {

        delete(this.Operators, sock.Id())
      } else if (t == "client") {

        delete(this.Clients, sock.Id())
      }
    })
  }
}


func (this SocketListener) onClientMessage(ds *store.InMemory, sock socketio.Socket) func (msg string) {
  return func (msg string) {

    log.Println(DEBUG, "client", fmt.Sprintf("%s: %s", sock.Id(), msg))

    // Create Message
    m := chat.Message{
      Timestamp: time.Now(),
      Content: msg,
      Author: this.Clients[sock.Id()].StoreKey(),
      Chat: this.Chats[sock.Id()].Uid,
    }

    // Save Message to DB
    ds.Put(m)

    // Update Operators of the new messages
    this.emitToOperators("client:message", msg)
  }
}


func (this SocketListener) onOperatorMessage(ds *store.InMemory, sock socketio.Socket) func (msg string) {
  return func (msg string) {
    // TODO: Get chat from message, and then update it/send to correct client

    log.Println(DEBUG, "operator", fmt.Sprintf("%s: %s", sock.Id(), msg))

    // Create Message
    // m := chat.Message{
    //   Timestamp: time.Now(),
    //   Content: msg,
    //   Author: this.Operators[sock.Id()].StoreKey(),
    //   Chat: ch.ID,
    // }

    // Save Message to DB
    // ds.Put(m)
  }
}


func (this SocketListener) onError(sock socketio.Socket, err error) {

  // TODO: Write some proper error handling here
  log.Println(ERROR, "socket:", err)
}