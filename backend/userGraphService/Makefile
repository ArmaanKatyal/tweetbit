run:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto

test:
	@go test -coverprofile cover.out -cover github.com/ArmaanKatyal/tweetbit/backend/userGraphService \
	github.com/ArmaanKatyal/tweetbit/backend/userGraphService/utils \
	github.com/ArmaanKatyal/tweetbit/backend/userGraphService/services \
	github.com/ArmaanKatyal/tweetbit/backend/userGraphService/helpers

	@go tool cover -func cover.out | grep total

cover:
	@go tool cover -html=cover.out