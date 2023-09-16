
test:
	go test ./pkg/thyrsus

fmt:
	go fmt $$(go list ./...)
