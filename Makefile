all: prepare build

prepare:
	go get -u github.com/rakyll/statik

build:
	statik -f -src=swagger_dist -p saunter
	mv saunter/statik.go statik.go
	rm -rf saunter
