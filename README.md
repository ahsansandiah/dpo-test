# Depoguna Bangunan Online TEST

This project is an initial test to enter Depoguna Bangunan Online as a Backend Engineer

## Folder Structure

- **`api`**: Contains all the API-related code.
  - **`customer`**: Specific to the "customer" domain.
    - **`delivery`**: Handles data delivery, such as API endpoints and controllers.
    - **`domain`**: Contains business entities and models.
    - **`repository`**: Manages data access and database queries.
    - **`usecase`**: Business logic and application rules for the "customer" domain.

- **`cmd`**: Contains the main entry points of the application. Typically, this includes the main application file.

- **`helpers`**: Utility functions and helpers used throughout the application.
  - **`error`**: Error handling utilities.
  - **`trace`**: Tracing utilities for debugging and monitoring.

- **`migrations`**: Database migration scripts for managing schema changes.

- **`packages`**: Various modules and packages used within the application.
  - **`client`**: Client interactions and external service communications.
  - **`config`**: Application configuration settings.
  - **`json`**: Utilities for JSON data handling.
  - **`log`**: Logging utilities.
  - **`manager`**: Manages various aspects of the application, such as database connections.
  - **`server`**: Server-related code for running the application.
  - **`storage`**: Data storage management.

## Installation

### Setup Env
In the `packages > config` folder there is `placholder.env`, please copy it and set the respective values ​​for the existing variables

### Not Using Docker
#### Run application:

`$ make run`

#### Run unit test:

`$ make test`

### Using Docker

Run Docker Image:

`$ docker build -t dpo-test .` 

Run Docker Container:

`$ docker run -p 8080:8080 dpo-test`

Run Docker Host (to use the running container host)

`docker run --network host -d dpo-test`

## Migrate Database
I'am using [Goose](https://github.com/pressly/goose) to track changes on database.

### Usage
1. Set your working directory to `/migrations`
2. Then execute this command to create new migration script `$ goose create changes-description sql`.
3. Above commnad will generate new `.sql` file and now write your queries inside below section.
    ```
    -- +goose Up
    -- +goose StatementBegin
        YOUR SQL Query Should be here!!!
    -- +goose StatementEnd
    ```
    Make sure your queries work like a charm before doing the next steps
4. Run command `make migrate-up-local`
<br />

### Entity Diagrams

https://dbdiagram.io/d/DPO-Diagram-6686aa879939893dae0ef2c8

<img src="https://i.ibb.co.com/b71thP8/Screenshot-2024-07-04-at-21-30-41.png" width="550">

## Copyright


Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by:
* [Ahsan Sandiah](https://www.linkedin.com/in/sansandiah/)