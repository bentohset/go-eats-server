# Go Eats Web Server
This repo houses the web server for Go Eats. The main repo can be found [here](https://github.com/bentohset/go-eats).

The web server acts as the REST API Gateway between clients and the PostgreSQL Database. It provides methods of manipulating the data in the database.

Furthermore, it houses the main infrastructure of the microservice on the Azure ecosystem using Terraform.

## Table of Contents
- [Setup](#setup)
- [Requirements](#requirements)
- [Usage](#usage)
    - [API Endpoints](#api-endpoints)
    - [Modelling](#modelling)
- [Testing](#testing)
- [Deployment](#deployment)
- [Todo](#todo)

## Setup
**Setup go environment:**
```
go mod download
```

**Setup .env variables:**
```
DB_HOST=<db uri>
DB_USERNAME=postgres
DB_PASSWORD=
DB_NAME=eats
TEST_DB_NAME=eats_test
DOCKER_PASSWORD=
```

**Run the server:**
```
go run .
```
or 
```
go build
./go-eats-server
```

## Requirements
PostgreSQL
Gorilla MUX API - acts as the router in this web server

## Usage
### API endpoints

**GET /health**
- retrieves the health of the server
- getHealth()

It should return
```
{
    "result": "Server is up and running"
}, 200
```

**GET /places**
- returns all the places in the DB as an array of objects
- getPlaces()

It should return
```
[
    {
        "id":
        "name":
        "budget":
        "location":
        "mood":
        "cuisine":
        "mealtime":
        "rating":
        "approved":
    },
    ...
], 200
```

**POST /places**
- creates a new place
- createPlace()

It accepts:
```
{
        "name":
        "budget":
        "location":
        "mood":
        "cuisine":
        "mealtime":
        "rating":
}
```

It should return:
```
{
        "name":
        "budget":
        "location":
        "mood":
        "cuisine":
        "mealtime":
        "rating":
}, 201
```

**GET /places/approved**
- returns all the approved places with approved field as `true` as an array of objects
- getApprovedPlaces()

It should return
```
[
    {
        "id":
        "name":
        "budget":
        "location":
        "mood":
        "cuisine":
        "mealtime":
        "rating":
        "approved": true
    },
    ...
], 200
```

**GET /places/requested**
- returns all the requested places with approved field as `false` as an array of objects
- getRequestedPlaces()

It should return
```
[
    {
        "id":
        "name":
        "budget":
        "location":
        "mood":
        "cuisine":
        "mealtime":
        "rating":
        "approved": false
    },
    ...
], 200
```

**GET /places/{id}**
- returns an individual places for the given `id`
- getPlace()

It should return
```
{
    "id": {id}
    "name":
    "budget":
    "location":
    "mood":
    "cuisine":
    "mealtime":
    "rating":
    "approved":
}, 200
```

**PUT /places/{id}**
- updates a place with a new place for the given `id`
- updatePlace()

It accepts
```
{
    "id": {id}
    "name":
    "budget":
    "location":
    "mood":
    "cuisine":
    "mealtime":
    "rating":
    "approved":
}
```
It should return
```
{
    "id": {id}
    "name":
    "budget":
    "location":
    "mood":
    "cuisine":
    "mealtime":
    "rating":
    "approved":
}, 200
```


**DELETE /places/{id}**
- deletes a place row for the given `id`
- deletePlace()

It should return
```
{
    "result":"success"
}, 200
```


**PATCH /places/{id}/approve**
- updates the approved field of the row to `true` for the given `id`
- approvePlace()

It should return
```
{
    "id": {id}
    "name":
    "budget":
    "location":
    "mood":
    "cuisine":
    "mealtime":
    "rating":
    "approved": true
}, 200
```

**PATCH /places/{id}/disapprove**
- updates the approved field of the row to `false` for the given `id`
- disapprovePlace()

It should return
```
{
    "id": {id}
    "name":
    "budget":
    "location":
    "mood":
    "cuisine":
    "mealtime":
    "rating":
    "approved": false
}, 200
```

### Modelling
`places` table to store places:
| Field    | Data type   | Remarks                    |
|----------|-------------|----------------------------|
| ID       | serial/int  | primary key                |
| name     | text/string |                            |
| budget   | numeric/int |                            |
| location | text/string |                            |
| mood     | text/string | string separated by commas |
| cuisine  | text/string | string separated by commas |
| mealtime | text/string | string separated by commas |
| rating   | numeric/int | 0 to 5                     |
| approved | boolean     | default false              |

`reviews` table to store reviews of restaurants:
| Field    | Data type   | Remarks                    |
|----------|-------------|----------------------------|
| ID       | serial/int  | primary key                |
| name     | text/string |                            |
| review   | text/string |                            |

The same name might have multiple entries, this is resolved on the recommendation side where data is processed first before loading the training.

## Testing
Tests are inside the `main_test.go` file
- Each endpoint consists of at least 1 test

To run tests:
```
go test
```




## Deployment
Deploying changes to Docker image:
```
docker-compose build
docker-compose push
```

Deploying as services to Kubernetes:
```
kubectl apply -f ingress.yaml
kubectl apply -f server.yaml
```

It utilises AKS HTTP Application Routing which is Azure's community nginx based ingress controller

go-eats-server service is deployed as a ClusterIP on port 81

The ingress creates a Ingress Controller to route any requests made to the DNS to our service on port 81

DNS is provided with Azure DNS Zone and hosted on the following DNS:
```
go-eats.816b0060473d4e7c99bf.southeastasia.aksapp.io
```


## Todo
- Settle TLS with cert-manager
- Store reviews