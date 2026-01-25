from fastapi import FastAPI, Depends, HTTPException, Query
from fastapi.middleware.cors import CORSMiddleware
from fastapi.openapi.utils import get_openapi
from sqlalchemy.orm import Session
from app.application.ledger_use_case import LedgerUseCase
from app.infrastructure.repositories import SQLAlchemyLedgerRepository
from app.infrastructure.database import get_db, init_db, SessionLocal, LedgerModel
from pydantic import BaseModel, Field
from typing import List, Optional
import logging
import os
from datetime import datetime
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Create FastAPI app with comprehensive OpenAPI configuration
app = FastAPI(
    title="E-Ticketing Ledger Service",
    description="""
    # Double-Entry Accounting Ledger Service

    A robust financial ledger system for the e-ticketing platform that implements 
    double-entry bookkeeping principles. Every payment transaction is recorded as a 
    balanced pair of entries (debit and credit) ensuring financial integrity and 
    audit compliance.

    ## Key Features

    - **Double-Entry Accounting**: Every transaction creates balanced debit/credit entries
    - **ACID Transactions**: Atomic commits ensure data consistency
    - **Audit Trail**: Complete immutable transaction history
    - **RESTful API**: Simple JSON-based endpoints
    - **CORS Enabled**: Accessible from multiple frontend services
    - **Production Ready**: Error handling, logging, and validation

    ## Account Types

    - **Cash_Asset**: Represents money received from payment methods
    - **Ticket_Revenue**: Represents service revenue from ticket sales

    ## Double-Entry Example

    When a customer purchases a 200 AED Premium ticket:
    ```
    Debit   | Cash_Asset      | 200 AED  (Money in)
    Credit  | Ticket_Revenue  | 200 AED  (Service out)
    ```

    This ensures: **Total Debits = Total Credits** (always balanced)
    """,
    version="1.0.0",
    contact={
        "name": "E-Ticketing Team",
        "email": "support@eticketing.local"
    },
    license_info={
        "name": "MIT License"
    }
)

# CORS Configuration
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Initialize database on startup
@app.on_event("startup")
def startup():
    init_db()
    logger.info("Database initialized")

# ==================== Request Models ====================

class PaymentEvent(BaseModel):
    """
    Request model for recording a new payment transaction.
    
    This creates two balanced ledger entries:
    - Debit: Cash_Asset (payment received)
    - Credit: Ticket_Revenue (service provided)
    """
    transaction_id: str = Field(
        ..., 
        description="Unique transaction identifier from payment service",
        example="txn-20260125-001",
        min_length=1,
        max_length=50
    )
    amount: float = Field(
        ..., 
        description="Payment amount in AED (United Arab Emirates Dirham)",
        example=200.00,
        gt=0
    )

    class Config:
        json_schema_extra = {
            "example": {
                "transaction_id": "txn-20260125-001",
                "amount": 200.00
            }
        }

# ==================== Response Models ====================

class LedgerEntryResponse(BaseModel):
    """Single ledger entry in the accounting ledger."""
    id: Optional[int] = Field(None, description="Database record ID")
    transaction_id: str = Field(..., description="Associated transaction ID")
    account_name: str = Field(..., description="Account name (Cash_Asset or Ticket_Revenue)")
    entry_type: str = Field(..., description="Entry type: DEBIT or CREDIT")
    amount: float = Field(..., description="Entry amount in AED")
    created_at: str = Field(..., description="Entry creation timestamp (ISO 8601)")
    
    class Config:
        from_attributes = True
        json_schema_extra = {
            "example": {
                "id": 1,
                "transaction_id": "txn-20260125-001",
                "account_name": "Cash_Asset",
                "entry_type": "DEBIT",
                "amount": 200.00,
                "created_at": "2026-01-25T10:30:00"
            }
        }

class LedgerDataResponse(BaseModel):
    """Response wrapper for ledger entries list."""
    data: List[LedgerEntryResponse] = Field(..., description="Array of ledger entries")
    
    class Config:
        json_schema_extra = {
            "example": {
                "data": [
                    {
                        "id": 1,
                        "transaction_id": "txn-001",
                        "account_name": "Cash_Asset",
                        "entry_type": "DEBIT",
                        "amount": 200.00,
                        "created_at": "2026-01-25T10:30:00"
                    },
                    {
                        "id": 2,
                        "transaction_id": "txn-001",
                        "account_name": "Ticket_Revenue",
                        "entry_type": "CREDIT",
                        "amount": 200.00,
                        "created_at": "2026-01-25T10:30:00"
                    }
                ]
            }
        }

