# 📝 Chat App (Golang)

## 📦 Installation
Make sure you have **nodemon** installed:
```sh
npm install -g nodemon
```

### Seed Database (Users & Rooms)
```sh
go run ./server seed
```

## 🚀 Running the Project

### Start the Server
```sh
nodemon --exec "go run ./server start" --watch server --ext go --signal SIGTERM
```

### Start the Client
```sh
nodemon --exec "go run ./client" --watch client --ext go --signal SIGTERM
```

### Connect to Chat as a Specific User
*(Enable `customizeState` in `chat/client/client.go` to use this feature)*
```sh
nodemon --exec "go run ./client --user Sandor" --watch client --ext go --signal SIGTERM
nodemon --exec "go run ./client --user Arya" --watch client --ext go --signal SIGTERM
```

### Start Logger (Server & Client)
```sh
nodemon --exec "go run ./logger" --watch logger --ext go --signal SIGTERM
```



---

## 🧑‍💻 Default Users
| Username        | Password  |
|---------------|------------|
| Sandor Clegane | `Test1234` |
| Arya Stark     | `Test1234` |

---


