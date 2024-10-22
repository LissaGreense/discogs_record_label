# Discogs Label Releases

An application to search through artist, style or genre of releases for selected record label.

## Starting the Application
### Using Docker Compose:

1. Environment Configuration
To run the application, you need to set up two environment configuration files: *.env.backend* and 
*.env.db*. These files contain the necessary variables for the backend service and the database 
connection, respectively.

### .env.backend
Create a file named .env.backend in the root of your project with the following configuration:

```
dotenv
CORS_ORIGIN=http://localhost:3000               # The allowed origin for CORS requests
POSTGRES_USER=youruser                          # Your PostgreSQL username
POSTGRES_PASSWORD=yourpassword                  # Your PostgreSQL password
POSTGRES_DB=yourdb                              # The name of your PostgreSQL database
POSTGRES_HOST=db                                # The host where PostgreSQL is running (usually `db` in Docker setup)
POSTGRES_PORT=5432                              # The port for PostgreSQL (default is 5432)
DISCOGS_APP_NAME=PROVIDE_APP_NAME               # The name of your application (used in User-Agent)
DISCOGS_KEY=API_KEY                             # Your Discogs API key
DISCOGS_SECRET=API_SECRET                       # Your Discogs API secret
SELECTED_LABEL=5                                # Label id to fetch from Discogs API
```

### .env.db
Create a separate file named .env.db in the root of your project with the following configuration:

```
dotenv
POSTGRES_USER=youruser                     # Your PostgreSQL username
POSTGRES_PASSWORD=yourpassword             # Your PostgreSQL password
POSTGRES_DB=yourdb                         # The name of your PostgreSQL database
POSTGRES_HOST=db                           # The host where PostgreSQL is running (usually `db` in Docker setup)
POSTGRES_PORT=5432                         # The port for PostgreSQL (default is 5432)
```

2. Run docker-compose
Run the following command to build and start the services:

```
bash
docker-compose up --build
```
This will start your backend service and Postgres database. The React frontend will be served through Nginx.


3. Access the Application:

Open your browser and go to http://localhost:3000 to view the application

## Backend Tests
Navigate to your Go backend directory:

```
bash
cd backend
```
Run the following command to execute tests:
```
bash
go test ./...
```

## Frontend Tests

Navigate to your Go backend directory:

```
bash
cd tone_addiction_frontend
```
Run the following command to execute tests:
```
bash
npm test
```
