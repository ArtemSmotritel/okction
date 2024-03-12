run:
	@templ generate
	@go run main.go

build:
	@templ generate
	@go build -o ./tmp/main .

air:
	@air