class RecordPaymentResponse(BaseModel):
    """Response for successful payment recording."""
    message: str = Field(..., description="Success message")
    transaction_id: str = Field(..., description="Recorded transaction ID")
    
    class Config:
        json_schema_extra = {
            "example": {
                "message": "Ledger entries recorded successfully",
                "transaction_id": "txn-20260125-001"
            }
        }

class RecordPaymentDataResponse(BaseModel):
    """Wrapper response for payment recording."""
    data: RecordPaymentResponse

class ErrorResponse(BaseModel):
    """Standard error response."""
    detail: str = Field(..., description="Error message")
    
    class Config:
        json_schema_extra = {
            "example": {
                "detail": "Failed to record ledger entries"
            }
        }

class HealthResponse(BaseModel):
    """Health check response."""
    status: str = Field(..., description="Service status")
    service: str = Field(..., description="Service name")


# ==================== Endpoints ====================

@app.get(
    "/health",
    response_model=HealthResponse,
    tags=["Health"],
    summary="Service Health Check",
    description="Check if the ledger service is operational and connected to the database."
)
def health():
    """
    Health check endpoint for monitoring and orchestration.
    
    Returns:
        HealthResponse: Service status and name
    """
    return {"status": "healthy", "service": "ledger-service"}


@app.post(
    "/api/v1/ledger",
    response_model=RecordPaymentDataResponse,
    status_code=201,
    tags=["Ledger Entries"],
    summary="Record Payment Transaction",
    description="Record a successful payment as balanced double-entry ledger entries.",
    responses={
        201: {
            "description": "Payment ledger entries recorded successfully",
            "model": RecordPaymentDataResponse,
            "content": {
                "application/json": {
                    "example": {
                        "data": {
                            "message": "Ledger entries recorded successfully",
                            "transaction_id": "txn-20260125-001"
                        }
                    }
                }
            }
        },
        400: {
            "description": "Invalid request data",
            "model": ErrorResponse
        },
        500: {
            "description": "Server error while recording entries",
            "model": ErrorResponse
        }
    }
)
def record_payment(
    event: PaymentEvent,
    db: Session = Depends(get_db)
):
    """
    Record a payment transaction as balanced ledger entries.
    
    This endpoint implements double-entry bookkeeping by creating two entries:
    1. **Debit** entry to Cash_Asset account (money received)
    2. **Credit** entry to Ticket_Revenue account (service provided)
    
    Both entries are created atomically - either both succeed or both fail.
    
    Args:
        event (PaymentEvent): Payment details including transaction ID and amount
        db (Session): Database session (injected dependency)
    
    Returns:
        RecordPaymentDataResponse: Confirmation of recorded entries
    
    Raises:
        HTTPException: 500 if database operation fails
    
    Example:
        ```
        POST /api/v1/ledger
        Content-Type: application/json
        
        {
            "transaction_id": "txn-20260125-001",
            "amount": 200.00
        }
        
        Response:
        {
            "data": {
                "message": "Ledger entries recorded successfully",
                "transaction_id": "txn-20260125-001"
            }
        }
        ```
    """
    try:
        repo = SQLAlchemyLedgerRepository(db)
        use_case = LedgerUseCase(repo)
        
        use_case.record_successful_payment(event.transaction_id, event.amount)
        logger.info(f"Ledger entries recorded for transaction {event.transaction_id}")
        
        return {
            "data": {
                "message": "Ledger entries recorded successfully",
                "transaction_id": event.transaction_id
            }
        }
    except Exception as e:
        logger.error(f"Error recording ledger entries: {str(e)}")
        raise HTTPException(status_code=500, detail="Failed to record ledger entries")


