docker load -i pfc300dynmodapi.tar


docker run -d -p 8080:8080 --name my-api dynmodapi:arm64
