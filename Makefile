run-k6:
	./k6 run loadtesting/index.js

build-k6:
	xk6 build master --with github.com/NAlexandrov/xk6-tcp

run-toy-redis-on-debug-mode:
	TOY_REDIS_LOG_LEVEL="-4" go run .

run-toy-redis:
	go run .
