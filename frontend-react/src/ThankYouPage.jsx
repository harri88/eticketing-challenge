import React from 'react';
import './App.css'; 

const ThankYouPage = ({ receipt, cart, tickets, onReset, email }) => {
  
  // 1. Helper: Find ticket details safely
  const getTicket = (id) => tickets.find((t) => t.id === id);

  // 2. Helper: Calculate total one last time for display
  const totalAmount = Object.keys(cart).reduce((sum, id) => {
    const t = getTicket(id);
    return sum + (t ? t.price * cart[id] : 0);
  }, 0);

  return (
    <div className="page-container fade-in">
      <div className="thank-you-card">
        
        {/* --- Header Section --- */}
        <div className="success-header">
          <div className="checkmark-wrapper">
            <svg className="checkmark" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 52 52">
              <circle className="checkmark-circle" cx="26" cy="26" r="25" fill="none"/>
              <path className="checkmark-check" fill="none" d="M14.1 27.2l7.1 7.2 16.7-16.8"/>
            </svg>
          </div>
          <h1 className="success-title">Payment Successful!</h1>
          <p className="success-subtitle">
            A confirmation email has been sent to <br/>
            <strong>{email || 'your registered email'}</strong>
          </p>
        </div>

        {/* --- Receipt Body --- */}
        <div className="receipt-body">
          <h3>Order Summary</h3>
          
          <div className="ticket-list">
            {Object.keys(cart).map((id) => {
              const qty = cart[id];
              if (qty === 0) return null;
              const t = getTicket(id);
              if (!t) return null;

              return (
                <div key={id} className="summary-row">
                  <div className="item-info">
                    <span className="qty-badge">{qty}x</span>
                    <span className="item-name">{t.name}</span>
                  </div>
                  <span className="item-price">{t.price * qty} AED</span>
                </div>
              );
            })}
          </div>

          <div className="dashed-divider"></div>

          {/* Transaction Metadata */}
          <div className="meta-info">
            <div className="meta-row">
              <span className="label">Transaction ID</span>
              <span className="value mono">{receipt?.transactionId || 'N/A'}</span>
            </div>
            <div className="meta-row">
              <span className="label">Reference</span>
              <span className="value mono">{receipt?.paymentRef || 'N/A'}</span>
            </div>
            <div className="meta-row">
              <span className="label">Time</span>
              <span className="value">{receipt?.timestamp}</span>
            </div>
          </div>

          <div className="total-divider"></div>

          <div className="total-row-large">
            <span>Amount Paid</span>
            <span className="success-color">{totalAmount} AED</span>
          </div>
        </div>

        {/* --- Footer Buttons --- */}
        <div className="card-footer">
          <button className="btn-secondary full-width" onClick={() => alert("Simulating PDF Ticket Download...")}>
            📥 Download E-Tickets
          </button>
          <button className="btn-primary full-width" onClick={onReset}>
            Buy More Tickets
          </button>
        </div>

      </div>
    </div>
  );
};

export default ThankYouPage;