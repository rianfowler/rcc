# RCC Test Go Application

This is a simple Go application used for testing the RCC build system. It provides a basic HTTP server that responds with a greeting message.

## Running the Application

```bash
go run main.go
```

The server will start on port 8080 and can be accessed at http://localhost:8080/

## Building with RCC

To build this application using RCC:

```bash
# Using GoReleaser
rcc build go v1.0.0

# Using Cloud Native Buildpacks
rcc build pack . --name go-test-app
```