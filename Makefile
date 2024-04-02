build: 
	@go build

clean:
	@rm -rf data && mkdir data

run: 
	@go build && ./tld
