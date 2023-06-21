MODULE=github.com/hjwalt/flows

test:
	go test ./... -cover -coverprofile cover.out
	
testv:
	go test ./... -cover -coverprofile cover.out -v

cov: test
	go tool cover -func cover.out

htmlcov: test
	go tool cover -html cover.out -o cover.html

# --------------------

tidy:
	go mod tidy
	go mod vendor
	go fmt ./...

update:
	go get -u ./...
	go mod tidy
	go mod vendor
	go fmt ./...

# --------------------

run:
	./script/run.sh

reset:
	./script/reset.sh

listen:
	./script/listen.sh

group-delete:
	./script/group-delete.sh
	
# --------------------

mocks: RUN
	mockgen -source=test_helper/interfaces.go -destination=test_helper/implementations.go -package=test_helper ;\
	mockgen -source=stateful_bun/connection.go -destination=test_helper/stateful_bun_connection.go -package=test_helper ;\

proto: RUN
	rm -rf $$GOPATH/$(MODULE)/ ;\
	protoc -I=. --go_out=$$GOPATH **/*.proto ;\
	cp -r $$GOPATH/$(MODULE)/* .

# --------------------

RUN:
