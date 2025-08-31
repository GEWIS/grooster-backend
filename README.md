# grooster-backend
A GEWIS rooster (roster/schedule) maker. This is the backend repository.

## Setup
Make sure you can run go and have this repository open in your favourite IDE.

Next run ```go mod tidy``` to install the required packages.

If you want to seed the database run ```go run ./cmd/seeder/main.go```

To start the application run ```go run ./cmd/src/main.go```

## Generating Docs and Client
To generate the docs go into the cmd/src folder and run the command: ```swag init --pd --parseInternal```

To generate the client which will be uploaded to npm, first go into the cmd/src/docs folder, then run the command: 
```openapi-generator-cli generate -i swagger.yaml -g typescript-axios -o ../../../client/src```
