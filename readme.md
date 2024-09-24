# Материалы для семинаров 2024 года

## Генерация документации по аннотациям в коде

``` swag init -d cmd -o api  ```

## Генерация кода из описания API

``` go generate ./...```

## Генерация кода из proto файла

```
protoc \
-I ./libs \
\
--grpc-gateway_out=logtostderr=true:./api  \
--grpc-gateway_opt=paths=source_relative \
\
--openapiv2_out ./api \
--openapiv2_opt logtostderr=true \
--openapiv2_opt output_format=yaml \
\
--go_out=./api \
--go_opt=paths=source_relative \
\
--go-grpc_out=./api \
--go-grpc_opt=paths=source_relative \
\
--proto_path=./api \
api/*.proto
```

 ## Ссылки
- [Postman](https://www.postman.com)
- [curl](https://github.com/curl/curl/blob/master/docs/MANUAL.md)

#### --

- [Онлайн редактор Swagger](https://editor.swagger.io)
- [Генерация описания по аннотациям в коде.  Swag](https://github.com/swaggo/swag)

#### --

- [Спецификация Swagger](https://docs.swagger.io/spec.html)
- [Спецификация OpenAPI](https://github.com/OAI/OpenAPI-Specification)

#### -- 

- [Генератор кода из OpenAPI описания](https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file)

#### --

- [Документация Protobuf](https://protobuf.dev/getting-started/gotutorial/)
- [Генератор кода из proto файлов](https://protobuf.dev/reference/go/go-generated/)

#### -- 

- [gRPC-to-HTTP](https://github.com/grpc-ecosystem/grpc-gateway)
- [Опции прото файла для gRPC-to-HTTP](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto#L46)

#### --

-[Статья про документацию API](https://habr.com/ru/articles/496574/)
