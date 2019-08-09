run:
	# https://github.com/lambci/docker-lambda
	GOOS=linux go build -o app
	docker-compose up
	rm app
zip:
	# https://docs.aws.amazon.com/lambda/latest/dg/lambda-go-how-to-create-deployment-package.html