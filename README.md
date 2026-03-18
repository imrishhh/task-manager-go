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
