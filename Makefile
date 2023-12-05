run-k6:
	./k6 run ./loadtesting/dist/main.bundle.js

build-k6:
	xk6 build master --with github.com/NAlexandrov/xk6-tcp

run-toy-redis-on-debug-mode:
	TOY_REDIS_LOG_LEVEL="-4" go run .

run-toy-redis:
	go run .

deploy:
	GOOS=linux GOARCH=amd64 go build
	scp toy-redis kohei@192.168.1.12:/home/kohei
