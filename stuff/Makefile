.PHONE: all build docker-push clean docker-push-testing

OUTPUTNAME=stuff
REMOTE_REPO=testing.lan:5000
BUILDFOLDER=build/

build: clean test
	@echo 'Compiling...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILDFOLDER)$(OUTPUTNAME) ./cmd/stuff/
	@echo 'Done'

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

clean:
	@echo 'Cleaning...'
	@rm -rf $(BUILDFOLDER)
	@echo 'Done'

test:
	@go test ./pkg/service/
	@echo "Tests Passed Succesfully"
