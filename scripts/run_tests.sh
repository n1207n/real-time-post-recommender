docker-compose -f deployments/docker-compose.test.yaml up migrate
docker-compose -f deployments/docker-compose.test.yaml up backend-test --build --remove-orphans
docker-compose -f deployments/docker-compose.test.yaml down