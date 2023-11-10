install: uninstall
	GO111MODULE=on go mod tidy -go=1.21
	GO111MODULE=on go build -ldflags "${LDFLAGS}" -o sem
	mv sem ${GOPATH}/bin

uninstall:
	-rm ${GOPATH}/bin/sem
