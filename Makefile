docker-onu:
	docker-compose -f docker-compose.yml up -d onu

docker-continente:
	docker-compose -f docker-compose.yml up continente

docker-oms:
	docker-compose -f docker-compose.yml up oms

docker-datanode1:
	docker-compose -f docker-compose.yml up datanode1

docker-datanode2:
	docker-compose -f docker-compose.yml up datanode2

docker-down:
	docker-compose -f docker-compose.yml down

docker-clean:
	docker system prune -a
