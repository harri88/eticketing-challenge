# E-Ticketing Challenge System

Welcome to the E-Ticketing Challenge System! This is a full-stack, polyglot microservices application designed to demonstrate modern software architecture principles, including Clean Architecture, containerization, and automated CI/CD.

## 🏗 Architecture Diagram

The system follows a microservices architecture pattern, where each domain is isolated with its own database and business logic. The services interact over HTTP REST APIs.



```mermaid
graph TD
    User((Customer)) -->|HTTP/Browsers| FE_React[frontend-react]
    Admin((Admin)) -->|HTTP/Browsers| FE_Back[frontend-backoffice]

    subgraph "Docker Network"
        FE_React -->|REST| TICKETS[backend-dotnet-tickets]
        FE_React -->|REST| PAYMENTS[backend-go-payment]
        
        FE_Back -->|REST| TICKETS
        FE_Back -->|REST| PAYMENTS
        FE_Back -->|REST| LEDGER[backend-python-ledger]

        TICKETS -->|DB Protocol| DB_T[(Ticket DB)]
        PAYMENTS -->|DB Protocol| DB_P[(Payment DB)]
        LEDGER -->|DB Protocol| DB_L[(Ledger DB)]

        PAYMENTS -.->|REST Internal| TICKETS
        PAYMENTS -.->|REST Internal| LEDGER
    end
```

## 🗄️ Database Entity Relationship Diagram (ERD)

The following diagram illustrates the data models across the three isolated databases and their logical relationships. Note that cross-service references (dotted lines) are logical links, not foreign keys.

```mermaid
erDiagram
    %% Ticket Service Database (PostgreSQL)
    TICKET_DB_TICKETS {
        varchar id PK
        varchar name
        decimal price
        varchar currency
        int quota
        int held_quota
    }
    TICKET_DB_ORDERS {
        varchar id PK
        varchar customer_email
        decimal total_amount
        varchar currency
        varchar status
        timestamp created_at
        timestamp paid_at
    }
    TICKET_DB_ORDER_ITEMS {
        serial id PK
        varchar order_id FK
        varchar ticket_id FK
        int quantity
        decimal price_at_purchase
    }

    TICKET_DB_ORDERS ||--|{ TICKET_DB_ORDER_ITEMS : "contains"
    TICKET_DB_TICKETS ||--o{ TICKET_DB_ORDER_ITEMS : "reserved_in"

    %% Payment Service Database (PostgreSQL)
    PAYMENT_DB_TRANSACTIONS {
        serial id PK
        varchar transaction_id UK "Unique Payment Ref"
        varchar order_id "Ref to TICKET_DB_ORDERS"
        varchar payment_method
        varchar payment_ref "Provider Ref"
        decimal amount
        varchar currency
        varchar status
        timestamp created_at
    }

    %% Logical link: Transaction pays for an Order
    TICKET_DB_ORDERS ||--o| PAYMENT_DB_TRANSACTIONS : "paid_via (Logical)"

    %% Ledger Service Database (PostgreSQL)
    LEDGER_DB_ENTRIES {
        serial id PK
        varchar transaction_id "Ref to PAYMENT_DB_TRANSACTIONS"
        varchar account_name
        varchar entry_type "DEBIT/CREDIT"
        decimal amount
        text description
        timestamp created_at
    }

    %% Logical link: Transaction generates Ledger Entries
    PAYMENT_DB_TRANSACTIONS ||--|{ LEDGER_DB_ENTRIES : "audited_by (Logical)"
```

## 📜 Sequence Diagrams

### 1. User Journey (Get Tickets → Create Order → Payment)

The following sequence diagram illustrates the complete flow of a user selecting tickets, reserving them (Creating Order), and eventually paying for them (which triggers Ledger recording).

