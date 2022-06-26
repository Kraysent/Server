# Example web server
This is example web server written in Go. It has REST API with the following handlers:

- `/ping` - healthcheck. Should respond `pong`.
- `/register` - create new user with given login, description and password.
- `/login` - check if given login and password are correct. 

It uses PostgreSQL as database management system. 