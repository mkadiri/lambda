run:
	# https://github.com/lambci/docker-lambda
	GOOS=linux go build main.go
	docker run --rm -v ${CURDIR}:/var/task lambci/lambda:go1.x main '{"What is your name?": "John", "How old are you?": 9}'