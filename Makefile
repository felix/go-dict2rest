
GOPATH := ${PWD}
export GOPATH

default: build

build: dict2rest

dict2rest:
	go install dict2rest

clean:
	rm bin/*

.PHONY: build dict2rest clean
