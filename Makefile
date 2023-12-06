install: uninstall
	GO111MODULE=on go mod tidy
	GO111MODULE=on go build -ldflags "${LDFLAGS}" -o sem
	mv sem ${GOPATH}/bin

uninstall:
	-rm ${GOPATH}/bin/sem
