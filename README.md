## Running The Project

To run the project use the following commands.

```
# Install Tooling
$ make dev-gotooling
$ make dev-brew
$ make dev-docker

# Run Tests
$ make test

# Run Benchmarking:
$ make benchmark

# Shutdown Tests
$ make test-down

# Run Project with k8s
$ make dev-up
$ make dev-update-apply
$ make token
$ export TOKEN=<COPY TOKEN>
$ make users

# Run Project with compose
$ make compose-up


# Run Load
$ make load

# Run Tooling
$ make grafana
$ make statsviz

# Shut Project k8s
$ make dev-down

# Shut Project compose
$ make compose-down
```
