# Dev task in Go Lang
<em>App setup and usage description.</em>

## Build up the Docker containers
1. Download the project to a desired location.

```
git clone https://github.com/killedit/2025-07-23-dev-task-go-lang
```

2. Build the containers from the project folder. Of course `docker` engine and `docker compose` are required. The task does not require the use of Docker, but have I decided this is the easiest way to ensure my application runs the same way on any platform.

```
docker compose up --build
```

Skipping the `-d` flag will not detach and will output the logs in the terminal and there's no need to tail the logs `docker container logs -f --tail 0 dev_task_go_lang-app` or opening Docker Desktop (GUI).

```
dev_task_go_lang-app  | Creating DB.
dev_task_go_lang-app  | Building app binary...
dev_task_go_lang-app  | Runing DB seeder.
dev_task_go_lang-app  | 2025/07/26 13:53:26 Connected to PostgreSQL database
dev_task_go_lang-app  | DB connected: true
dev_task_go_lang-app  | Seeding the DB
dev_task_go_lang-app  | ===
dev_task_go_lang-app  | Quantum chaos: Put(cat, meow) silently failed
dev_task_go_lang-app  | Quantum chaos: Put(dog, bark) stored as (chaos_dog, bark)
dev_task_go_lang-app  | Seeding complete.
dev_task_go_lang-app  | Starting DB.
dev_task_go_lang-app  | 2025/07/26 13:53:26 Connected to PostgreSQL database
dev_task_go_lang-app  | DB connected: true
dev_task_go_lang-app  | DB is running...
dev_task_go_lang-app  | Available commands:
dev_task_go_lang-app  |   -dump     : Dump the DB contents
dev_task_go_lang-app  |   -test     : Run tests to demonstrate quantum chaos
dev_task_go_lang-app  |   -example  : Run example usage
dev_task_go_lang-app  | DB is still alive... (Press Ctrl+C to exit)
```

Please note that if the db container is failing to build/run the port might be in use by another process. Just pick another port and map it to 5432 on this line `{free_port}:5432` in `docker-compose.yaml`.

3. Connect to PostgreSQL in DBeaver to monitor the records.

```
Host:       db
Port:       5432
Database:   dev_task_go_lang
Username:   go_user
Password:   go_password
```

You should see a DB with just one table `dev_task_go_lang > Databases > Schemas > public > Tables > table_key_value`.

In Laravel/PHP db migrations and seeding is done at `docker compose up` step and this is the approach I was after. I have used an `entrypoint.sh` bash script.

4. App usage.

```
docker ps -a
docker exec -it dev_task_go_lang-app bash
go run .
go run . -dump
go run . -test
go run . -example
```

5. Task requirements:
- Put (key, value) -> store data
- Get (key) -> retrieve data
- Delete (key) -> delete data

```
go run . -put-key=test -put-value=test
go run . -get-key=test
go run . -delete-key=test
```

6. This part of the code introduces about 30% chance of any operation misbehaving. There is a test to confirm it later.

```
func quantumChaos() bool {
	return rand.Float32() < 0.3
}
```

Then in each function we have cases that follow the task's requirements:
- Put - 30% it silently fails, stores wrong key or store wrong value.
- Get - 30% of the time will return a random result you did not ask for.
- Delete - 30% of the time will delete a random key instead of one you asked for.

7. Tests are included in `schrodinger_test.go`.

7.1. TestSchrodingerStore_BasicOperations():
- Tests Put, Get and Delete for a single key.
- Checks if value is stored, retrieved, and deleted as expected.
- May fail if there is quantumChaos, when Put, Get, Delete in main.go have `if quantumChaos() {...}`.

7.2. TestSchrodingerStore_QuantumChaos():
- Runs many Get operations to check how ofthen chaos is occuring.

7.3. TestSchrodingerStore_Dump
- Prints out the DB state with Dump().

7.4. TestSchrodingerStore_ConcurrentOperations
- Runs concurent Put and Get operations.

7.5. TestSchrodingerStore_EdgeCases
- Tests empty keys/values, long keys/values, and special characters.

7.6. BenchmarkSchrodingerStore_Put and BenchmarkSchrodingerStore_Get
- Performance of Get and Put.

<em>This command runs all benchmark tests `go test -bench=.`, but instead it can be run as `go test -bench=BenchmarkSchrodingerStore_Put`</em>

```
go test -v
go test -v -cover
go test -bench=.
go test -bench=. -benchmem
go test -v -run TestSchrodingerStore_BasicOperations
go test -v -run TestSchrodingerStore_QuantumChaos
go test -v -run TestSchrodingerStore_Dump
go test -v -run TestSchrodingerStore_ConcurrentOperations
go test -v -run TestSchrodingerStore_EdgeCases
go test -v -run BenchmarkSchrodingerStore_Put
go test -v -run BenchmarkSchrodingerStore_Get

```

Put, Get, Delete randomnly misbehave in `TestSchrodingerStore_BasicOperations()`. Delete(test_key) might delete another key or Get(test_key) might return the value for another key.

If you need the result from Dump() to be sorted by last created record please replace `ORDER BY key` in the query with `ORDER BY created_at DESC`.

8. Exit the app and remove the containers.

```
exit+Enter
docker compose down
```

P.S. Thank you in advance! I will appreciate any feedback.