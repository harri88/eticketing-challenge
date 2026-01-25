const TICKET_SERVICE_URL = process.env.REACT_APP_TICKET_SERVICE || 'http://localhost:7175';
const PAYMENT_SERVICE_URL = process.env.REACT_APP_PAYMENT_SERVICE || 'http://localhost:8081';
const LEDGER_SERVICE_URL = process.env.REACT_APP_LEDGER_SERVICE || 'http://localhost:8000';

const handleResponse = async (response) => {
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }
  return response.json();
};

export const fetchTicketAvailability = async () => {
  try {
    const response = await fetch(`${TICKET_SERVICE_URL}/api/v1/tickets`);
    const data = await handleResponse(response);
    return data;
  } catch (error) {
    console.error('Error fetching tickets:', error);
    return { data: [] }; // Return empty array on error
  }
};

export const fetchTransactions = async () => {
  try {
    const response = await fetch(`${PAYMENT_SERVICE_URL}/api/v1/transactions`);
    const data = await handleResponse(response);
    return data;
  } catch (error) {
    console.error('Error fetching transactions:', error);
    return { data: [] }; // Return empty array on error
  }
};

export const fetchLedgerEntries = async () => {
  try {
    const response = await fetch(`${LEDGER_SERVICE_URL}/api/v1/ledger`);
    const data = await handleResponse(response);
    return data;
  } catch (error) {
    console.error('Error fetching ledger entries:', error);
    return { data: [] }; // Return empty array on error
  }
};