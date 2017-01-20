SMS Service
============

Сервис для смс-рассылок

Installation
------------

```sh
make
make install
```

Configuration
-------------

Format:

```yml
project: phone
listen: 3006

rabbit:
  host: 127.0.0.1
  port: 5672
  username: guest
  password: guest

sms:
  -
    codes: []
    enable: true
    settings:
      url: 'https://example.com'
      login: login
      password: password
      from: FROM
    transport: example

```

Starting service
----------------

```sh
bin/sms-service --c cfg/config.yml
```