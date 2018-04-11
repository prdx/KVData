# About

Very naive and simple implementation of distributed key value database in GO. Created in purpose of learning Golang and system design.

# How to run

1. Use `./run_10_servers.sh` and then run `./run_test.sh`
2. Use make all and then `./server -ip <IP> -port <PORT>` to run the server and `./proxy <SERVER_IP>:<SERVER_PORT>`

# Load balancing design

Data assigned at random server and then stored in a struct called `addressBook` on the proxy.

# TODO

1. Reduce redundancy in the code
2. Create unit testing
3. Implement replication