```mermaid
sequenceDiagram
    autonumber
    actor U as User (Frontend)
    participant T as Ticket Service (.NET)
    participant DB_T as Ticket DB
    participant P as Payment Service (Go)
    participant DB_P as Payment DB
    participant L as Ledger Service (Python)
    participant DB_L as Ledger DB

    %% 1. Get Tickets
    Note over U, T: Step 1: Browse Available Tickets
    U->>T: GET /api/v1/tickets
    T->>DB_T: Select * from tickets
    DB_T-->>T: [Gold, Premium, VIP...]
    T-->>U: JSON { data: [tickets...] }

    %% 2. Create Order
    Note over U, T: Step 2: Checkout / Reserve
    U->>T: POST /api/v1/checkout/orders
    Note right of U: { customer_email, cart_items }
    T->>DB_T: Check Quota & Create Order (Pending)
    DB_T-->>T: Order Created (ID: ORD-123)
    T-->>U: JSON { order_id: "ORD-123", total: 500 }

    %% 3. Payment
    Note over U, P: Step 3: Process Payment
    U->>P: POST /api/v1/payments
    Note right of U: { order_id, amount, method, Idempotency-Key }
    
    rect rgb(240, 248, 255)
        Note right of P: Idempotency Check
        P->>DB_P: Create Transaction (PENDING)
        
        P->>P: Process Gateway (Stripe/Mock)
        
        alt Payment Success
            P->>DB_P: Update Transaction (SUCCESS)
            
            par Async Notifications
                P->>T: POST /api/v1/webhooks/payment-success
                T->>DB_T: Update Order -> PAID
                
                P->>L: POST /api/v1/ledger
                L->>DB_L: Insert Debit/Credit Entries
            end

            P-->>U: JSON { status: "SUCCESS", tx_id: "TXN-123" }
        else Payment Failed
            P->>DB_P: Update Transaction (FAILED)
            P-->>U: 400 Bad Request
        end
    end
```

## 🧩 Service Responsibilities

