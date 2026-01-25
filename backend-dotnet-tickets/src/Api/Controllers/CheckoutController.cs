using Application.DTOs;
using Application.Services;
using Microsoft.AspNetCore.Mvc;
using System;
using System.Threading.Tasks;

namespace Api.Controllers
{
    [ApiController]
    [Route("api/v1/checkout")]
    public class CheckoutController : ControllerBase
    {
        private readonly ICheckoutService _checkoutService;

        public CheckoutController(ICheckoutService checkoutService)
        {
            _checkoutService = checkoutService;
        }

        [HttpPost("orders")]
        public async Task<IActionResult> CreateOrder([FromBody] CheckoutRequest request)
        {
            try
            {
                if (request.CartItems == null || request.CartItems.Count == 0)
                    return BadRequest("Cart is empty.");

                var result = await _checkoutService.ProcessCheckoutAsync(request);
                
                // Return 201 Created
                return StatusCode(201, result);
            }
            catch (Exception ex)
            {
                // In production, log this error
                return BadRequest(new { error = ex.Message });
            }
        }

        [HttpPost("confirm-payment")]
        public async Task<IActionResult> ConfirmPayment([FromBody] ConfirmPaymentReq request)
        {
            try
            {
                if (request.OrderId == null)
                    return BadRequest("Order ID is required.");

                var result = await _checkoutService.ProcessConfirmPaymentAsync(request);
                
                // Return 201 Created
                return StatusCode(201, result);
            }
            catch (Exception ex)
            {
                // In production, log this error
                return BadRequest(new { error = ex.Message });
            }
        }
    }
}