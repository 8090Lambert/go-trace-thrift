BINARY_FILE := "tracer_thrift"
SOURCE_FILE := `find . -name "*.go" | grep -v "vendor/*" | grep -v "example/*"`
VET_PACKAGES := `go list ./... | grep -v "vendor/"`
TEST_FOLDER := `go list ./... | grep -v "vendor/" | grep -v "example/"`

all: install

install: clean build
	go mod vendor

fmt:
	gofmt -s -w ${SOURCE_FILE} > /dev/null

fmt_check:
	@diff=$$(gofmt -s -d ${SOURCE_FILE}); \
	if [ -n "$${diff}" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

test:
	echo "mode: count" > coverage.out
	for d in ${TEST_FOLDER}; do \
		go test -v -covermode=count -coverprofile=profile.out $$d > tmp.out > /dev/null; \
		cat tmp.out; \
		if grep -q "^--- FAIL" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "build failed" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "setup failed" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		fi; \
		if [ -f profile.out ]; then \
			cat profile.out | grep -v "mode:" >> coverage.out; \
			rm profile.out tmp.out; \
		fi; \
	done

clean:
	@if [ -f ${BINARY_FILE} ]; then rm ${BINARY_FILE} ; fi

vet:
	go vet $(VET_PACKAGES) > /dev/null

build:
	go build
