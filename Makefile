.PHONY: vendor
vendor: glide.yaml
	@glide update --strip-vendor
	@glide-vc --use-lock-file --no-tests --only-code
