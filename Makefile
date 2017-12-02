all: build/rdbfs

build/%:
	go build -o $@ -tags netgo -ldflags "-w" main.go

clean:
	rm -rf build

.PHONY: all clean