# SonarQube / SonarCloud Integration Guide

This repository is configured to automatically analyze code with SonarQube/SonarCloud on every push.

## 1. Prerequisites

You need a running SonarQube instance (e.g., on your EC2, Docker, or SonarCloud.io).

### Docker Helper (Run Locally/EC2)
```bash
docker run -d --name sonarqube -e SONAR_ES_BOOTSTRAP_CHECKS_DISABLE=true -p 9000:9000 sonarqube:latest
```
*Login with admin/admin and change password.*

## 2. GitHub Secrets Setup

Go to **Settings** -> **Secrets and variables** -> **Actions**.

Add the following secrets:

| Secret Name | Description |
|---|---|
| `SONAR_HOST_URL` | Your SonarQube server URL. <br>Example: `http://3.106.254.198:9000` or `https://sonarcloud.io` |
| `SONAR_TOKEN` | A global analysis token (User Token) generated in SonarQube. <br>Go to **User > My Account > Security > Generate Token**. |

## 3. Project Configuration

The workflows use these project keys by default. You should create these projects in SonarQube manually if using the Community Edition, or let them be auto-created if using Enterprise/Cloud with provisioning.

-   `eticketing-backend-dotnet-tickets`
-   `eticketing-backend-go-payment`
-   `eticketing-backend-python-ledger`
-   `eticketing-frontend-react`
-   `eticketing-frontend-backoffice`

## 4. How it works

-   **Backend .NET**: Uses `dotnet-sonarscanner`. It installs the tool, begins analysis, attempts to build the project, and ends analysis to upload results.
-   **Others (Go, Python, React)**: Uses the generic `sonarsource/sonarqube-scan-action`. It reads the `sonar-project.properties` file in each directory and uploads the source code for analysis.

## 5. View Results

After a GitHub Action run completes, check your `SONAR_HOST_URL` dashboard to see bugs, vulnerabilities, and code smells.
