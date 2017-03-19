package main

import (
  "log"
  "fmt"
  "time"
  "strings"

  "github.com/enodata/faker"
  "github.com/wayn3h0/go-uuid" // UUID (RFC 4122)

  "github.com/minimalchat/mnml-daemon/operator"
  "github.com/minimalchat/mnml-daemon/client"
  "github.com/minimalchat/mnml-daemon/chat"
  "github.com/minimalchat/mnml-daemon/person"
)

func generate_Chat() chat.Chat {
  cl := generate_Client()
  op := generate_Operator()
  creation := faker.Time().Between(time.Unix(1451606400,0), time.Now())

  return chat.Chat{
    ID: faker.Internet().Slug(),
    Client: &cl,
    Operator: &op,
    CreationTime: creation,
    UpdatedTime: faker.Time().Between(creation, time.Now()),
    Open: faker.Number().Number(5) > faker.Number().Number(3),
  }
}

func generate_Client() client.Client {
  p := generate_Person()
  uuid, _ := uuid.NewRandom()

  return client.Client{
    Person: p,
    Name: fmt.Sprintf("%s %s", p.FirstName, p.LastName),
    Uuid: uuid.String(),
  }
}

func generate_Operator() operator.Operator {
  p := generate_Person()

  return operator.Operator{
    Person: p,
    UserName: fmt.Sprintf("%s-%s", strings.ToLower(p.FirstName), strings.ToLower(p.LastName)),
  }
}

func generate_Person() person.Person {
  return person.Person{
    FirstName: faker.Name().FirstName(),
    LastName: faker.Name().LastName(),
  }
}


func main() {
  // log.Println("DEBUG", "generate", "operator", generate_Operator())
  // log.Println("DEBUG", "generate", "client", generate_Client())
  log.Println("DEBUG", "generate", "chat", generate_Chat())
}
