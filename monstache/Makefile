.PHONE: docker-push

OUTPUTNAME=monstache_custom
REMOTE_REPO=testing.lan:5000
BUILDFOLDER=build/


docker-build:
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
