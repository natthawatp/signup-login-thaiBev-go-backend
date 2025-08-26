VERSION=v0.0.1

version:
	echo $(VERSION)

run:
	go run -ldflags "-X 'main.VERSION=$(VERSION)'" .
