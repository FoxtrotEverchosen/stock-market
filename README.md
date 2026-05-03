# Simplified Stock Market Implementation

This project implements a simple stock market represented as single entity bank, wallets and audit logs. It allows for creating a buy/sell actions by wallets and creating stocks by bank. Project utilizes Go for core implementation, Docker, Nginx and Redis for persistence.

## Architecture
This implementation consists of four main components running as docker containers.

- Go application - two instances of simple HTTP server responsible for handling API requests
- Redis - data storage, mutable state shared between app instances 
- nginx - load balancer diverting traffic to free instance of application

Two instances of the application handle the same requests and store data in shared redis database. Any of them can handle any request. This, combined with nginx load balancing allows for high availability. If one instance cannot handle the request, it will be passed to another one. Additionally, any time the app instance is killed with POST /chaos, docker will try to restart it again.

To avoid possible race conditions by two simultaneous buy/sale operations, both of them are handled on Redis side with Lua scripts. Since these are executed atomically on Redis, the shared state is always modified only by one app instance.

## Running the Application
Application requires docker to run. Depending on the operating system, the application can be run with one command.

On Windows:
```shell 
./start.bat 8080
 ```

On Linux/MacOS (you may need to make the script executable):
```bash 
chmod +x start.sh
./start.sh 8080
 ```

## Endpoint
You can run API calls to this app following the provided structure.

| Method| Path                                      | Description                                   |
|:------|:------------------------------------------|:--------------                                |
| POST  | /stocks                                   | Set bank stock levels                         |
| POST  | /wallets/{wallet_id}/stocks/{stock_name}  | Buy or sell a single stock                    |
| POST  | /chaos                                    | Kill the running instance                     |
| GET   | /stocks                                   | Get current bank stock levels                 |
| GET   | /wallets/{wallet_id}                      | Get wallet state                              |
| GET   | /wallets/{wallet_id}/stocks/{stock_name}  | Get quantity of a specific stock in a wallet  |
| GET   | /log                                      | Get full audit log                            |

### Examples:
```bash
# Set bank stocks
curl -X POST http://localhost:8080/stocks \
  -H "Content-Type: application/json" \
  -d '{"stocks": [{"name": "AAPL", "quantity": 10}]}'

# Get bank stocks
curl http://localhost:8080/stocks

# Buy a stock
curl -X POST http://localhost:8080/wallets/wallet1/stocks/AAPL \
  -H "Content-Type: application/json" \
  -d '{"type": "buy"}'

# Sell a stock
curl -X POST http://localhost:8080/wallets/wallet1/stocks/AAPL \
  -H "Content-Type: application/json" \
  -d '{"type": "sell"}'

# Get wallet state
curl http://localhost:8080/wallets/wallet1

# Get specific stock quantity in wallet
curl http://localhost:8080/wallets/wallet1/stocks/AAPL

# Get audit log
curl http://localhost:8080/log

# Chaos
curl -X POST http://localhost:8080/chaos
```

## Design Decisions
- Redis as data storage - data model matches Redis relatively well, using hashes for stocks, list for audit log, set for wallet tracking. Additionally since Redis works in memory, all actions are executed fast. Similar could be achieved with other kinds of databases, but resulting in more overhead efficiency loss.

- Lua scripts for buy/sell atomicity - with two app instances running concurrently, a read-then-write approach creates a race condition where both instances could read the same bank quantity and both modify it. Redis executes Lua scripts atomically, so the entire check-and-update sequence is guaranteed to run on its own.

- Vendored dependencies - dependencies copied into the repo. This makes the Docker build fully self-contained. It requires no internet during build, guaranteeing reproducibility on any machine.

- Stateless app instances - the app itself holds no state, everything lives in Redis. This makes scaling easy - adding a third instance requires one line in docker-compose.

## Known Limitations
- Redis persistence - by default Redis is in-memory only. If the Redis container crashes, all data is lost. Production would require configuring Redis persistence (RDB snapshots or AOF logging).

- No authentication - endpoints are open. Production would require at minimum some kind of API key validation.

- Single Redis instance - Redis itself is a single point of failure in this setup. Production would use Redis Sentinel or Redis Cluster for high availability at the data layer too.

- The `proxy_next_upstream` in nginx retries failed requests on the other instance. If a buy/sell Lua script completes successfully but the response doesn't reach nginx before the instance dies, nginx could retry the request on the other instance, executing the trade twice. This is an inherent trade-off between availability and exactly-once delivery.