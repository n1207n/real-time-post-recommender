# real-time-post-recommender
Real-time post recommender based on old hackernews ranking algo

## TODO
- [ ] Gravity adjustment based on ranking stats
- [ ] Simulator based on read-write ratio?
- [ ] Code & Test refactorings

## Prerequisites
- Go 1.20
- Docker & docker-compose
- golang-migrate

## Setup
- Create `.env` in `/deployments` folder

## Commands
### Run docker-compose
`docker-compose -f deployments/docker-compose.yaml up`

### Stop docker-compose
`docker-compose -f deployments/docker-compose.yaml stop`

### Get Postgres shell
`docker-compose -f deployments/docker-compose.yaml exec -it db psql -U app -d app`

### Run tests
`./scripts/run_tests.sh`

### Generate Fake Posts
`go run cmd/driver/post_generator_driver.go`

### DB Migrations
#### Forward migration
`docker-compose -f deployments/docker-compose.yaml --profile tools run migrate up`
#### Backward migration
`docker-compose -f deployments/docker-compose.yaml --profile tools run migrate down`
#### Fix migration
`docker-compose -f deployments/docker-compose.yaml --profile tools run migrate force <VERSION>`