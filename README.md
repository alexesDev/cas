# Work in progress

Базовый код для заливки данных о товаре в CAS CL5000J/CL3000. Сам
протокол довольно просто и описан тут https://github.com/alexesDev/cas/blob/master/docs/protocol.pdf

Все написано на коленке, поэтому стоит использовать только на свой страх и
риск. Из фишек, которые постараюсь реализовать в свободное время:

 - [ ] обработку всех ошибок
 - [ ] web-server, который принимает JSON с массовом заданий
 - [x] cli, который принимает JSON с массовом заданий
 - [ ] тесты
 - [ ] добавить exe для Windows

## cascli

Утилита, выполняющее JSON задание следующего вида:

```json
{
  "Addr": "192.168.88.250:20000",
  "Plan": [{
    "Type": "ErasePLU",
    "Input": {
      "DepartmentNumber": 1,
      "PLUNumber": 1
    }
  }, {
    "Type": "DownloadPLU",
    "Input": {
      "ScaleId": 0,
      "Data": {
        "DepartmentNumber": 1,
        "PLUName1": "Привет мир cp1251",
        "PLUType": 1,
        "PLUNumber": 1
      }
    }
  }, {
    "Type": "UploadPLU",
    "Input": {
      "ScaleId": 0,
      "PLUNumber": 1
    }
  }, {
    "Type": "UploadPLU",
    "Input": {
      "ScaleId": 0,
      "PLUNumber": 2
    }
  }]
}
```

Доступные поля можно смотреть тут https://github.com/alexesDev/cas/blob/master/pkg/cas/main.go#L29

Пример запуска:
```bash
go get github.com/alexesDev/cas/cmd/cascli
cascli example_task.json
```

Или в Docker:
```bash
docker run --rm -it \
  -v $(pwd)/example_task.json:/task.json cascli /task.json
  alexes/cas cli /task.json
```
