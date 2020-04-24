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

<html>
    <head>
        <title>Hello</title>
    </head>
    
    <body>
        <h1>Hello</h1>
        <p>Hello from server</p>
    </body>
</html>
```
#### Сбросить предела запросов
##### POST /admin/reset?login=admin&password=[ПАРОЛЬ] 
Сброс предела запросов от клиента в подсети /24 (IP берется из заголовка X-FORWARDED-FOR)  
При запуске сервера с помощью флага --pass можно задать пароль (по умолчанию он 123456)   

```
HTTP/1.1 200 OK
Date: Wed, 15 Apr 2020 14:54:47 GMT
Content-Length: 11
Content-Type: text/plain; charset=utf-8

<html>
    <head>
        <title>Reset</title>
    </head>
    
    <body>
        <h1>Reset</h1>
        <p>Limit was reset</p>
    </body>
</html>
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

Ответ, когда выставлен неправильный IP в заголовке X-Forwarded-For
```
HTTP/1.1 403 Forbidden
Content-Type: text/html
Date: Fri, 24 Apr 2020 20:23:32 GMT
Content-Length: 143

<html>
        <head>
                <title>Wrong IP</title>
        </head>
        <body>
                <h1>Wrong IP</h1>
                <p>X-Forwarded-For header include wrong ID</p>
        </body>
</html>
```

Ответ, когда не введены или введены неправильно логин и пароль админа
```
HTTP/1.1 403 Forbidden
Content-Type: text/html
Date: Fri, 24 Apr 2020 20:25:32 GMT
Content-Length: 170
<html>
    <head>
        <title>No access to admin request</title>
    </head>
    
    <body>
        <h1>No access to admin request</h1>
        <p>Enter admin login and password</p>
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
  -pass string
        enter admin password (default "123456")
  -period duration
        enter period in seconds (default 1m0s)
  -port int
        enter port (default 8080)
  -wait duration
        enter wait in seconds (default 2m0s)
```

