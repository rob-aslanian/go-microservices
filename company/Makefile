.PHONY: clean build run clean docker-build docker-push-testing docker-clean

OUTPUTNAME=company
REMOTE_REPO=testing.lan:5000
BUILDFOLDER=build/

docker-build: build
	@echo 'Building container...'
	@docker build --rm -q -t $(OUTPUTNAME) .
	@echo 'Done'

docker-push-testing: docker-build
	@echo 'Pushing Container to '$(REMOTE_REPO)'...'
	@docker tag $(OUTPUTNAME) $(REMOTE_REPO)/$(OUTPUTNAME):latest
	@docker push $(REMOTE_REPO)/$(OUTPUTNAME):latest
	@echo 'Done'

docker-clean:
	@docker rmi $$(docker images -f "dangling=true" -q)
	@docker rm $$(docker ps -f "status=exited" -q)

build: clean test
	@echo 'Build...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILDFOLDER)$(OUTPUTNAME) cmd/company/main.go
	@echo 'Done'

run:
	@echo 'Run...'
	go run cmd/company/main.go
	@echo 'Done'

clean:
	@echo 'Cleaning...'
	@rm -rf $(BUILDFOLDER)
	@echo 'Done'
	
test: 
	@go test ./pkg/service/
	@echo "Tests Passed Succesfully"
