mock:
	mockgen -source ./catalog_interface.go -package mock_api -destination ./mock_api/catalog.go Catalog

test:
	go test -v -coverprofile=coverage.out
