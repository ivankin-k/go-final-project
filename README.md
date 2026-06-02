# Go Practicum Final Project (ToDo Web App)
Nice and small todo task app with following features:
- add/delete/edit tasks
- set repeat rules (daily, weekly, monthly, yearly)
- mark as done
- search tasks by date or text (in title or comment)
- password-based authentication (optionally enabled using `TODO_PASSWORD` env variable)

#### All ⭐ tasks are completed:
- `TODO_PORT` env
- `TODO_DBFILE` env
- `w` and `m` repeat rules
- search
- authentication (`TODO_PASSWORD`)
- Dockerfile

## 📁 Clone repository
```
git clone git@github.com:ivankin-k/go-final-project.git && cd go-final-project
```

## 🧪 Run & Test
### Test settings
Adjust test settings (if needed) in file `tests/settings.go`
```
var Port = 7540                 # default port
var DBFile = "../scheduler.db"  # default DB file name
var FullNextDate = true         # run tests for all repeat rules (d, w, m, y)
var Search = true               # run search tests
var Token = ...                 # token to test authentication
```
> [!NOTE]
> Token is currently set for password=`pass123`, so use default run command below to test.

### Start the app
```
TODO_PASSWORD=pass123 go run .
```
Go to [http://localhost:7540/](http://localhost:7540/)
### Run tests
```
go test -v ./tests
```

## 🐳 Docker

### Build
```
docker build --rm -t go-final-project:v1.0 .
```

### Run using defaults
```
docker run --rm -p 7541:7540 go-final-project:v1.0
```
- port: `7540`
- password: <none> (authentication disabled)
- DB filename: `scheduler.db` (already contains dummy tasks to play with)

Go to [http://localhost:7541/](http://localhost:7541/)

### Run using custom settings
```
docker run --rm -p <host_port>:<app_port> -e TODO_PORT=<app_port> -e TODO_PASSWORD=<password> -e TODO_DBFILE=<filename> go-final-project:v1.0
```
- port: `TODO_PORT=<app_port>`
- password: `TODO_PASSWORD=<password>`
- DB filename: `TODO_DBFILE=<filename>`