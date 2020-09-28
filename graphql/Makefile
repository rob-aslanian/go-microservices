.PHONY: all schema resolver container build clean docker-push proto docker-clean

BUILDFOLDER=build/
OUTPUTNAME=graphql_server
BUILDVERSION=$(shell date +'%d%m%y-%H%M')
REMOTE_REPO=testing.lan:5000
# TAG=$(OUTPUTNAME):$(BUILDVERSION)

all: fixes test

schema:
	@echo 'Generating GraphQL schema...'
	@go generate ./schema/
	@echo 'Done'

resolver: schema
	@echo 'Generating GraphQL resolver...'
	@graphql-gen-go ./schema/graphql/* --out_dir ./ --pkg resolver
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

build: clean fixes test
	@echo 'Compiling...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILDFOLDER)$(OUTPUTNAME)
	@echo 'Done'

clean:
	@echo 'Cleaning...'
	@rm -rf $(BUILDFOLDER)
	@echo 'Done'

proto:
	@echo 'Generating grpc proto files...'
	@make ./../grpc-proto/Makefile
	@echo 'Done'

test:
	@echo 'Testing...'
	@go test
	@echo 'Done'

fixes: resolver
	@echo 'Applying fixes...'
	@sed -i -- 's/float32/float64/g' ./resolver/resolver.gql.go		# fix type error

	@sed -i -- 's/Is_hidden bool/Is_hidden *bool/g' ./resolver/resolver.gql.go		# fix type error
	@sed -i -- 's/return &r.R.Is_hidden/return r.R.Is_hidden/g' ./resolver/resolver.gql.go		# fix type error

	@sed -i -- 's/import (/import (\n	"gitlab.lan\/Rightnao-site\/microservices\/graphql\/graphiql"/g' ./resolver/server.gql.go	# add import
	@sed -i -- 's/import (/import (\n	"time"/g' ./resolver/resolver.gql.go	# add time import
	@sed -i -- 's/mux.Handle("\/graphql", srv)/mux.Handle("\/graphql", srv)\n	mux.Handle("\/graphiql", graphiql.GetGraphiql())/g' ./resolver/server.gql.go	# add graphiql
	@sed -i -- 's/GqlResolver/interface{}/g' ./resolver/server.gql.go	#	change interface type

	@sed -i -- 's/type ProfileResolver struct {/type ProfileResolver struct {\n language string/g' ./resolver/resolver.gql.go	# add language in ProfileResolver
	@sed -i -- 's/return ProfileResolver{/return ProfileResolver{R: /g' ./resolver/resolver.gql.go	# fix
	@sed -i -- 's/return SkillProfileResolver{/return SkillProfileResolver{R: /g' ./resolver/resolver.gql.go	# fix

	@sed -i -- 's/type SkillProfileResolver struct {/type SkillProfileResolver struct {\n language string/g' ./resolver/resolver.gql.go	# add language in SkillProfileResolver
	@sed -i -- 's/(items, ProfileResolver{/(items, ProfileResolver{R: /g' ./resolver/resolver.gql.go	# fix
	@sed -i -- 's/(items, SkillProfileResolver{/(items, SkillProfileResolver{R: /g' ./resolver/resolver.gql.go	# fix
	@sed -i -- 's/return &ProfileResolver{/return \&ProfileResolver{R: /g' ./resolver/resolver.gql.go	# fix

	@tr '\n' '@' < ./resolver/resolver.gql.go > ./resolver/resolver.gql.go.tmp

	@sed -i -E 's/for\s_,\sitm\s:=\srange\sr\.R\.(\w+)\s+\{\s*@\s*items\s*=\s*append\(items,\s*(\w+)Resolver\{&itm\}\)/for i := range r\.R\.\1 {@    items = append(items, \2Resolver{\&r\.R\.\1[i]})/g' ./resolver/resolver.gql.go.tmp # fixing loops
	@sed -i -E 's/for\s_,\sitm\s:=\srange\sr\.R\.(\w+)\s+\{\s*@\s*items\s*=\s*append\(items,\s*(\w+)Resolver\{R\:\s*&itm\}\)/for i := range r\.R\.\1 {@    items = append(items, \2Resolver{R:\&r\.R\.\1[i]})/g' ./resolver/resolver.gql.go.tmp # fixing loops for Profile
	@sed -i -E 's/return\s\&graphql\.Time\{Time:\s\*r\.R\.(\w+)\}/if r.R.\1 == nil {@ return nil@}@ return \&graphql\.Time\{Time: *r.R.\1\}/g' ./resolver/resolver.gql.go.tmp # fixing getting *time.Time value by address

	@#sed -i -E 's/type\s+NetworkInfoInUserProfile\sstruct\s+\{@\s+\w+\sint32@\s+\w+\sint32@\s+\w+\s+int32@\s+\w+\sint32@\s+\w+\s+int32/type NetworkInfoInUserProfile struct {@\sConnections *int32@ Followings *int32@ Followers \*int32@ Recommendations *int32@ Reviews *int32/g' ./resolver/resolver.gql.go.tmp # fixes pointer in NetworkInfoInUserProfile

	@tr '@' '\n' < ./resolver/resolver.gql.go.tmp > ./resolver/resolver.gql.go.2.tmp
	@rm ./resolver/resolver.gql.go.tmp
	@mv ./resolver/resolver.gql.go.2.tmp ./resolver/resolver.gql.go

	@gofmt -w ./resolver/resolver.gql.go

	@echo 'Done'
