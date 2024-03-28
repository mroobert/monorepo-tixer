
start-api:
	go run github.com/mroobert/monorepo-tixer/cmd 

start-infra:
	docker compose up -d

stop-infra:
	docker compose down