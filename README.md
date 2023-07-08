# Auctioning App

An auction service for AD space using supply/demand side platform. This app
consists of 4 microservices.

- **auction_db_svc:** Specififes the `ad_space_auction` schema for MySQL
  database and allows saving/retreiving of data by other auction services

  - Table `space` : id, name, base_price
  - Table `bidder` : id, name,
  - Table `auction` : id, space_id, starts_on, ends_on
  - Table `auction_bid` : id, auction_id, bidder_id, bid_price, bid_on

- **auction_supply_svc:** Provides endpoints to perform CRUD operations for AD
  spaces

  - `POST /space `: Create a new space for auction
  - `GET /space `: List all the spaces available for auction
  - `GET /space/:id `: Get the info of a specific space
  - `PATCH /space/:id `: Update the info of a specific space
  - `DELETE /space/:id `: Delete a specific space

- **auction_demand_svc:** Provides endpoints to perform CRUD operations for
  potential bidders

  - `POST /bidder `: Create a new bidder
  - `GET /bidder `: List all the bidders
  - `GET /bidder/:id `: Get the info of a specific bidder
  - `PATCH /bidder/:id `: Update the info of a specific bidder
  - `DELETE /bidder/:id `: Delete a specific bidder

- **auction_main_svc:** Provides endpoints to create/get an auction and allows
  to make a bid in the auction
  - `POST /auction `: Create a new auction
  - `GET /auction `: List all the auctions
  - `GET /auction/:id `: Get the info of a specific auction
  - `POST /auction/:id/bid `: Make a bid in a specific auction

## Installation

1. Build the docker images

```
make build
```

OR

```
docker build -t apsgangwar/auction_db_svc:latest ./db/
docker build -t apsgangwar/auction_supply_svc:latest ./supply/
docker build -t apsgangwar/auction_demand_svc:latest ./demand/
docker build -t apsgangwar/auction_main_svc:latest ./auction/
```

2.  Start the services

It will take about a minute to setup database and start the services.

```
make install
```

OR

```
docker-compose up -d
```

3. Stop the services

```
make uninstall
```

OR

```
docker-compose down
```

## Testing

Currently, tests are added for `supply/svc` package in `supply/svc/svc_test.go`
file. To test the code, run these command:

```
cd supply
go mod download
go test supply/svc
```
