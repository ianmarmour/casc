.PHONY: build
build: 
	cd ./cmd/casc-extract && go build

.PHONY: test
test:
	go test -race -v ./...

.PHONY: online
online: build
	cd ./cmd/casc-extract && ./casc-extract -app $(CASC_APP) -region $(CASC_REGION) -cdn $(CASC_REGION) -pattern "$(CASC_PATTERN)" -o "extract/$(CASC_APP)/online" -v

.PHONY: local
local: build
	cd ./cmd/casc-extract && ./casc-extract -dir "$(CASC_DIR)" -pattern "$(CASC_PATTERN)" -o "extract/$(CASC_APP)/local" -v

.PHONY: testslow
testslow:
	go test -race -slow -v ./... -timeout 120m -app $(CASC_APP)

.PHONY: testupdate
testupdate:
	go test -race -slow -update -v ./... -timeout 120m -app $(CASC_APP)