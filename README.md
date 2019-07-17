# ProWD Back End

The back-end service of [ProWD]((http://prowd.id)) is used as database wrapper to save CFA profiles which ProWD's users want to analyze

## Prerequisites & Setup

- Linux machine with root access
- Golang - [Go installation](https://golang.org/doc/install)
- [Govendor](https://github.com/kardianos/govendor)

```
git clone https://github.com/avcwisesa/wikidata-profiling-prototype-back
```

## Running the code

1. Copy environment variables
```
cp .env.example .env
```

2. Modify environment variable to fit database spec & back-end server port

3. Install dependencies
```
govendor fetch +outside
```
4. For running program, use command below:
```
go run main.go
```

## Development Guide

### Libraries
- [Gin-Gonic](https://github.com/gin-gonic)
- [GORM](http://gorm.io/)

### Modifying code
- To maintain lean structure and maintainability the code is currently only in [main.go](https://github.com/avcwisesa/wikidata-profiling-prototype-back/blob/master/main.go), future needs may dictates otherwise.
- To modify database field, modify `type profile struct` in main.go, and modify `createProfile` and `updateProfile` functions accordingly
- Migration will be done automatically after the back-end server has been restarted.
- To add new feature/endpoint refer to Gin-Gonic docs for routing and handling requests
