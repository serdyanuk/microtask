# Microtask Application
## Docker development usage:
### Start app

```
$ docker-compose up 
```

### Down app
```
$ docker-compose down
```
### Swagger UI
```
http://localhost:3000/
```

## Services configuration:
```
./config/config.yaml
```

## Solution notes

1. Images are uploaded via Files-api and saved to disk, after which Process Service
receives a message about the arrival of a new file, which contains: file ID, Width, Height and physical size in kilobytes

2. Service API loads the required image from disk, and applies optimization to it, resizing it with the specified ResizePower parameter from the configuration file, after which the result and new image data are output to the console

3. Two image formats are supported: jpg / jpeg and png
## Task
https://softcery.notion.site/Test-Task-Junior-Go-Developer-0b649d4310454c8588610b15a63eb376