.PHONY: build install uninstall push all

build: build_db_svc build_supply_svc build_demand_svc build_auction_svc

build_db_svc:
	docker build -t apsgangwar/auction_db_svc:latest ./db/

build_supply_svc:
	docker build -t apsgangwar/auction_supply_svc:latest ./supply/

build_demand_svc:
	docker build -t apsgangwar/auction_demand_svc:latest ./demand/

build_auction_svc:
	docker build -t apsgangwar/auction_main_svc:latest ./auction/

install:
	docker-compose up -d

uninstall:
	docker compose down

push:
	docker push apsgangwar/auction_db_svc:latest
	docker push apsgangwar/auction_supply_svc:latest
	docker push apsgangwar/auction_demand_svc:latest
	docker push apsgangwar/auction_main_svc:latest

all: build install