# Load Tests


## gRPC Units Load Tester
A utility for load testing gRPC services with flexible configuration for scenarios, 
concurrency, duration, and report formats.

### ⚙️ Features
Run one or multiple scenarios concurrently

- Test by number of requests or duration
- Control concurrency levels
- Support for limiting requests per second (RPS)
- Multiple report formats: print, json, html

All reports are saved to the ./reports directory

### 🚀 Usage
```bash
go run main.go [flags]
```

Examples:
1. Simple run by number of requests:
   ```bash
    go run main.go -addr=localhost:50051 -scenarios=users_auth_register -requests=100 -concurrency=10
   ```
2. Run by duration (without RPS limitation):
   ```bash
    go run main.go -addr=localhost:50051 -scenarios=users_auth_register -duration=10s -concurrency=20
   ```
3. Run with fixed RPS:
   ```bash
    go run main.go -addr=localhost:50051 -scenarios=users_auth_register -duration=10s -rpc=100 -concurrency=10
   ```
4. Multiple scenarios in a single run:
   ```bash
    go run main.go -addr=localhost:50051 -scenarios="users_auth_register,users_auth_login" -requests=500 -concurrency=20
   ```
5. HTML report:
   ```bash
    go run main.go -addr=localhost:50051 -scenarios=users_auth_register -requests=1000 -report=html
   ```

### 🧾 Supported Flags
   Flag	Description
   - ```-addr```	gRPC server address (default: localhost:50051)
   - ```scenarios```	Comma-separated list of scenario keys (registered in the code)
   - ```concurrency```	Number of concurrent workers per scenario (default: 10)
   - ```requests```	Total number of requests (if specified, duration is ignored)
   - ```duration```	Time duration to run the test (e.g., 10s)
   - ```rpc```	Requests per second (for duration mode)
   - ```timeout```	Timeout for each request (default: 3s)
   - ```report```	Report format: print, json, html (default: print)

### 📁 Reports
All reports are saved in the ./reports directory:

- report_<scenario>.json — structured result
- report_<scenario>.html — a nicely formatted visual report (if -report=html is selected)

### 🧱 Scenario Structure
Scenarios are registered in the code and may contain:

- Name
- Main function (Call(ctx))
- An optional Init() function (for data setup)
- Dependency on another scenario (DependsOn)
