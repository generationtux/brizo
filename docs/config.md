# Config

### Database

Brizo uses MySQL to store data. You can configure the connection details by passing the following flags to the `run` command, or setting environment variables.

| Flag               | Env                    | Default   |
| ------------------ | ---------------------- | --------- |
| `--mysql-host`     | `BRIZO_MYSQL_HOST`     | localhost |
| `--mysql-port`     | `BRIZO_MYSQL_PORT`     | 3306      |
| `--mysql-user`     | `BRIZO_MYSQL_USER`     | root      |
| `--mysql-password` | `BRIZO_MYSQL_PASSWORD` |           |
| `--mysql-database` | `BRIZO_MYSQL_DATABASE` | brizo     |
