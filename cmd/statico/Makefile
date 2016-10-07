AUTHOR = dethi
PROJECT = statico

all: compile image publish

compile:
	GOOS=linux GOARCH=amd64 go build -v -a -ldflags '-s -w'

image:
	docker build -t ${AUTHOR}/${PROJECT} .

publish:
	docker push ${AUTHOR}/${PROJECT}

clean:
	go clean

.PHONY: all compile image publish clean
