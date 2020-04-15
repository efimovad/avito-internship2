# avito-internship
### Задача
Необходимо создать HTTP-сервис, способный ограничивать количество запросов (rate limit) из одной подсети IPv4. Если ограничения отсутствуют, то нужно выдавать одинаковый статический контент.

### Запуск проекта
`$:docker-compose up` - в докере  
`$:make && ./server` - локально  
`$:make test` - тесты  

**Параметры по умолчанию:**
- ограничение 100 запросов в минуту
- время ожидания после ограничения 2 минуты
- подсеть /24
- порт 8080

**Запросы:**
#### Получение статического контента
##### GET /
```
HTTP/1.1 200 OK
Date: Wed, 15 Apr 2020 14:52:02 GMT
Content-Length: 12
Content-Type: text/plain; charset=utf-8

HELLO FROM SERVER!
```
#### Сбросить предела запросов
##### POST /reset 
Сброс предела запросов от клиента в подсети /24 (IP берется из заголовка X-FORWARDED-FOR)  

```
HTTP/1.1 200 OK
Date: Wed, 15 Apr 2020 14:54:47 GMT
Content-Length: 11
Content-Type: text/plain; charset=utf-8

RESET LIMIT
```

Ответ, когда клиент превысил лимит запросов
```
HTTP/1.1 429 Too Many Requests
Content-Type: text/html
Retry-After: 1m59.339902611s
Date: Wed, 15 Apr 2020 14:56:02 GMT
Content-Length: 200

<html>
       <head>
               <title>Too Many Requests</title>
       </head>
       <body>
               <h1>Too Many Requests</h1>
               <p>I only allow 100 requests per period to this Web site per net. Try again soon.</p>
       </body>
</html>
```

Запуск с разными параметрами
```
Usage of ./server:
  -cidr int
        enter mask (default 24)
  -limit int
        enter requests limit per period (default 100)
  -period duration
        enter period in seconds (default 1m0s)
  -port int
        enter port (default 8080)
  -wait duration
        enter wait in seconds (default 2m0s)
```

