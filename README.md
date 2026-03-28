# Task Manager Go

**This project is a simple CRUD demonstration while exploring
various go concepts like;error handling, dependency injection in go,
interface/abstraction like in OOPs (satisfaction?), implementing JWTauth
via cookie and also authorization header.**

Basically, I made this project to have more deeper understanding about back-end
and Go language itself.

This project is based on RESTful API:

### END POINTS

**Auth End Points**

| METHOD | END POINTS | USE CASE |
| ------ | ---------- | -------  |
| POST   | /api/auth/register | To register user |
| POST   | /api/auth/login    | To login user    |
| POST   | /api/auth/refresh  | For generating refresh token |

**Tasks End Points**

| METHOD | END POINTS | USE CASE |
| ------ | --------   | -------- |
| POST   | /api/task/create | Create a task |
| PUT    | /api/task/update | Update a task |
| GET    | /api/task/by-task-id/{taskId} | Get task by it's id |
| GET    | /api/task/by-user-id/{userId} | Get tasks by user id |
| GET    | /api/task/all | Get all tasks |
| DELETE | /api/task/{taskId} | Delete the task |

### How to run?

**Setting .env file**

- Copy example.env to .env and fill the values accordingly

```
DB_USER=myuser
DB_PASSWORD=mypassword
DB_PORT=5432
DB_NAME=mydb
DB_HOST=localhost

# Generate using command (if linux): openssl rand -base64 32
JWT_KEY=YOUR_JWT_KEY
```

**Migration**

```bash
go run migration.go
```

**Running the project**

```bash
go run main.go
```

**Building the project**

```bash
go build -o taskmanager .
```
