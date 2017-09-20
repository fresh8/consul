mock:
	mockgen -source ./api_if/catalog.go -package mock_api -destination ./mock_api/catalog.go Catalog

test:
	go test -v
