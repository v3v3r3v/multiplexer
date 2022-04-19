Curl:
```shell
curl --location --request POST 'http://localhost:8990/multiplex' \
--header 'Content-Type: application/json' \
--data-raw '[
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000",
  "http://stub:8991/stub?limit=1000"
]'
```