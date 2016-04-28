gh-create-issues: gh-create-issues.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $@ $<
	strip $@

lint:
	@ go get -v github.com/golang/lint/golint
	@for file in $$(git ls-files '*.go' | grep -v '_workspace/'); do \
		export output="$$(golint $${file} | grep -v 'type name will be used as docker.DockerInfo')"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
	done; \
	exit $${status:-0}

vet:
	go vet gh-create-issues.go

imports:
	goimports -d gh-create-issues.go

test: lint vet imports
