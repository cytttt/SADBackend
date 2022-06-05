# SAD 2022 Final Project Backend

## requirement
`go 1.16`

## run server
`make install`

`./SADBackend`

## run unit test
`make test`

## directory structure

```
root/ -- routers/                     # routing path
      |- controllers/ -- v1/          # API definition
      |               |- service/     # utility function for API
      |- pkg/mongodb/                 # mongodb query method
      |- repo/                        # database interface
      |- model/                       # data schema
      |- constant/                    # constant variables and response msg
      |- env/                         # environment variables
      |- main.go
```
