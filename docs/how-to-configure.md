# Table of Content

- [Table of Content](#table-of-content)
  - [How to configure](#how-to-configure)
    - [Generic placeholders](#generic-placeholders)
    - [Use environment variables in the configuration](#use-environment-variables-in-the-configuration)
    - [Configuration parameters](#configuration-parameters)
      - [positionConfig](#positionconfig)
      - [categories](#categories)
      - [feeds\_config](#feeds_config)
- [Example config](#example-config)

## How to configure

### Generic placeholders

`<boolean>`: a boolean that can take the values true or false
`<int>`: any integer matching the regular expression [1-9]+[0-9]*
`<string>`: a string
`<url>`: a URL
`<filepath>`: a string containing an absolute or relative path and filename to a file on disk

### Use environment variables in the configuration

You can use environment variable references in the YAML configuration file to set values that need to be configurable during deployment.

Each variable reference is replaced at startup by the value of the environment variable. The replacement is case-sensitive and occurs before the YAML file is parsed. References to undefined variables are replaced by empty strings unless you specify a default value or custom error text.

To specify a default value, use `${VAR}`.

### Configuration parameters

#### positionConfig

```yaml
# Backend storage to use. Supported backends are: filesystem, sqlite
positionConfig:
  backend: <string>

  # Define path of position.yaml file
  filesystem:
    path: <filepath>

  # Define path of position.yaml file
  sqlite:
    path: <filepath>
    database: <string>
```

#### categories

```yaml
# Define category with delivery channel. Supported channel telegram, prometheus, slack
# Replace category_name with you category
categories:
  [<string>]:
    telegram:
      enabled: <boolean>
      chatId: <string>
      botToken: <string>
    prometheus:
      enabled: <boolean>
      url: <url>
      baicAuth:
        username: <string>
        password: <string>
    slack:
      enabled: <boolean>
      webhookUrl: <string>
    feeds: [list of key value> | default = []]
```

#### feeds_config

```yaml
name: <string>
url: <url>
```

# Example config

```yaml
positionConfig:
  backend: filesystem

  filesystem:
    path: position.yaml

categories:
  technology:
    telegram:
      enabled: true
      chatId: ${TELEGRAM_TECHNOLOGY_CHAT_ID}
      botToken: ${TELEGRAM_BOT_TOKEN}
    prometheus:
      enabled: true
      url: ${GRAFANA_CLOUD_METRIC_ENDPOINT}
      baicAuth:
        username: ${GRAFANA_CLOUD_METRIC_USER}
        password: ${GRAFANA_CLOUD_METRIC_TOKEN}
    slack:
      enabled: true
      webhookUrl: ${SLACK_WEBHOOK_URL}
    feeds:
      - name: Kubernetes Blog
        url: https://kubernetes.io/feed.xml
      - name: Grafana Blog
        url: https://grafana.com/blog/news.xml
      - name: Prometheus Blog
        url: https://prometheus.io/blog/feed.xml
      - name: Substack Kubernetes Weekly Blog
        url: https://learnkubernetesweekly.substack.com/feed
```
