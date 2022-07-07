# Example web server
This is example web server written in Go. It has REST API with the following handlers:

- `/ping` - healthcheck. Should respond `pong`.
- `/register` - create new user with given login, description and password. Creates cookie.
- `/login` - check if given login and password are correct. Creates cookie.
- `/profile` - get info about the user. Requires auth cookie.
- `/search` - get info about the users with login according to the pattern. Requires auth cookie.

It uses PostgreSQL as database management system. 

There is also an example of Github Action that runs `go test`.