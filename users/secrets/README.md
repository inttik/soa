# Secrets

Каталог содержит секреты, которые, по хорошему, не следовало бы загружать на гитхаб.

## JWT tokens
Для jwt токенов используются `signature.pem` и `signature.pub`. Генерируются с помощью следующих строчек:

Создаем приватный ключ:
```
openssl genrsa -out signature.pem 2048
```

По приватнму ключу получаем публичный:
```
openssl rsa -in signature.pem -outform PEM -pubout -out signature.pub
```
