import React, { useState, useEffect } from 'react';
import './App.css'; 
import ThankYouPage from './ThankYouPage';

// --- MOCK API SERVICE ---
const apiService = {
    getTickets: async () => {
    try {
      // 1. Call the real .NET API
      const response = await fetch('http://localhost:5020/api/v1/tickets');
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      // 2. Parse JSON
      const result = await response.json();
      
      // 3. Return the 'data' array from the wrapper { data: [...] }
      return result.data; 
    } catch (error) {
      console.error("Could not fetch tickets:", error);
      alert("Error connecting to API. Is the backend running on port 5020?");
      return [];
    }
  },
  createOrder: async (email, cart, tickets) => {
    try {
      // Transform cart object to cart_items array
      const cartItems = tickets
        .filter(t => cart[t.id] && cart[t.id] > 0)
        .map(t => ({
          ticket_id: t.id,
          quantity: cart[t.id]
        }));

      const payload = {
        customer_email: email,
        cart_items: cartItems
      };

      const response = await fetch('http://localhost:5020/api/v1/checkout/orders', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const result = await response.json();
      return result;
    } catch (error) { 
      console.error("Could not create order:", error);
      alert("Error connecting to API. Is the backend running on port 5020?");
      return null;
    }
  }
};

const generateReceipt = () => ({
  transactionId: `TRX-${Math.floor(Math.random() * 100000000)}`,
  paymentRef: `REF-${Math.random().toString(36).substring(7).toUpperCase()}`,
  timestamp: new Date().toLocaleString(),
});

// ==========================================
// COMPONENT 1: TICKET SELECTION (Tidied UI)
// ==========================================
const TicketSelection = ({ tickets, onProceed }) => {
  const [cart, setCart] = useState({});

  const handleQuantityChange = (ticketId, delta, maxQuota) => {
    setCart((prev) => {
      const current = prev[ticketId] || 0;
      const next = current + delta;
      if (next < 0 || next > maxQuota) return prev;
      return { ...prev, [ticketId]: next };
    });
  };

  const calculateTotal = () => tickets.reduce((total, t) => total + (t.price * (cart[t.id] || 0)), 0);
  const totalAmount = calculateTotal();

  return (
    <div className="page-container">
      <header className="app-header">
        <h1>🎟️ Event E-Ticketing</h1>
        <p>Secure your spot today</p>
      </header>
      
      <div className="ticket-cards-container">
        {tickets.map((t) => {
          const quantity = cart[t.id] || 0;
          return (
            <div key={t.id} className={`ticket-card ${quantity > 0 ? 'selected' : ''}`}>
              <div className="card-top">
                <h3>{t.name}</h3>
                <span className="price-tag">{t.price} {t.currency}</span>
              </div>
              <p className="ticket-desc">{t.description || 'General Admission'}</p>
              <div className="quota-badge">Only {t.quota} left</div>
              
              <div className="ticket-controls">
                <button 
                  className="ctrl-btn minus" 
                  onClick={() => handleQuantityChange(t.id, -1, t.quota)}
                  disabled={quantity === 0}
                >−</button>
                <span className="qty-display">{quantity}</span>
                <button 
                  className="ctrl-btn plus" 
                  onClick={() => handleQuantityChange(t.id, 1, t.quota)}
                  disabled={quantity >= t.quota}
                >+</button>
              </div>
            </div>
          );
        })}
      </div>

      <div className="bottom-bar">
        <div className="total-display">
          <span>Total:</span>
          <strong>{totalAmount} AED</strong>
        </div>
        <button 
          className="btn-primary" 
          onClick={() => onProceed(cart, totalAmount)} 
          disabled={totalAmount === 0}
        >
          Review Order &rarr;
        </button>
      </div>
    </div>
  );
};

// ==========================================
// COMPONENT 2: ORDER SUMMARY (Tidied UI)
// ==========================================
const OrderSummaryPage = ({ cart, tickets, total, onCreateOrder, onBack }) => {
  const [email, setEmail] = useState('');
  const [isCreating, setIsCreating] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsCreating(true);
    const response = await apiService.createOrder(email, cart, tickets);
    console.log("Create Order Response:", response);
    if (response && response.order_id) {
      onCreateOrder(response.order_id);
    } else {
      alert("Failed to create order");
      setIsCreating(false);
    }
  };

  // Filter only selected tickets
  const selectedTickets = tickets.filter(t => (cart[t.id] || 0) > 0);

  return (
    <div className="page-container fade-in">
      <button className="btn-text-back" onClick={onBack}>&larr; Back to Selection</button>
      
      <div className="summary-card">
        <h2>Order Summary</h2>
        <div className="summary-items">
          {selectedTickets.map(t => (
            <div key={t.id} className="summary-row">
              <div className="item-info">
                <strong>{t.name}</strong>
                <span className="item-qty">x {cart[t.id]}</span>
              </div>
              <span className="item-price">{t.price * cart[t.id]} AED</span>
            </div>
          ))}
          <div className="divider"></div>
          <div className="total-row-large">
            <span>Total Payable</span>
            <span>{total} AED</span>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="buyer-form">
          <div className="input-group">
            <label htmlFor="email">Where should we send your tickets?</label>
            <input 
              id="email"
              type="email" 
              placeholder="name@example.com" 
              required 
              value={email} 
              onChange={(e) => setEmail(e.target.value)} 
              className="email-input"
            />
          </div>
          <button type="submit" className="btn-primary full-width" disabled={isCreating}>
            {isCreating ? (
              <span className="spinner">↻ Creating Order...</span>
            ) : (
              "Confirm & Proceed to Payment"
            )}
          </button>
        </form>
      </div>
    </div>
  );
};

