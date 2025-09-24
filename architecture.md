Architecture of the Metrics Collection System for Delivery Dashboard
1. Overview

The goal of this project is to create a service in Golang to collect, process, and provide development team performance metrics from Jira and GitHub for subsequent visualization in Grafana.

The service will run as a background application (daemon) that periodically queries the Jira and GitHub APIs, transforms the received data into Prometheus metrics, and exposes them via an HTTP endpoint for Prometheus to scrape.
2. System Components

The system consists of the following key components:

    Data Collectors:

        Jira Collector: A module responsible for connecting to the Jira API. The andygrunwald/go-jira library will be used for API interaction. It will fetch data about issues, sprints, statuses, and time tracking.

        GitHub Collector: A module for connecting to the GitHub API. The official google/go-github client will be used for this purpose. It will gather information on commits, pull requests, reviews, and branch activity.

    Application Core:

        Configuration: Manages connection settings for APIs (tokens, URLs) and data collection parameters (repositories, Jira projects).

        Scheduler: Periodically runs the data collectors (e.g., every 5 minutes).

        Metrics Processor: Aggregates and transforms data from the collectors into the Prometheus metrics format.

    Storage & Export:

        Prometheus Exporter: Provides an HTTP endpoint /metrics that Prometheus will use to scrape data. The prometheus/client_golang library is used for this.

        Internal Cache (optional): For temporarily storing data between collection cycles to avoid frequent API requests.

    Visualization:

        Prometheus: A system for collecting and storing time-series data (metrics).

        Grafana: A platform for visualizing data stored in Prometheus.

3. Data Flow Diagram

   The Scheduler triggers the GitHub Collector and Jira Collector on a schedule.

   The GitHub Collector sends requests to the GitHub API to get data on commits and pull requests.

   The Jira Collector sends requests to the Jira API to get data on issues and sprints.

   The collected raw data is passed to the Metrics Processor.

   The processor converts the raw data into Prometheus metrics. For example:

        github_pull_requests_open_total{repository="...", user="..."}

        jira_issues_in_progress_total{project="...", assignee="..."}

        github_commits_total{repository="...", user="..."}

   The Prometheus Exporter updates the metrics available at the /metrics endpoint.

   The Prometheus Server periodically scrapes the /metrics endpoint and stores the data.

   Grafana connects to Prometheus as a data source and builds dashboards.

4. Project Structure (Go)
```
delivery-dashboard/
├── cmd/
│   └── app/
│       └── main.go         # Application entry point (app.go in our root)
├── pkg/
│   ├── collector/
│   │   └── collector.go    # Data collection orchestrator
│   ├── config/
│   │   └── config.go       # Configuration management
│   ├── github/
│   │   └── github.go       # Client for the GitHub API
│   ├── jira/
│   │   └── jira.go         # Client for the Jira API
│   └── storage/
│       └── storage.go      # Handles Prometheus metrics
├── go.mod
├── go.sum
└── .env                    # Environment variables file
```

5. Key Libraries

   GitHub API Client: google/go-github — The official, well-supported client for working with the GitHub API.

   Jira API Client: andygrunwald/go-jira — A popular and feature-rich library for Jira integration.

   Prometheus Metrics: prometheus/client_golang — The standard library for creating and exporting Prometheus metrics.

   Environment Files: joho/godotenv — For convenient handling of .env files.

6. Metric Definitions

Some of the key metrics to be collected include:

GitHub:

    Number of open/closed/merged Pull Requests (by user, repository).

    Pull Request lifetime (from opening to merging).

    Number of commits (by user, repository).

    Number of comments in Pull Requests.

Jira:

    Number of issues in different statuses (by project, sprint, assignee).

    Time an issue spends in a status.

    Number of created/completed issues per sprint.

    Lead time / Cycle time.

This architecture ensures modularity and extensibility, allowing for the future addition of new data sources (e.g., GitLab) or new metrics without significant changes to the core system.