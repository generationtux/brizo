# Config

### Kubernetes

The easiest way to run Brizo is inside a pod running inside a Kubernetes cluster. Doing so will allow Brizo to use the the default service account to authenticate with Kubernetes and no other configuration is needed.

If you need to run Brizo outside of the cluster and configure access externally you may use a `kubeconfig` file.

| Flag               | Env                    |           |
| ------------------ | ---------------------- | --------- |
| `--k8s-external`   | `BRIZO_K8S_EXTERNAL`   | *optional connect to external cluster, defaults false* |
| `--k8s-config`     | `BRIZO_K8S_CONFIG`     | *optional path to kubeconfig file: if not present, will look for ./kubeconfig or ~/.kube/config (--k8s-external must be used)* |
| `--k8s-context`    | `BRIZO_K8S_CONTEXT`    | *optional context name to use (--k8s-external must be used)* |

### Database

Brizo uses MySQL to store data. You can configure the connection details by passing the following flags to the `run` command, or setting environment variables.

| Flag               | Env                    | Default   |
| ------------------ | ---------------------- | --------- |
| `--mysql-host`     | `BRIZO_MYSQL_HOST`     | localhost |
| `--mysql-port`     | `BRIZO_MYSQL_PORT`     | 3306      |
| `--mysql-user`     | `BRIZO_MYSQL_USER`     | root      |
| `--mysql-password` | `BRIZO_MYSQL_PASSWORD` |           |
| `--mysql-database` | `BRIZO_MYSQL_DATABASE` | brizo     |
