# Laboratorio 2 #

* Vicente Muñoz Rojas - 202073557-3
* Carlos Lagos - 202073571-9
* Carlos Kuhn - 202073574-3


## Instrucciones 
Para iniciar el servidor OMS usando docker usamos el siguente comando.
```sh
make docker-oms
```
Para iniciar un servidor ONU usando docker usamos el siguiente comando. 
```sh
make docker-onu
```
Ademas para acceder a la terminal del docker usamos.
```sh
docker exec -it grupo17-laboratorio-2_onu_1 sh
```
Para iniciar los servidores de los continentes se usa el siguente comando.
```sh
make docker-continente
```
Para iniciar los servidores de los datanodes se usa el siguente comando.
```sh
make docker-datanode1
make docker-datanode2
```
