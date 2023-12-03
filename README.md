<h1 align="center">NotVulnApp</h1>
<p align="center" >Защищённая версия <a href="https://github.com/Stepashka20/vuln-app">приложения</a></p>

## 📃 Инструкция по сборке и запуску приложения

### Требования

Для успешной сборки этого приложения вам потребуется:

- Go (версия 1.18 и выше)
- gcc (установленный и доступный в системе, необходим для sqlite3)

### Сборка и запуск

Для сборки приложения выполните следующие команды:

```bash
git clone https://github.com/Stepashka20/not-vuln-app
cd vuln-app
go build -o vuln-app
```

Перед запуском скопируйте и заполните `.env`:
```bash
cp .env.example .env
```

Для запуска приложения выполните следующую команду:

```bash
./vuln-app
```

После этого приложение будет доступно по адресу http://localhost:8080

## 👨🏻‍💻 Комментарии к исправлениям

### 1. XSS [Коммит](https://github.com/Stepashka20/not-vuln-app/commit/c9ee1b384339aa1d273f08d16b60c87e1c6b2f38)
1. Для защиты от XSS атаки я использовал функцию `template.HTMLEscapeString`, которая экранирует все спецсимволы в строке. 
2. Также можно рендерить темплейт через `gin.H{"username": username, "filename": favoriteFilename}`, который также автоматически экранирует спецсимволы

### 2. SQLI [Коммит](https://github.com/Stepashka20/not-vuln-app/commit/41ad0d0c984da5366909a63c69160dcee3801cae)
1. Для защиты от SQLI был использован подготовленный запрос, который не даст внедрить SQL код в запрос:
```go
("SELECT username, password FROM users WHERE username = ?", username)
```

### 3. Brute force [Коммит](https://github.com/Stepashka20/not-vuln-app/commit/97d3cda5910ea1a6c34199b1493089a6efb54235)
1. Для защиты от брутфорса я добавил middleware c ограничением на количество запросов с одного IP в секунду. Теперь код ниже выполняет всего 2 запроса в секунду, а остальные 48 запросов возвращаются с 429 ошибкой.

```bash
seq 50 | xargs -P50 -I{} curl -s -o /dev/null -w "%{http_code}\n" 'http://localhost:8080/login' --data-raw 'username=admin&password={}'
```

### 4. Path Traversal [Коммит](https://github.com/Stepashka20/not-vuln-app/commit/f621df051af2a972ca81ed84c03ade6de3d0a1c5)
1. Для защиты от этой уязвимости я не сразу отдаю файл пользователю. Я очищаю путь от спец символов, а затем проверю находится ли файл в разрешённой директории `uploads`:
2. Теперь файлы возможно скачать только из папки `uploads`

### 5. OS command injection [Коммит](https://github.com/Stepashka20/not-vuln-app/commit/9776b63dba02df50ca280158a23ebbde87d32dea)
1. Для защиты от OS command injection я сначала валидирую пользовательский ввод, а затем использую `exec.Command` с утилитой `ping`, и подставлю адрес как один из параметров.

### 6. IDOR [Коммит](https://github.com/Stepashka20/not-vuln-app/commit/81b3a2febf616e0f682e03a5e2dd5c0ff654f11c)
1. Для защиты от IDOR я добавил проверку на то, что пользователь может удалить только свой аккаунт, сверяя name из GET параметра и из куков.