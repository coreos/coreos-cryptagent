.PHONY: vendor
vendor: glide.yaml
	@glide update --strip-vendor
