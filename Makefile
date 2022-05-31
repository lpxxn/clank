build:
	cd ./clank && go build -ldflags "-X main.version=`git tag --sort=-version:refname | head -n 1`" && go install