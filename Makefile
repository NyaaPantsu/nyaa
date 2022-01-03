OUT=nyaa

all: build

build: $(OUT)

$(OUT):
	godep go build -ldflags "-X main.buildversion=$(git rev-parse HEAD)"

clean:
	go clean 
