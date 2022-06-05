# SAD 2022 Final Project Backend

## Requirement
`go 1.16`

## API Document
please refer to `Document.md`
## Run Server
`make install`

`./SADBackend`

## Run Unit Test
`make test`

## Directory Structure

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
