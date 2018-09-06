# Websocket Event Documentation

---

## Connection

Simply connect to the websocket by follwoing URL:
```
ws://IPADDRESS:PORT/ws
```

---

## Data Types

### `Author`

Value | Type | Description
------|------|-------------
`id` | `int64` | Unique ID of the author
`username` | `string` | Username of the user
`color` | `string` | Chat color string of the user

### `Message`

Value | Type | Description
------|------|-------------
`id` | `int64` | Unique ID of the message
`author` | `Author` | Author of the message
`content` | `string` | Content string of the message
`timestamp` | `int64` | UnixNano timestamp of the message

### `MessageEvent`

Value | Type | Description
------|------|-------------
`event` | `string` | "message"
`data` | `Message` | Message object of the message event

---

## Receiving Events

Ingoing events needs to be send as JSON objects like following example scheme:
```json
{
    "event": "login",
    "data": {
        "username": "zekro",
        "password": "thisisabadpassword"
    }
}
```

### `login`

> Create Account session by username with password. If the username was never used before, the entered username will be protected by the entered password.

Parameter | Type | Description
----------|------|-------------
`data.username` | `string` | Username of the user
`data.password` | `string` | Password used for the username

**Emits:**
- `clientConnect`
- `connected`
- `connect_rejected`

### `checkUsername`

> Checks if the passed username is registered with a password.

Parameter | Type | Description
----------|------|-------------
`data` | `string` | Username of the user

**Emits:**
- `usernameState`

### `message`

> Send chat message to the socket.

Parameter | Type | Description
----------|------|-------------
`data` | `string` | Message content string

**Emits:**
- `spamTimeout`
- `message`

## `deleteMessage`

> Delete messsage by ID.

Parameter | Type | Description
----------|------|-------------
`data.msgid` | `int64` | ID of the message

**Emits:**
- `messageDeleted`

---

# Sending Events

Send events are also in JSON format in following example scheme:
```json
{
    "event": "clientConnect",
    "data": {
        "author": {
            "id": 123123123,
            "username": "zekro",
            "color": "#f4a3dc"
        },
        "nclients": 12,
        "clients": { ... },
        "history": [ ... ]
    }
}
```

### `clientConnect`

> Connection confirmation from socket after login.

Value | Type | Description
----------|------|-------------
`data.author` | `Author` | Author object of the logged in user
`data.nclients` | `int` | Number of currently connected sockets
`data.clients` | `map[string:string]` | Map of username:color of all connected users
`data.history` | `array[MessageEvent]` | Array of last send messages

### `connected`

> Broadcast message that a user connected to the chat room.

Value | Type | Description
----------|------|-------------
`data.author` | `Author` | Author object of the logged in user
`data.nclients` | `int` | Number of currently connected sockets
`data.clients` | `map[string:string]` | Map of username:color of all connected users

### `connect_rejected`

> Response that login request was rejected because of wrong username-password combination

*No Response Data.*

### `usernameState`

> Response of usernameCheck.

Value | Type | Description
----------|------|-------------
`data` | `boolean` | True if username is still registered with a password.

### `connect_rejected`

> Information response that the send message was rejected because of spam protection.

*No Response Data.*

### `message`

> Broadcast event if message was send to socket server.

Value | Type | Description
----------|------|-------------
`data.author` | `Author` | Author object of the logged in user
`data.content` | `string` | Message content text
`data.timestamp` | `int64` | UnixNano timestamp of the message
`data.id` | `int64` | Unique ID of the message

### `messageDeleted`

> Broadcast message that message was deleted on socket

Value | Type | Description
----------|------|-------------
`data` | `MessageEvent` | MessageEvent of the deleted message

### `disconnected`

> Broadcast message that client disconnected from socket.

Value | Type | Description
----------|------|-------------
`data.name` | `string` | Username of the disconnected user
`data.clients` | `map[string:string]` | Map of username:color of all connected users

---