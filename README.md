# gitinsight

## Usage


- gen config
```bash
gitinsight config gen
```

- config.yaml
```yaml
debug: false
server:
    address: 0.0.0.0:8080
    database:
        type: sqliteshim
        dsn: file:gitinsight.db
insight:
    auths:
        - domain: github.com
          username: robotism
          password: robotism
    repos:
        - url: https://github.com/robotism/gitinsight.git
          user: robotism
          password: robotism
    authors:
        - name: robotism
          email: robotism@robotism.com
          nickname: robotism
    reset: false
    cache:
        path: ./.repos
    interval: 15m
    since: ""

```

- docker

> https://github.com/robotism/gitinsight/pkgs/container/gitinsight

```bash
# generate config
touch config.yaml
docker run --rm \
-v $(pwd)/config.yaml:/config.yaml \
ghcr.io/robotism/gitinsight:20251007-296c4ca \
/gitinsight config gen -f /config.yaml -o
# vim config.yaml

# run
docker run --rm \
-p 8088:8080 \
-v $(pwd)/config.yaml:/config.yaml \
-v $(pwd)/.repos:/.repos \
ghcr.io/robotism/gitinsight:20251007-296c4ca

```

- docker-compose

```bash


```

## Screenshot

![](screenshots/home.png)
![](screenshots/analyzer1.png)
![](screenshots/analyzer2.png)
![](screenshots/contributors.png)