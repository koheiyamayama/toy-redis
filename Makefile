run-k6:
	./k6 run loadtesting/index.js

build-k6:
	xk6 build master --with github.com/NAlexandrov/xk6-tcp
