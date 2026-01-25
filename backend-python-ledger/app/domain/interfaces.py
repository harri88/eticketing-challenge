from abc import ABC, abstractmethod
from app.domain.models import LedgerEntry

class ILedgerRepository(ABC):
    @abstractmethod
    def save_balanced_entries(self, entries: list[LedgerEntry]):
        """Ensures atomic persistence of double-entry records"""
        pass