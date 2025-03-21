# Users
Цель сервиса - регистрировать новых пользователей, аутентифицировать их, а еще хранить персональную информацию о них, например ФИО или почту.

Есть следующие ручки:

* `(POST)` `/v1/register` - Регистрация нового пользователя. Клиент отправляет логин, пароль, почту, потом сервер регистрирует пользователя и возвращает его `id`.
* `(POST)` `/v1/login` - Вход.  Клиент отправляет логин и пароль, сервер возвращает `JWT` токен для пользователя, если пароль верный
* `(GET)` `/v1/user/{login}` - Получение `user_id`. Возвращает `id` пользователя с таким логином, или 404, если такого пользователя нет.
* `(GET)` `/v1/profile/{user_id}` - Доступ к странице пользователя. Возвращается страница пользователя без некоторых данных, если клиент не имеет к ней доступа.
* `(POST)` `/v1/profile/{user_id}` - Доступ к странице пользователя. Обновляется страница и возвращается новая с актуальной информацией.

* `(GET)` `/v1/friends/{user_id}` - (not implemented) Возвращает список всех (только публичных, если клиент не имеет доступа) друзей пользователя.

* `(POST)` `/v1/friends/{user_id}` - (not implemented) Позволяет обновить информацию о дружбе с другим пользователем.
