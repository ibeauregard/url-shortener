build:
	docker build --tag=url_shortener:latest .
	docker image prune -f

run:
	docker run --interactive --tty --name=url_shortener_container \
    		--volume=url_shortener_db:/home/urlshortener/db --publish=8080:8080 --rm url_shortener:latest

restart: clean run

clean:
	@docker rm -f url_shortener_container &>/dev/null && echo "Removed any existing container"
