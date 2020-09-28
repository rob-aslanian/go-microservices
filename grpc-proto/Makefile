.PHONE: all

all:
	protoc --go_out=plugins=grpc:. ./userRPC/*.proto
	protoc --go_out=plugins=grpc:. ./chatRPC/*.proto
	protoc --go_out=plugins=grpc:. ./authRPC/*.proto
	protoc --go_out=plugins=grpc:. ./infoRPC/*.proto
	protoc --go_out=plugins=grpc:. ./mailRPC/*.proto
	protoc --go_out=plugins=grpc:. ./companyRPC/*.proto
	protoc --go_out=plugins=grpc:. ./networkRPC/*.proto
	protoc --go_out=plugins=grpc:. ./searchRPC/*.proto
	protoc --go_out=plugins=grpc:. ./servicesRPC/*.proto
	protoc --go_out=plugins=grpc:. ./jobsRPC/*.proto
	protoc --go_out=plugins=grpc:. ./notificationsRPC/*.proto
	protoc --go_out=plugins=grpc:. ./advertRPC/*.proto
	protoc --go_out=plugins=grpc:. ./statisticsRPC/*.proto
	protoc --go_out=plugins=grpc:. ./stuffRPC/*.proto
	protoc --go_out=plugins=grpc:. ./newsfeedRPC/*.proto
	protoc --go_out=plugins=grpc:. ./groupsRPC/*.proto
	protoc --go_out=plugins=grpc:. ./shopRPC/*.proto
	protoc --go_out=plugins=grpc:. ./rentalRPC/*.proto

	@# fixing import path for search.proto
	@sed -i -- 's/userRPC \"userRPC\"/userRPC "gitlab.lan\/Rightnao-site\/microservices\/grpc-proto\/userRPC"/' ./searchRPC/search.pb.go
	@sed -i -- 's/jobsRPC \"jobsRPC\"/jobsRPC "gitlab.lan\/Rightnao-site\/microservices\/grpc-proto\/jobsRPC"/' ./searchRPC/search.pb.go
	@sed -i -- 's/companyRPC \"companyRPC\"/companyRPC "gitlab.lan\/Rightnao-site\/microservices\/grpc-proto\/companyRPC"/' ./searchRPC/search.pb.go
	@sed -i -- 's/servicesRPC \"servicesRPC\"/servicesRPC "gitlab.lan\/Rightnao-site\/microservices\/grpc-proto\/servicesRPC"/' ./searchRPC/search.pb.go

	@sed -i -- 's/userRPC \"userRPC\"/userRPC "gitlab.lan\/Rightnao-site\/microservices\/grpc-proto\/userRPC"/' ./advertRPC/advert.pb.go
	@sed -i -- 's/companyRPC \"companyRPC\"/companyRPC "gitlab.lan\/Rightnao-site\/microservices\/grpc-proto\/companyRPC"/' ./advertRPC/advert.pb.go

	@sed -i -- 's/companyRPC \"companyRPC\"/companyRPC "gitlab.lan\/Rightnao-site\/microservices\/grpc-proto\/companyRPC"/' ./userRPC/user.pb.go

	@sed -i -- 's/userRPC \"userRPC\"/userRPC "gitlab.lan\/Rightnao-site\/microservices\/grpc-proto\/userRPC"/' ./networkRPC/network.pb.go
	@sed -i -- 's/companyRPC \"companyRPC\"/companyRPC "gitlab.lan\/Rightnao-site\/microservices\/grpc-proto\/companyRPC"/' ./networkRPC/network.pb.go
