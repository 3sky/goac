# GOAC(GO Applications Control)

App contains basic http auth with `3sky:test` hardcoded in app.go, feel free to change it :)

## Install

- Fetach deps(I used go modules)

  ```bash
  go get
  ```

- Build app

  ```bash
  make build
  ```

- Run server !

  ```bash
  make run
  ```

### API Base command usage

- Almost everywhere `-u user:password` is mandatory
- API accpect fallwoing methods:

    - GET `curl -s http://127.0.0.1:5000/api/app/1`
    - POST `curl -H "Content-Type: application/json" -X POST -d '{"app_name": "GO_API", "app_version": "2.123", "environment": "stg"}' http://127.0.0.1:5000/api/app/new`
    - PUT `curl -i -H "Content-Type: application/json" -X PUT -d '{"app_name": "GO_API", "app_version": "0.95"}' http://127.0.0.1:5000/api/app/2`
    - DELETE `curl -s -X DELETE http://127.0.0.1:5000/api/app/2`

- POST/PUT available parameters

    - app_name, app_version, environment, ip, branch, update_by


### CLI base usage

- Client configure `client/.creds`

    ```json
    {
    "creditional": {
        "user": "3sky",
        "password": "test"
    },
    "server": {
        "ip": "127.0.0.1",
        "port": 5000
    }
    }
    ```

- Get help

    ```bash
    ./appClient -help
    ```

- Sample search

    ```bash
   ./appClient -action search -app GO_API -env dev
    ```
