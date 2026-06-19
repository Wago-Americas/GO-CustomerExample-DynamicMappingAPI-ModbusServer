docker load -i pfc300dynmodapi.tar


docker run -d -p 8080:8080 --name my-api dynmodapi:arm64

curl -X POST http://localhost:8080/ModuleOneWord      -H "Content-Type: application/json"      -d '{"value": "7"}'

//So far only works for ModuleOneWord to ModuleThreeWord

curl -X GET http://localhost:8080/OutputTotal
