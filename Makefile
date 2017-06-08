all: test install run
install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_rsync/*.go
test:
	GO15VENDOREXPERIMENT=1 go test -cover `glide novendor`
vet:
	go tool vet .
	go tool vet --shadow .
lint:
	golint -min_confidence 1 ./...
errcheck:
	errcheck -ignore '(Close|Write)' ./...
check: lint vet errcheck
run:
	backup_rsync \
	-logtostderr \
	-v=4 \
	-one-time \
	-host=bborbe.devel.lf.seibert-media.net \
	-user=backup \
	-port=22 \
	-privatekey=/Users/bborbe/Documents/backup-ssh-keys/id_rsa \
	-source=/opt/apache-maven-3.3.9/ \
	-target=/backup/
format:
	find . -name "*.go" -exec gofmt -w "{}" \;
	goimports -w=true .
prepare:
	npm install
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	glide install
update:
	glide up
clean:
	rm -rf vendor
