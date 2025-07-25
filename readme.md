# Dev task in Go Lang

## Build up the Docker containers
1. Download the project to a desired location.

    git clone https://github.com/killedit/2025-07-23-dev-task-go-lang

2. Build the containers from the project folder. The task does not require the use of Docker, but have I decided this is the easiest way to ensure my application runs the same way on any platform.

Thinking about it now, I should have chosen SQLite. This is an OK approach for dev. Any writing process is locking the db file even for reading.

    docker compose up -d --build

Please note that if a container is failing to build/run the port might be in use by anther process.

3. Connect to PostgreSQL locally/container in DBeaver

    Host:       localhost/db
    Port:       5432
    Database:   dev_task_go_lang
    Username:   go_user
    Password:   go_password



