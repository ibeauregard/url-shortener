build:
	docker build --tag=url_shortener:latest .
	docker image prune -f

run:
	docker run --interactive --tty --name=url_shortener_container \
    		--volume=url_shortener_db:/home/urlshortener/db/data --publish=8080:8080 --rm \
    		--env APP_HOST=localhost url_shortener:latest

restart: stop run

clear-db:
	docker exec -it url_shortener_container sh db/clear-db.sh

delete-db:
	@docker volume rm -f url_shortener_db &>/dev/null && echo "Deleted any existing database"

connect:
	docker exec -it url_shortener_container sh

connect-to-db:
	docker exec -it url_shortener_container sqlite3 db/data/url-mappings.db

stop:
	@docker rm -f url_shortener_container &>/dev/null && echo "Stopped any existing container"