| Service | Stack | Port | Description |
| :--- | :--- | :--- | :--- |
| **backend-dotnet-tickets** | .NET 10 (C#) | `5020` | Core domain for Ticket inventory and management. Implements **Clean Architecture**. |
| **backend-go-payment** | Go 1.25 | `8081` | Handles payment processing and gateway integrations (Mock/Stripe). |
| **backend-python-ledger** | Python 3.11 | `8000` | Double-entry accounting ledger to record all financial transactions auditably. |
| **frontend-react** | React 19 | `3000` | Customer-facing application for browsing tickets, cart management, and checkout. |
| **frontend-backoffice** | React 19 | `3001` | Administrative dashboard for viewing sales, inventory, and ledger entries. |

## 🌐 Live Demo

The functionality is deployed and accessible at the following public endpoints:

- **Customer Frontend**: [http://3.106.254.198:3000](http://3.106.254.198:3000)
- **Backoffice Dashboard**: [http://3.106.254.198:3001](http://3.106.254.198:3001)
- **APIs & Documentation**:
  - Ticket Service (Swagger): [http://3.106.254.198:5020/swagger/index.html](http://3.106.254.198:5020/swagger/index.html)
  - Payment Service (Docs): [http://3.106.254.198:8081/swagger/index.html](http://3.106.254.198:8081/swagger/index.html)
  - Ledger Service (Docs): [http://3.106.254.198:8000/docs](http://3.106.254.198:8000/docs)

## 🚀 How to Run Locally

### Prerequisites
- Docker & Docker Compose installed
- Git

### Steps
1. **Clone the repository**:
   ```bash
   git clone <repo-url>
   cd eticketing-challenge
   ```

2. **Configure Environment**:
   Create a `.env` file at the root (or verify the existing one):
   ```bash
   PUBLIC_IP=localhost
   # PUBLIC_DOMAIN=localhost (optional)
   
   # These help the frontends find the backends locally
   REACT_APP_API_URL=http://localhost:5020
   PAYMENT_API_URL=http://localhost:8081
   LEDGER_API_URL=http://localhost:8000
   ```

3. **Start the System**:
   ```bash
   docker-compose up -d --build
   ```

4. **Access the Applications**:
   - Customer App: [http://localhost:3000](http://localhost:3000)
   - Backoffice: [http://localhost:3001](http://localhost:3001)
   - Ticket API Swagger: [http://localhost:5020/swagger](http://localhost:5020/swagger)
   - Payment API Docs: [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html) (if enabled)
   - Ledger API Docs: [http://localhost:8000/docs](http://localhost:8000/docs)

## ☁️ How to Deploy (EC2 / Production)

We use **GitHub Actions** for CI/CD. The deployment pipeline automatically deploys changes to an AWS EC2 instance.

1.  **Provision EC2**: Ubuntu Server.
2.  **Setup Server**: Run keys setup script inside `scripts/setup_ec2.sh`.
3.  **Secrets**: Add `EC2_HOST`, `EC2_USER`, and `EC2_SSH_KEY` to GitHub Repo Secrets.
4.  **Deploy**: Push to `main` branch.

See full guide in: [DEPLOY_EC2.md](./DEPLOY_EC2.md)

## 📦 Package Distribution

Create distributable packages for all services with version control:

```bash
cd scripts
VERSION=1.0.0 ./package-all.sh
```

This creates `.tar.gz` archives in the `dist/` directory for deployment to any environment. Each package includes compiled binaries, configuration templates, and documentation.

See full guide in: [PACKAGING.md](./PACKAGING.md)

## 🔑 Environment Variables

The system relies on a central `.env` on the production server (generated by CI/CD) or locally.

| Variable | Description | Default (Local) |
| :--- | :--- | :--- |
| `PUBLIC_IP` | The public IP of the host/server. | `localhost` |
| `REACT_APP_API_URL` | Base URL for the Ticket Service used by Frontends. | `http://localhost:5020` |
| `TICKET_API_URL` | Internal/External URL for Ticket Service. | `http://localhost:5020` |
| `PAYMENT_API_URL` | Internal/External URL for Payment Service. | `http://localhost:8081` |
| `LEDGER_API_URL` | Internal/External URL for Ledger Service. | `http://localhost:8000` |

## 📡 API Examples

### 1. Get Available Tickets
```bash
curl -X GET "http://localhost:5020/api/v1/tickets" -H "accept: application/json"
```

### 2. Create an Order (Checkout)
```bash
curl -X POST "http://localhost:5020/api/v1/checkout/orders" \
     -H "Content-Type: application/json" \
     -d '{
           "customer_email": "john@example.com",
           "cart_items": [
             { "ticket_id": "gold_001", "quantity": 2 }
           ]
         }'
```

### 3. Process Payment
```bash
curl -X POST "http://localhost:8081/api/v1/payments" \
     -H "Content-Type: application/json" \
     -H 'Idempotency-Key: b4e62cc4-cc09-4b07-a00a-42dcd5cae281' \
     -d '{
           "order_id": "ORD-12345",
           "amount": 200.00,
           "currency": "AED",
           "payment_method": "credit_card"
         }'
curl 'http://3.106.254.198:8081/api/v1/payments' \
      -H 'Referer: http://3.106.254.198:3000/' \
      -H 'Idempotency-Key: b4e62cc4-cc09-4b07-a00a-42dcd5cae281' \
      -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36' \
      -H 'Content-Type: application/json' \
      --data-raw '{"order_id":"ORD-19237671A0","amount":200,"payment_method":"qr_scan","card_number":"","card_expiry":"","card_cvc":""}'
```

## ⚖️ Assumptions & Trade-offs

1.  **Shared Database vs Isolated DB**: 
    -   *Decision*: We chose **Isolated Databases** (Database-per-service pattern).
    -   *Trade-off*: Higher complexity in infrastructure (3 Postgres containers), but ensures loose coupling and independent scaling.

2.  **Communication Protocol**:
    -   *Decision*: Synchronous **HTTP/REST**.
    -   *Trade-off*: Easier to implement and debug than event-driven (RabbitMQ/Kafka), but introduces temporal coupling (if Ticket Service is down, Payment might fail to validate inventory).

3.  **Authentication**:
    -   *Assumption*: For this challenge scope, APIs are public or use simple internal checks. A production system would require OAuth2/OIDC (e.g., Auth0 or Keycloak).

4.  **Frontend Configuration**:
    -   *Decision*: Build-time environment variables (`ARG`) in Docker.
    -   *Trade-off*: Requires rebuilding the container to change the API URL. Runtime configuration (fetching `config.json` on startup) would be more flexible but slightly more complex to setup.

5.  **Polyglot Stack**:
    -   *Reason*: To demonstrate proficiency in multiple modern languages (.NET, Go, Python).
    -   *Trade-off*: Increased operational complexity (need meaningful CI pipelines for all 3 stacks).

## 🛠 CI/CD & Code Quality

-   **GitHub Actions**: Pipelines created for all 5 services ([.github/workflows](./.github/workflows)).
-   **SonarQube**: Integration added for static code analysis. See [SONARQUBE_SETUP.md](./SONARQUBE_SETUP.md).

---
*Generated for E-Ticketing Challenge*


[def]: https://ibb.co/h1DGkNGd
