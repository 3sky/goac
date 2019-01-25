# sample_api_go

First try to convert https://github.com/3sky/sample-python-api into golang

## Base command usage

- Get one app

  ```commandline
  curl -s http://127.0.0.1:5000/api/app/1 | jq .
  ```

- Get all app

  ```commandline
  curl -s http://127.0.0.1:5000/api/apps | jq .
  ```

- Add new app

  ```commandline
  curl -i -H "Content-Type: application/json" -X POST -d '{"app_name": "Some App", "app_version": "2.123", "updated_by": "Kuba"}' http://127.0.0.1:5000/api/app/new
  ```

- Update data
  
  ```commandline
  curl -i -H "Content-Type: application/json" -X PUT -d '{"app_name": "GO API", "app_version": "0.95"}' http://127.0.0.1:5000/api/app/2
  ```

- Delete data

  ```commandline
  curl -s -X DELETE http://127.0.0.1:5000/api/app/2 | jq .
  ```

- Get HTML table

  ```commandline
  curl -s http://127.0.0.1:5000/
  ```
