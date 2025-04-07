# StockOverflow

### Usage For Compose

```bash
### simply compose up would trigger the load test
### this would automatically open the loadtest with the server
sudo docker-compose up 

### However We provide a server only entry
sudo docker-compose up app
```

### Usage for Test

```bash
### inside docker-deploy/testing/cores_*_test.sh
chmod 777 *.sh
### To reproduce the test results We need a 8 cores vm
./cores_*_test.sh
### forward vm port to connected machine's 8089 port
### Or you can view on http://vcm-xxxxx.vm.duke.edu:8089/
### So you can do browser views for graphs
### You can adjust the max user and user spawn rate according to the website