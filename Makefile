GOSWAGGER_DIR = ./tools/go-swagger/
GOSWAGGER=swagger_linux_amd64

check-go-swagger:
ifeq ($(strip $(GOSWAGGER)),)
	$(error no go-swagger for your platform ?)
endif

notes-api-generate: | check-go-swagger
	rm -rf notes-api/gen/restapi/operations/
	rm -rf notes-api/gen/restapi/doc.go
	rm -rf notes-api/gen/restapi/embedded_spec.go
	rm -rf notes-api/gen/restapi/server.go
	rm -rf notes-api/gen/models/
	$(GOSWAGGER_DIR)$(GOSWAGGER) generate server -t notes-api/gen -A notes -f notes-api/swagger/swagger.yaml --exclude-main
	make notes-client-generate

notes-client-generate:
ifeq ($(strip $(GOSWAGGER)),)
	$(error no go-swagger for your platform ?)
endif
	rm -rf notes-api/gen/client/
	$(GOSWAGGER_DIR)$(GOSWAGGER) generate client -t notes-api/gen -A notes  -f notes-api/swagger/swagger.yaml

cct-seed: 
	docker-compose up -d