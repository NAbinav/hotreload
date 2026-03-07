build:
	go build -o hotreload .

demo: build
	./hotreload --root ./testserver --build "cd testserver && go build -o ../bin/http ." --exec "./bin/http"

