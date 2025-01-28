# DeBeAndo Agent for MySQL

Database monitoring tool designed for small environments, adapted for Kubernetes and send metrics to InfluxDB.

## Image Description

This image is maintained by DeBeAndo and will be updated regularly on best-effort basis. The image is based on Alpine Linux and only contains the build result of this repository.

## Run

To run container:

```bash
docker run \
	--name debeando-agent-mysql \
	--env DEBUG=true \
	--env INTERVAL=10 \
	--env INFLUXDB_HOST="http://com-env-influxdb-observability-node01.aws.com" \
	--env INFLUXDB_TOKEN="abc123cde456==" \
	--env MYSQL_HOST="com-env-mysql-stack-node01.aws.com" \
	--env MYSQL_USER="monitor" \
	--env MYSQL_PASSWORD="<monitor-pass>" \
	--env SERVER="com-env-mysql-stack-node01" \
	debeando/agent-mysql
```

## MySQL Config

Create a `monitor` user to allow access to agent.

```sql
CREATE USER monitor@'%' IDENTIFIED BY '<monitor-pass>';
ALTER USER monitor@'%' WITH MAX_USER_CONNECTIONS 5;
GRANT REPLICATION CLIENT ON *.* TO monitor@'%';
GRANT PROCESS ON *.* TO monitor@'%';
GRANT SELECT ON *.* TO monitor@'%';
```

Please, change default password `<monitor-pass>`.

## Environment Variables

When you start the `agent-mysql` image, you can adjust the configuration of the agent instance by passing one or more environment variables on the docker run command line.

- **DEBUG:** Enable debug mode with `true` value, by default value is `false`.
- **DISABLE:** Disable specific metric, list separated by comma, by default is empty.
- **INTERVAL:** Interval time in second, by default value is `10`.
- **INFLUXDB_HOST:** The HTTP hostname or IP address.
- **INFLUXDB_TOKEN:** The authentication token for connecting to the InfluxDB instance.
- **MYSQL_HOST:** The hostname or IP address of the MySQL server to be monitored.
- **MYSQL_USER:** The MySQL username to connect with.
- **MYSQL_PASSWORD:** The password associated with the specified MySQL user.
- **SERVER:** The name of running the agent.
