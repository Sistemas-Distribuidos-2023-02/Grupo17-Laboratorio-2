docker-onu:
    docker-compose -f docker-compose.yml up onu

docker-continente:
    docker-compose -f docker-compose.yml up continente

docker-oms:
    docker-compose -f docker-compose.yml up oms

docker-datanode:
	docker-compose -f docker-compose.yml up datanode

docker-down:
    docker-compose -f docker-compose.yml down

docker-clean:
    docker system prune -a