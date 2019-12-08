.PHONY: test fuzz

test:
	go test -coverprofile=profile.out -coverpkg=github.com/enkogu/goldmark,github.com/enkogu/goldmark/ast,github.com/enkogu/goldmark/extension,github.com/enkogu/goldmark/extension/ast,github.com/enkogu/goldmark/parser,github.com/enkogu/goldmark/renderer,github.com/enkogu/goldmark/renderer/html,github.com/enkogu/goldmark/text,github.com/enkogu/goldmark/util ./...

cov: test
	go tool cover -html=profile.out

fuzz:
	which go-fuzz > /dev/null 2>&1 || (GO111MODULE=off go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build; GO111MODULE=off go get -d github.com/dvyukov/go-fuzz-corpus; true)
	rm -rf ./fuzz/corpus
	rm -rf ./fuzz/crashers
	rm -rf ./fuzz/suppressions
	rm -f ./fuzz/fuzz-fuzz.zip
	cd ./fuzz && go-fuzz-build
	cd ./fuzz && go-fuzz
