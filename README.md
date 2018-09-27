# Work in progress

Базовый код для заливки данных о товаре в CAS CL5000J/CL3000. Сам
протокол довольно просто и описан тут https://github.com/alexesDev/cas/blob/master/docs/protocol.pdf

Все написано на коленке, поэтому стоит использовать только на свой страх и
риск. Из фишек, которые постараюсь реализовать в свободное время:

 - [ ] обработку всех ошибок
 - [ ] web-server, который принимает JSON с массовом заданий
 - [x] cli, который принимает JSON с массовом заданий
 - [ ] тесты

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

Пример запуска:
```bash
go get github.com/alexesDev/cas/cmd/cascli
cascli example_task.json
```
