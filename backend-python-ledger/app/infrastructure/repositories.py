from sqlalchemy.orm import Session
from app.domain.interfaces import ILedgerRepository
from app.domain.models import LedgerEntry
from app.infrastructure.database import LedgerModel # SQLAlchemy Model

class SQLAlchemyLedgerRepository(ILedgerRepository):
    def __init__(self, db: Session):
        self.db = db

    def save_balanced_entries(self, entries: list[LedgerEntry]):
        try:
            for entry in entries:
                db_entry = LedgerModel(
                    transaction_id=entry.transaction_id,
                    account_name=entry.account_name,
                    entry_type=entry.entry_type,
                    amount=entry.amount
                )
                self.db.add(db_entry)
            self.db.commit() # Atomic commit for both rows
        except Exception:
            self.db.rollback()
            raise