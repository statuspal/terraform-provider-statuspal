default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ENV=TEST TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
