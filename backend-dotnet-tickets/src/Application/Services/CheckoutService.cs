using Application.DTOs;
using Domain.Entities;
using Domain.Interfaces;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace Application.Services
{
    public interface ICheckoutService
    {
        Task<CheckoutResponse> ProcessCheckoutAsync(CheckoutRequest request);
        Task<ConfirmPaymentResponse> ProcessConfirmPaymentAsync(ConfirmPaymentReq request);
    }

    public class CheckoutService : ICheckoutService
    {
        private readonly ICheckoutRepository _repository;

        public CheckoutService(ICheckoutRepository repository)
        {
            _repository = repository;
        }

        public async Task<CheckoutResponse> ProcessCheckoutAsync(CheckoutRequest request)
        {
            var orderId = "ORD-" + Guid.NewGuid().ToString("N").Substring(0, 10).ToUpper();
            var order = new Order
            {
                Id = orderId,
                CustomerEmail = request.CustomerEmail,
                Status = "created", // or pending_payment
                Currency = "AED",
                CreatedAt = DateTime.UtcNow
            };

            decimal totalAmount = 0;

            // 1. Loop through items to Validate & Calculate Logic
            foreach (var item in request.CartItems)
            {
                if (item.Quantity <= 0) continue;

                // A. Fetch freshly from DB (never trust frontend price)
                var ticket = await _repository.GetTicketByIdAsync(item.TicketId);

                if (ticket == null)
                    throw new Exception($"Ticket ID {item.TicketId} not found.");

                // B. Check Quota Availability
                if (ticket.Quota < item.Quantity)
                    throw new Exception($"Not enough quota for {ticket.Name}. Remaining: {ticket.Quota}");

                // C. Calculate Total (Server-side)
                decimal lineTotal = ticket.Price * item.Quantity;
                totalAmount += lineTotal;

                // D. Build Order Item
                order.Items.Add(new OrderItem
                {
                    TicketId = ticket.Id,
                    Quantity = item.Quantity,
                    PriceAtPurchase = ticket.Price // Snapshot the price
                });
            }

            order.TotalAmount = totalAmount;

            // 2. Persist: Save Order & Update Quotas in one transaction
            await _repository.CreateOrderWithReservationAsync(order);

            // 3. Return Response with Payment URL
            return new CheckoutResponse
            {
                OrderId = order.Id,
                Status = "pending_payment",
                Summary = new CheckoutSummary
                {
                    Subtotal = totalAmount,
                    Currency = "AED",
                    ItemCount = order.Items.Sum(i => i.Quantity)
                },
                RedirectUrl = $"https://payment-gateway.com/checkout/{order.Id}"
            };
        }

        public async Task<ConfirmPaymentResponse> ProcessConfirmPaymentAsync(ConfirmPaymentReq request)
        {
            var order = await _repository.GetOrderByIdAsync(request.OrderId);

            if (order == null)
                throw new Exception("Order not found.");

            if (order.Status != "created") 
                throw new Exception("Order is not in a payable state.");

            // Simulate payment confirmation logic
            order.Status = "paid";
            order.PaidAt = DateTime.UtcNow;

            await _repository.UpdateOrderAndTicketAsync(order);

            return new ConfirmPaymentResponse
            {
                OrderId = order.Id,
                Status = order.Status
            };
        }   
    }

    
}