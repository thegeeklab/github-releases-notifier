# github-releases-notifier

GitHub release notification bot

[![Build Status](https://img.shields.io/drone/build/thegeeklab/github-releases-notifier?logo=drone)](https://cloud.drone.io/thegeeklab/github-releases-notifier)
[![Docker Hub](https://img.shields.io/badge/dockerhub-latest-blue.svg?logo=docker&logoColor=white)](https://hub.docker.com/r/thegeeklab/github-releases-notifier)
[![Quay.io](https://img.shields.io/badge/quay-latest-blue.svg?logo=docker&logoColor=white)](https://quay.io/repository/thegeeklab/github-releases-notifier)
[![Go Report Card](https://goreportcard.com/badge/github.com/thegeeklab/github-releases-notifier)](https://goreportcard.com/report/github.com/thegeeklab/github-releases-notifier)
[![Source: GitHub](https://img.shields.io/badge/source-github-blue.svg?logo=github&logoColor=white)](https://github.com/thegeeklab/github-releases-notifier)
[![License: MIT](https://img.shields.io/github/license/thegeeklab/github-releases-notifier)](<[LICENSE](https://github.com/thegeeklab/github-releases-notifier/blob/master/LICENSE)>)

Receive Slack notifications if a new release of your favorite software is available on GitHub.

## Setup

1. Get a URL to send WebHooks to your Slack.
2. Get a token for scraping GitHub: [https://help.github.com/](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token).

To watch repositories simply add them to the list of arguments e.g. `-r=kubernetes/kubernetes -r=prometheus/prometheus`.

### Docker

```Shell
docker run --rm \
    -e GITHUB_TOKEN=XXX \
    -e SLACK_HOOK=https://hooks.slack.com/... \
    thegeeklab/github-releases-notifier -r=kubernetes/kubernetes
```
