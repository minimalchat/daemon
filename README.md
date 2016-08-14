# Let's Chat

Tired of crap or expensive live chat platforms, this aims to be the API daemon service for an open source live chat implementation.

Lets keep it simple.

## We want to provide:

 - Websocket endpoint to pass messages from client browser to server
 - Transport messages from server to some chat operator service
  (irc, facebook, sms, ?)
 - Save messages and clients (tagging, classifying, pages they are visiting) in a manor that allows clients to be ephemeral.

## We can break these down into smaller more abstract bits:

### Front-end

A. Client browser script
  - Show user chat onscreen, without degrading user experience
  - Allow feedback of chat/operator

### Back-end

B. Client-server communication (Websockets)
  - Receive/Deliver events (messages)

C. Server database operations (Create/read/update/delete)
  - Keep record of chats for resumption
  - Keep record of clients and chats for analysis
  - Keep record of operator tagging

D. Server-operator communication
  - Interpret actions (Tag, messages)
  - Operation communication (IRC, facebook, sms, phone/mobile app, bot)


## Questions that need to be answered to move forward

How do websockets work on the Browser? Can we use socket.io?

How do websockets work with Go.

What database should be used? Is in-memory good enough?


## Steps to implement

1. Setup websocket server side (B)
2. Connect rudimentary client side to test server side (A)
3. Build out saving messages/chats/clients (C)
4. Create initial operator transport \[laptop\] (D)
5. Refresh rudimentary client side with production level implementation (A)
6. Beta test
7. Fix bugs
8. Build out secondary operator transport \[non-laptop\] (D)

### Bonus levels

9. Implement multi tenant implementation
10. Build out marketing and hosted solution
11. Create payment portal/admin Browser portal