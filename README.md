# go eats server
The server for Go Eats
REST API

Users can submit food places to the server pending approval from admin
Admin can see pending approvals. Edit and approve accordingly

docker build -t go-eats-server .
docker run -it --rm -p 8080:8080 go-eats-server