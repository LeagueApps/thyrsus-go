
test:
	@bash ./scripts/tests.sh

fmt:
	go fmt $$(go list ./...)