@app.get(
    "/api/v1/ledger",
    response_model=LedgerDataResponse,
    tags=["Ledger Entries"],
    summary="Retrieve All Ledger Entries",
    description="Fetch complete audit trail of all ledger entries with most recent first.",
    responses={
        200: {
            "description": "List of all ledger entries",
            "model": LedgerDataResponse
        },
        500: {
            "description": "Server error while retrieving entries",
            "model": ErrorResponse
        }
    }
)
def get_ledger_entries(db: Session = Depends(get_db)):
    """
    Retrieve all ledger entries sorted by creation date (newest first).
    
    This endpoint returns the complete financial audit trail showing all
    debit and credit entries for all transactions.
    
    Args:
        db (Session): Database session (injected dependency)
    
    Returns:
        LedgerDataResponse: Array of all ledger entries
    
    Raises:
        HTTPException: 500 if database query fails
    
    Example:
        ```
        GET /api/v1/ledger
        
        Response:
        {
            "data": [
                {
                    "id": 2,
                    "transaction_id": "txn-001",
                    "account_name": "Ticket_Revenue",
                    "entry_type": "CREDIT",
                    "amount": 200.00,
                    "created_at": "2026-01-25T10:30:00"
                },
                {
                    "id": 1,
                    "transaction_id": "txn-001",
                    "account_name": "Cash_Asset",
                    "entry_type": "DEBIT",
                    "amount": 200.00,
                    "created_at": "2026-01-25T10:30:00"
                }
            ]
        }
        ```
    """
    try:
        entries = db.query(LedgerModel).order_by(LedgerModel.created_at.desc()).all()
        return {
            "data": [
                {
                    "id": e.id,
                    "transaction_id": e.transaction_id,
                    "account_name": e.account_name,
                    "entry_type": e.entry_type,
                    "amount": float(e.amount),
                    "created_at": e.created_at.isoformat()
                }
                for e in entries
            ]
        }
    except Exception as e:
        logger.error(f"Error retrieving ledger entries: {str(e)}")
        raise HTTPException(status_code=500, detail="Failed to retrieve ledger entries")


@app.get(
    "/api/v1/ledger/{transaction_id}",
    response_model=LedgerDataResponse,
    tags=["Ledger Entries"],
    summary="Retrieve Transaction Ledger Entries",
    description="Get the debit and credit entries for a specific transaction.",
    responses={
        200: {
            "description": "Ledger entries for the specified transaction",
            "model": LedgerDataResponse
        },
        404: {
            "description": "Transaction not found",
            "model": ErrorResponse
        },
        500: {
            "description": "Server error while retrieving entries",
            "model": ErrorResponse
        }
    }
)
def get_transaction_ledger(
    transaction_id: str,
    db: Session = Depends(get_db)
):
    """
    Retrieve ledger entries for a specific transaction.
    
    Returns both the debit and credit entries that make up the
    complete double-entry record for a single payment transaction.
    
    Args:
        transaction_id (str): The transaction ID to lookup
        db (Session): Database session (injected dependency)
    
    Returns:
        LedgerDataResponse: Array containing debit and credit entries for the transaction
    
    Raises:
        HTTPException: 404 if transaction not found, 500 if query fails
    
    Example:
        ```
        GET /api/v1/ledger/txn-20260125-001
        
        Response:
        {
            "data": [
                {
                    "id": 2,
                    "transaction_id": "txn-20260125-001",
                    "account_name": "Ticket_Revenue",
                    "entry_type": "CREDIT",
                    "amount": 200.00,
                    "created_at": "2026-01-25T10:30:00"
                },
                {
                    "id": 1,
                    "transaction_id": "txn-20260125-001",
                    "account_name": "Cash_Asset",
                    "entry_type": "DEBIT",
                    "amount": 200.00,
                    "created_at": "2026-01-25T10:30:00"
                }
            ]
        }
        ```
    """
    try:
        entries = db.query(LedgerModel).filter(
            LedgerModel.transaction_id == transaction_id
        ).all()
        
        if not entries:
            raise HTTPException(
                status_code=404,
                detail=f"No ledger entries found for transaction {transaction_id}"
            )
        
        return {
            "data": [
                {
                    "id": e.id,
                    "transaction_id": e.transaction_id,
                    "account_name": e.account_name,
                    "entry_type": e.entry_type,
                    "amount": float(e.amount),
                    "created_at": e.created_at.isoformat()
                }
                for e in entries
            ]
        }
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error retrieving transaction ledger: {str(e)}")
        raise HTTPException(status_code=500, detail="Failed to retrieve transaction ledger")


# ==================== OpenAPI Customization ====================

def custom_openapi():
    """
    Custom OpenAPI schema with enhanced documentation.
    """
    if app.openapi_schema:
        return app.openapi_schema
    
    openapi_schema = get_openapi(
        title="E-Ticketing Ledger Service",
        version="1.0.0",
        description=app.description,
        routes=app.routes,
    )
    
    # Add server information
    openapi_schema["servers"] = [
        {
            "url": "http://localhost:8000",
            "description": "Local Development Server"
        },
        {
            "url": "https://api.eticketing.local",
            "description": "Production Server"
        }
    ]
    
    # Add security schemes if needed in future
    openapi_schema["components"]["schemas"]["LedgerEntryResponse"]["properties"]["entry_type"]["enum"] = ["DEBIT", "CREDIT"]
    
    app.openapi_schema = openapi_schema
    return app.openapi_schema

app.openapi = custom_openapi