## IMPORTANT NOTE
In project i have used smtp emulated server to test the email sending functionality. Actually it not sending email to the actual email address. It just emulating the email sending functionality.
https://ethereal.email/ is the emulated server i have used in this project.

## Project Structure

The project is structured as follows:

```
├── Dockerfile # Docker configuration
├── docker-compose.yml # Docker Compose configuration
├── main.go # Main application code
└── main_test.go # Tests for the application
```


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
    go run main.go
    ```

4. **Run the tests:**

    ```bash
    go test
    ```

## Tests

Tests are included in the `main_test.go` file. They cover the key functionalities of the application.
