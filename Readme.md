# CASE FOR SOFTWARE ENGINEERING SCHOOL 4.0

## Project Structure

The project is structured as follows:


```
├── cmd
│   └── app
│       └── main.go                 # Main application code
├── internal
│   ├── db
│   │   ├── db.go                   # Database related code
│   │   └── migration.go            # Database migration code
│   ├── email
│   │   └── email.go                # Email related code
│   ├── handlers
│   │   ├── rate.go                 # Handler for rate
│   │   └── subscribe.go            # Handler for subscribe
│   ├── scheduler
│   │   └── scheduler.go            # Scheduler related code
│   └── tests
│       ├── db_test.go              # Test for db
│       ├── email_test.go           # Test for email
│       ├── rate_test.go            # Test for rate handler
│       ├── scheduler_test.go       # Test for scheduler
│       └── subscribe_test.go       # Test for subscribe handler
├── migrations
│   ├── 0001_initial.down.sql       # Initial migration down
│   └── 0001_initial.up.sql         # Initial migration up
├── go.mod                          # Go module file
├── http
│   ├── rate.http                   # HTTP requests for rate API
│   └── subscribe.http              # HTTP requests for subscribe API
```

## Important Note

In this project, an SMTP emulated server is used to test the email-sending functionality. It does not send emails to actual email addresses; it only emulates the email-sending process. The emulated server used in this project is [Ethereal Email](https://ethereal.email/).

This is a test case for a school project. For more information, visit [SOFTWARE ENGINEERING SCHOOL 4.0](https://www.genesis-for-univ.com/registration-software-engineering-school-4).

## Requirements

- Docker
- Docker Compose
- Go (if running locally)

## Getting Started

### Docker

1. **Build the Docker image:**

    ```bash
    docker-compose build
    ```

2. **Run the application:**

    ```bash
    docker-compose up
    ```

### Locally

1. **Install Go:**

   Follow the instructions on the [official Go website](https://golang.org/doc/install) to install Go.

2. **Install gcc:**

   - **If using Linux:**

       ```bash
       sudo apt-get install gcc
       ```

   - **If using Windows:** Follow the instructions [here](https://code.visualstudio.com/docs/cpp/config-mingw#_installing-the-mingww64-toolchain).

3. **Run the application:**

    ```bash
    go run cmd/app/main.go
    ```

4. **Run the tests:**

    ```bash
    go test ./...
    ```

## Tests

Tests are included in the `internal/tests` directory. They cover the key functionalities of the application.

