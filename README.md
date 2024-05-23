# locsearch
## This is a capstone project for the Golang course

### Description
This is a simple system that processes user locations and provides the ability to search for clients by location (coordinates) and radius, with an additional functionality of calculating the distance traveled by a user in a given time range.

### Installation
1. Clone the repository
```bash
git clone github.com/surtexx/locsearch
```

2. Install the dependencies
```bash
go mod download
```

3. Create a `.env` file in the root of the project and add the following environment variables:
```bash
LOCATIONS_TABLE=<YOUR_DYNAMODB_LOCATION_TABLE>
LOCATION_HISTORY_TABLE=<YOUR_DYNAMODB_LOCATION_HISTORY_TABLE>
```

4. Source the `.env` file
```bash
source .env
```

5. Save AWS credentials in the `~/.aws/credentials` file
```bash
[default]
aws_access_key_id = <YOUR_AWS_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_AWS_SECRET_ACCESS_KEY>
```

6. (Optional) Save AWS configuration in the `~/.aws/config` file
```bash
[default]
region = <YOUR_AWS_REGION>
output = json
```

7. Run the application
```bash
go run cmd/locsearch/main.go
```

### Testing
1. Unit testing
The application has unit tests saved in the same folder with the components they test. To run the tests, execute the following command:
```bash
go test ./...
```

2. Integration testing
The application has integration tests saved in same folder with the components they test. To run the tests, execute the following command:
```bash
go test ./... -tags=integration
```

### Usage
The application has a REST API that can be accessed via the following endpoints:
1. `PUT /updateLocation` - Update a user location
2. `GET /searchUsers` - Get all users at a given location within a given radius
3. `GET /getDistanceTraveled` - Get total distance travelled by a user in a time range

### Example requests
1. Update a user location
```bash
curl -X PUT \
http://localhost:8080/updateLocation\?username\=rgheorghe\&newLocation\=39.13355,27.14538
```

2. Search for users at a given location within a given radius
```bash
curl -X GET \
http://localhost:8080/searchUsers\?coordinates\=39.13355,29.14538\&radius\=100
```

3. Get total distance traveled by a user in a time range
```bash
curl -X GET \
http://localhost:8080/getDistanceTraveled\?username\=rgheorghe\&startDate\=2023-01-01T00:00:00+00:00\&endDate\=2024-12-30T00:00:00+00:00
```

### Using Docker
1. Build the Docker image
```bash
docker build -t locsearch --build-arg LOCATIONS_TABLE=<YOUR_LOCATIONS_TABLE> --build-arg LOCATION_HISTORY_TABLE=<YOUR_LOCATION_HISTORY_TABLE> .
```

2. Run the Docker container
```bash
docker run -p 8080:8080 locsearch
```