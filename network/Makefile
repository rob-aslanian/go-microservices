.PHONE: build docker-push clean

OUTPUTNAME=network
REMOTE_REPO=testing.lan:5000
BUILDFOLDER=build/

build: clean
	@echo 'Compiling...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILDFOLDER)$(OUTPUTNAME) ./cmd/network/
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
