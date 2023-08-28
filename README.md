# Todo Application

This is a simple Todo application built in Go that provides a REST API to manage tasks. The application supports creating, listing, marking tasks as complete, and deleting tasks. It also includes a feature flag to enable or disable certain endpoints.

## Table of Contents

-   [Todo Application](#todo-application)
    -   [Table of Contents](#table-of-contents)
    -   [OpenAPI Specification](#openapi-specification)
    -   [Building the Application](#building-the-application)
    -   [Running the Application](#running-the-application)
    -   [Testing the Application](#testing-the-application)
    -   [Feature Flag](#feature-flag)
    -   [Dockerizing the Application](#dockerizing-the-application)
        -   [Building the Docker Image](#building-the-docker-image)
        -   [Running the Docker Container](#running-the-docker-container)

## OpenAPI Specification

The OpenAPI specification for the application's REST API is provided in the `openapi.yaml` file. This file describes the available endpoints, request and response structures, and their attributes. You can use tools like Swagger UI to visualize and interact with the API documentation. The OpenAPI spec file is located in the project root.

## Building the Application

To build the application, ensure you have Go installed on your machine. Then follow these steps:

1. Open a terminal and navigate to the project root directory.
2. Run the following command to build the application:

    ```sh
    go build -o todo-app .
    ```

This will generate an executable file named `todo-app` in the current directory.

## Running the Application

To run the application, follow these steps:

1. After building the application, run the following command:

    ```sh
    ./todo-app
    ```

2. The application will start and listen on port 9999 by default.

## Testing the Application

To run tests for the application and view coverage, follow these steps:

1. Open a terminal and navigate to the project root directory.
2. Run the following command to run tests and generate coverage:

    ```sh
    go test -coverprofile=coverage.out ./...
    ```

3. To view coverage in the terminal, run:

    ```sh
    go tool cover -func=coverage.out
    ```

4. To generate an HTML coverage report, run:

    ```sh
    go tool cover -html=coverage.out -o coverage.html
    ```

You can open the `coverage.html` file in a web browser to view the coverage report.

## Feature Flag

The application includes a feature flag that can be used to enable or disable certain endpoints. By default, the flag is disabled. To enable or disable the flag, you can use the `/toggleflag` endpoint using a POST request.

### Enabling the Feature Flag

To enable the feature flag, send a POST request to `/toggleflag` endpoint:

```sh
curl -X POST http://localhost:9999/toggleflag
```

### Disabling the Feature Flag

To disable the feature flag, send another POST request to the `/toggleflag` endpoint:

```sh
curl -X POST http://localhost:9999/toggleflag
```

## Dockerizing the Application

The application can also be containerized using Docker for easier deployment.

### Building the Docker Image

To build a Docker image of the application, follow these steps:

1. Open a terminal and navigate to the project root directory.
2. Run the following command to build the Docker image:

    ```sh
    docker build -t todoapp .
    ```

### Running the Docker Container

To run the Docker container from the built image, follow these steps:

1. After building the Docker image, run the following command to start a container:

    ```sh
    docker run -d --name todo-container -p 9999:9999 todoapp
    ```

2. The application will be accessible at `http://localhost:9999`.
