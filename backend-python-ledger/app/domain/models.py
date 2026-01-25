from dataclasses import dataclass
from datetime import datetime
from decimal import Decimal
from enum import Enum

class EntryType(str, Enum):
    DEBIT = "DEBIT"
    CREDIT = "CREDIT"

@dataclass
class LedgerEntry:
    transaction_id: str
    account_name: str
    entry_type: EntryType
    amount: Decimal
    created_at: datetime = datetime.now()