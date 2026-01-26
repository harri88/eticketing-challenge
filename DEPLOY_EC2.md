# Deployment Guide (EC2)

## 1. Prepare your EC2 Instance

1.  **SSH into your EC2 instance**:
    ```bash
    ssh -i your-key.pem ubuntu@3.106.254.198
    ```

2.  **Install Docker & Docker Compose**:
    (You can use the provided script)
    ```bash
    curl -o setup_ec2.sh https://raw.githubusercontent.com/harri88/eticketing-challenge/main/scripts/setup_ec2.sh
    chmod +x setup_ec2.sh
    ./setup_ec2.sh
    ```
    *Log out and log back in to apply group permissions.*

## 2. Configure GitHub Secrets

Go to your repository **Settings** -> **Secrets and variables** -> **Actions** -> **New repository secret**.

Add the following secrets:

| Secret Name | Value Example | Description |
|---|---|---|
| `EC2_HOST` | `3.106.254.198` | The public IP of your EC2 |
| `EC2_USER` | `ubuntu` | The SSH username (default for Ubuntu AMIs) |
| `EC2_SSH_KEY` | `-----BEGIN RSA PRIVATE KEY----- ...` | Your private SSH key content (open .pem file and copy all) |

## 3. Deployment

Simply push to the `main` branch. The GitHub Action will:
1.  Package your code (excluding large folders like node_modules).
2.  Copy it to the server.
3.  Generate the correct production `.env` file on the server.
4.  Run `docker-compose up -d --build` to rebuild and start services.

## Troubleshooting

-   **Frontend API calls failing**: Check if `REACT_APP_API_URL` is set correctly in the `.env` on the server and that the container was rebuilt.
-   **Permission denied**: Ensure `EC2_USER` is in the `docker` group (`sudo usermod -aG docker ubuntu`).
