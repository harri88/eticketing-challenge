from decimal import Decimal
from app.domain.models import LedgerEntry, EntryType
from app.domain.interfaces import ILedgerRepository

class LedgerUseCase:
    def __init__(self, repo: ILedgerRepository):
        self.repo = repo

    def record_successful_payment(self, transaction_id: str, amount: float):
        decimal_amount = Decimal(str(amount))
        
        # SOLID: Single Responsibility - Defining the double-entry rule
        debit = LedgerEntry(
            transaction_id=transaction_id,
            account_name="Cash_Asset",
            entry_type=EntryType.DEBIT,
            amount=decimal_amount
        )
        
        credit = LedgerEntry(
            transaction_id=transaction_id,
            account_name="Ticket_Revenue",
            entry_type=EntryType.CREDIT,
            amount=decimal_amount
        )

        # Persistence via interface (Dependency Inversion)
        self.repo.save_balanced_entries([debit, credit])