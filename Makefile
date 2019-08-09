run:
	# https://github.com/lambci/docker-lambda
	GOOS=linux go build main.go
	docker-compose up
	rm main
zip:
	# https://docs.aws.amazon.com/lambda/latest/dg/lambda-go-how-to-create-deployment-package.html