clean:
	@-rm position.yaml
	@-rm daily-news-feed

build:
	@-go build -o $(app)
