compile:
	echo "🧱 Compiling low latency key value database"
	go mod tidy
	go build -o llkvdb ./cmd/kvs  

run:
	./llkvdb start &

test:
	go test ./... -v

clean:
	rm -rf wal.txt 
	rm -rf llkvdb 