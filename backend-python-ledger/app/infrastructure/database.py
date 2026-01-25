import os
from dotenv import load_dotenv
from sqlalchemy import create_engine, Column, Integer, String, Numeric, DateTime, func
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
from datetime import datetime
import logging
from urllib.parse import quote

# Load environment variables from .env file
load_dotenv()

logger = logging.getLogger(__name__)

# Database configuration
DATABASE_HOST = os.getenv("DB_HOST", "localhost")
DATABASE_PORT = os.getenv("DB_PORT", "5434")
DATABASE_USER = os.getenv("DB_USER", "postgres")
DATABASE_PASSWORD = os.getenv("DB_PASSWORD", "postgres@12345")
DATABASE_NAME = os.getenv("DB_NAME", "ledger_db")

# URL encode the password to handle special characters like @
encoded_password = quote(DATABASE_PASSWORD, safe='')
DATABASE_URL = f"postgresql://{DATABASE_USER}:{encoded_password}@{DATABASE_HOST}:{DATABASE_PORT}/{DATABASE_NAME}"

# Create engine
engine = create_engine(
    DATABASE_URL,
    pool_size=10,
    max_overflow=20,
    pool_pre_ping=True
)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()

# SQLAlchemy Model
class LedgerModel(Base):
    __tablename__ = "ledger_entries"
    
    id = Column(Integer, primary_key=True, index=True)
    transaction_id = Column(String(50), index=True, nullable=False)
    account_name = Column(String(100), nullable=False)  # e.g., "Cash_Asset", "Ticket_Revenue"
    entry_type = Column(String(10), nullable=False)  # "DEBIT" or "CREDIT"
    amount = Column(Numeric(15, 2), nullable=False)
    created_at = Column(DateTime, server_default=func.now(), nullable=False)

def init_db():
    """Initialize database and create tables"""
    try:
        Base.metadata.create_all(bind=engine)
        logger.info("Database tables created successfully")
    except Exception as e:
        logger.error(f"Error initializing database: {str(e)}")
        raise

def get_db() -> Session:
    """Dependency for getting database session"""
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