// ==========================================
// COMPONENT 3: PAYMENT PAGE (Functional)
// ==========================================
const PaymentPage = ({ total, orderId, onPaymentSuccess }) => {
  const [method, setMethod] = useState('cc');
  const [qrTimer, setQrTimer] = useState(8);

  useEffect(() => {
    let interval;
    if (method === 'qr') {
      setQrTimer(8);
      interval = setInterval(() => {
        setQrTimer((prev) => {
          if (prev <= 1) {
            clearInterval(interval);
            onPaymentSuccess(generateReceipt());
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [method, onPaymentSuccess]);

  return (
    <div className="page-container">
      <div className="payment-card">
        <div className="order-badge">Order #{orderId}</div>
        <h2>Select Payment Method</h2>
        
        <div className="amount-due">
          <span>Amount Due</span>
          <strong className="amount-value">{total} AED</strong>
        </div>
        
        <div className="tabs">
          <button 
            className={`tab ${method === 'cc' ? 'active' : ''}`} 
            onClick={() => setMethod('cc')}
          >
            💳 Credit Card
          </button>
          <button 
            className={`tab ${method === 'qr' ? 'active' : ''}`} 
            onClick={() => setMethod('qr')}
          >
            📱 QR Scan
          </button>
        </div>

        <div className="payment-body">
          {method === 'cc' ? (
            <form onSubmit={(e) => { e.preventDefault(); onPaymentSuccess(generateReceipt()); }} className="cc-form">
              <input type="text" className="input-field" placeholder="Card Number (4242...)" required />
              <div className="row">
                <input type="text" className="input-field" placeholder="MM/YY" required />
                <input type="text" className="input-field" placeholder="CVC" required maxLength="3" />
              </div>
              <button type="submit" className="btn-primary full-width">Pay {total} AED</button>
            </form>
          ) : (
            <div className="qr-box">
              <div className="qr-image-placeholder">
                <img src={`https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${orderId}`} alt="QR" />
              </div>
              <p>Please scan with your banking app</p>
              <div className="timer-pill">Checking: {qrTimer}s</div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// ==========================================
// MAIN APP COMPONENT
// ==========================================
const App = () => {
  const [view, setView] = useState('selection'); 
  const [tickets, setTickets] = useState([]);
  const [cart, setCart] = useState({});
  const [total, setTotal] = useState(0);
  const [orderId, setOrderId] = useState(null);
  const [receipt, setReceipt] = useState(null);

  useEffect(() => {
    apiService.getTickets().then(setTickets);
  }, []);

  return (
    <div className="app-wrapper">
      {view === 'selection' && (
        <TicketSelection tickets={tickets} onProceed={(c, t) => { setCart(c); setTotal(t); setView('summary'); }} />
      )}

      {view === 'summary' && (
        <OrderSummaryPage 
          cart={cart} tickets={tickets} total={total} 
          onBack={() => setView('selection')}
          onCreateOrder={(id) => { setOrderId(id); setView('payment'); }} 
        />
      )}

      {view === 'payment' && (
        <PaymentPage 
          total={total} orderId={orderId} 
          onPaymentSuccess={(data) => { setReceipt(data); setView('success'); }} 
        />
      )}

      {view === 'success' && (
        <ThankYouPage receipt={receipt} cart={cart} tickets={tickets} onReset={() => setView('selection')} />
      )}
    </div>
  );
};

export default App;