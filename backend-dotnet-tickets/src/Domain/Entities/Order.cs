using System;
using System.Collections.Generic;

namespace Domain.Entities
{
    public class Order
    {
        public string Id { get; set; }
        public string CustomerEmail { get; set; }
        public decimal TotalAmount { get; set; }
        public string Currency { get; set; }
        public string Status { get; set; } // 'pending_payment', 'created'

        public DateTime PaidAt { get; set; }
        public DateTime CreatedAt { get; set; }
        public List<OrderItem> Items { get; set; } = new List<OrderItem>();
    }

    public class OrderItem
    {
        public string TicketId { get; set; }
        public int Quantity { get; set; }
        public decimal PriceAtPurchase { get; set; }
    }
}