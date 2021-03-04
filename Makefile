clean:
	rm -f whisper
all:
	rm -f whisper
	go build -o whisper  ./cmd/web
