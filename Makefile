GO=go

build:
	$(GO) build -i -o ./output/hec .

clean:
	rm -rf ./output