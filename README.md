docker load -i modbusapi.tar


docker run -d -p --network=host --name my-api modbusapipfc300

curl -X POST http://localhost:8080/ModuleOneWord      -H "Content-Type: application/json"      -d '{"value": "7"}'

//So far only works for ModuleOneWord to ModuleThreeWord

curl -X GET http://localhost:8080/OutputTotal
