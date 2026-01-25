import React, { useState, useEffect } from 'react';
import { fetchTicketAvailability, fetchTransactions, fetchLedgerEntries } from './services/api';
import './App.css';

const Backoffice = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [apiKey, setApiKey] = useState('');
  const [loginError, setLoginError] = useState('');
  const [data, setData] = useState({ tickets: [], transactions: [], ledger: [] });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState('overview');

  // Load data on successful authentication
  useEffect(() => {
    if (isAuthenticated) {
      loadAllData();
      // Refresh data every 30 seconds
      const interval = setInterval(loadAllData, 30000);
      return () => clearInterval(interval);
    }
  }, [isAuthenticated]);

  const loadAllData = async () => {
    setLoading(true);
    setError('');
    try {
      const [tickets, transactions, ledger] = await Promise.all([
        fetchTicketAvailability(),
        fetchTransactions(),
        fetchLedgerEntries()
      ]);
      setData({
        tickets: tickets.data || tickets || [],
        transactions: transactions.data || transactions || [],
        ledger: ledger.data || ledger || []
      });
    } catch (err) {
      setError('Failed to load data. Check that all backend services are running.');
      console.error('Data loading error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleLogin = (e) => {
    e.preventDefault();
    // Simplified auth: accept any non-empty API key
    if (apiKey.trim().length > 0) {
      setIsAuthenticated(true);
      setLoginError('');
      localStorage.setItem('backoffice_key', apiKey);
    } else {
      setLoginError('Please enter an API key');
    }
  };

  const handleLogout = () => {
    setIsAuthenticated(false);
    setApiKey('');
    localStorage.removeItem('backoffice_key');
  };

  // Login Screen
  if (!isAuthenticated) {
    return (
      <div className="login-container">
        <div className="login-box">
          <h1>E-Ticketing Backoffice</h1>
          <p className="subtitle">Admin Dashboard</p>
          
          <form onSubmit={handleLogin}>
            <div className="form-group">
              <label>API Key</label>
              <input
                type="password"
                placeholder="Enter admin API key"
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
                className="form-input"
              />
            </div>
            {loginError && <div className="error-message">{loginError}</div>}
            <button type="submit" className="btn-primary">Login</button>
          </form>
          
          <div className="login-hint">
            <p><strong>Demo Credentials:</strong></p>
            <p>API Key: <code>admin-key-123</code></p>
          </div>
        </div>
      </div>
    );
  }

  // Main Dashboard
  return (
    <div className="admin-dashboard">
      <header className="dashboard-header">
        <div className="header-left">
          <h1>E-Ticketing Backoffice</h1>
          <p className="status-indicator">● System Online</p>
        </div>
        <div className="header-right">
          <button onClick={loadAllData} className="btn-refresh" disabled={loading}>
            {loading ? '⟳ Loading...' : '⟳ Refresh'}
          </button>
          <button onClick={handleLogout} className="btn-logout">Logout</button>
        </div>
      </header>

      {error && <div className="error-banner">{error}</div>}

      <nav className="dashboard-nav">
        <button
          className={`nav-item ${activeTab === 'overview' ? 'active' : ''}`}
          onClick={() => setActiveTab('overview')}
        >
          Overview
        </button>
        <button
          className={`nav-item ${activeTab === 'transactions' ? 'active' : ''}`}
          onClick={() => setActiveTab('transactions')}
        >
          Transactions
        </button>
        <button
          className={`nav-item ${activeTab === 'ledger' ? 'active' : ''}`}
          onClick={() => setActiveTab('ledger')}
        >
          Ledger
        </button>
        <button
          className={`nav-item ${activeTab === 'tickets' ? 'active' : ''}`}
          onClick={() => setActiveTab('tickets')}
        >
          Ticket Availability
        </button>
      </nav>

      <main className="dashboard-content">
        {activeTab === 'overview' && (
          <div className="overview-section">
            <h2>Dashboard Overview</h2>
            <div className="stats-grid">
              <div className="stat-card">
                <div className="stat-number">{data.transactions.length}</div>
                <div className="stat-label">Total Transactions</div>
              </div>
              <div className="stat-card">
                <div className="stat-number">{data.tickets.reduce((sum, t) => sum + (t.quota || 0), 0)}</div>
                <div className="stat-label">Total Ticket Quota</div>
              </div>
              <div className="stat-card">
                <div className="stat-number">{data.ledger.length}</div>
                <div className="stat-label">Ledger Entries</div>
              </div>
              <div className="stat-card">
                <div className="stat-number">
                  {data.transactions.filter(t => t.status === 'success').length}
                </div>
                <div className="stat-label">Successful Payments</div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'tickets' && (
          <section className="card">
            <h2>Ticket Quota Management</h2>
            <table className="data-table">
              <thead>
                <tr>
                  <th>Ticket Type</th>
                  <th>Price (AED)</th>
                  <th>Total Quota</th>
                  <th>Remaining</th>
                  <th>Sold</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {data.tickets.length > 0 ? (
                  data.tickets.map(t => (
                    <tr key={t.id || t.name}>
                      <td className="strong">{t.name}</td>
                      <td>{t.price} AED</td>
                      <td>{t.quota || 0}</td>
                      <td className="quota-remaining">{t.remaining || t.quota || 0}</td>
                      <td>{(t.quota || 0) - (t.remaining || 0)}</td>
                      <td>
                        <span className={`status-badge ${(t.remaining || 0) > 0 ? 'available' : 'sold-out'}`}>
                          {(t.remaining || 0) > 0 ? 'Available' : 'Sold Out'}
                        </span>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan="6" className="no-data">No ticket data available</td>
                  </tr>
                )}
              </tbody>
            </table>
          </section>
        )}

        {activeTab === 'transactions' && (
          <section className="card">
            <h2>Payment Transactions</h2>
            <table className="data-table">
              <thead>
                <tr>
                  <th>Transaction ID</th>
                  <th>Payment Method</th>
                  <th>Amount (AED)</th>
                  <th>Status</th>
                  <th>Timestamp</th>
                </tr>
              </thead>
              <tbody>
                {data.transactions.length > 0 ? (
                  data.transactions.map(tx => (
                    <tr key={tx.id || tx.transaction_id}>
                      <td className="mono">{tx.transaction_id || tx.id}</td>
                      <td>
                        <span className="method-badge">
                          {tx.payment_method === 'credit_card' ? '💳 Credit Card' : '📱 QR/UPI'}
                        </span>
                      </td>
                      <td className="amount">{tx.amount || 0}</td>
                      <td>
                        <span className={`status-badge ${tx.status === 'success' ? 'success' : tx.status === 'pending' ? 'pending' : 'failed'}`}>
                          {tx.status}
                        </span>
                      </td>
                      <td className="text-muted">{tx.created_at ? new Date(tx.created_at).toLocaleString() : 'N/A'}</td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan="5" className="no-data">No transactions yet</td>
                  </tr>
                )}
              </tbody>
            </table>
          </section>
        )}

        {activeTab === 'ledger' && (
          <section className="card">
            <h2>Financial Ledger (Double-Entry Audit Trail)</h2>
            <table className="data-table">
              <thead>
                <tr>
                  <th>Transaction ID</th>
                  <th>Account</th>
                  <th>Entry Type</th>
                  <th>Amount (AED)</th>
                  <th>Timestamp</th>
                </tr>
              </thead>
              <tbody>
                {data.ledger.length > 0 ? (
                  data.ledger.map((entry, i) => (
                    <tr key={i} className={entry.entry_type === 'debit' ? 'debit-row' : 'credit-row'}>
                      <td className="mono">{entry.transaction_id || entry.id || `ENTRY-${i}`}</td>
                      <td className="strong">{entry.account_name || entry.account}</td>
                      <td>
                        <span className={`entry-badge ${entry.entry_type}`}>
                          {entry.entry_type === 'debit' ? '→ Debit' : '← Credit'}
                        </span>
                      </td>
                      <td className="amount">{entry.amount || 0}</td>
                      <td className="text-muted">{entry.created_at ? new Date(entry.created_at).toLocaleString() : 'N/A'}</td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan="5" className="no-data">No ledger entries yet</td>
                  </tr>
                )}
              </tbody>
            </table>
          </section>
        )}
      </main>
    </div>
  );
};

export default Backoffice;