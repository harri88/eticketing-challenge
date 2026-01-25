using Application.Services;
using Domain.Interfaces;
using Infrastructure.Data;

var builder = WebApplication.CreateBuilder(args);

// 1. Register Services (Dependency Injection)
// Bind the Interface to the Concrete Implementation
builder.Services.AddScoped<ITicketRepository, DapperTicketRepository>();
builder.Services.AddScoped<ITicketService, TicketService>();
builder.Services.AddScoped<ICheckoutRepository, DapperCheckoutRepository>();
builder.Services.AddScoped<ICheckoutService, CheckoutService>();



builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowReactApp",
        policy =>
        {
            policy.WithOrigins("http://localhost:3000") // URL of your React App
                  .AllowAnyHeader()
                  .AllowAnyMethod();
        });
});

var app = builder.Build();

// 2. Configure Pipeline
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
    app.UseHttpsRedirection();
}

app.UseHttpsRedirection();
app.UseCors("AllowReactApp"); 
app.UseAuthorization();
app.MapControllers();

app.Run();