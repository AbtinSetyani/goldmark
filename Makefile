.PHONY: test fuzz

test:
	go test -coverprofile=profile.out -coverpkg=github.com/anytypeio/goldmark,github.com/anytypeio/goldmark/ast,github.com/anytypeio/goldmark/extension,github.com/anytypeio/goldmark/extension/ast,github.com/anytypeio/goldmark/parser,github.com/anytypeio/goldmark/renderer,github.com/anytypeio/goldmark/renderer/html,github.com/anytypeio/goldmark/text,github.com/anytypeio/goldmark/util ./...

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
