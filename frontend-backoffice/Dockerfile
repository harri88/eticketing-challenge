# Stage 1: Build
FROM node:18-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY . .

# Pass environment variables during build
ARG REACT_APP_TICKET_SERVICE
ARG REACT_APP_PAYMENT_SERVICE
ARG REACT_APP_LEDGER_SERVICE
ENV REACT_APP_TICKET_SERVICE=$REACT_APP_TICKET_SERVICE
ENV REACT_APP_PAYMENT_SERVICE=$REACT_APP_PAYMENT_SERVICE
ENV REACT_APP_LEDGER_SERVICE=$REACT_APP_LEDGER_SERVICE


RUN echo "Building with REACT_APP_TICKET_SERVICE=${REACT_APP_TICKET_SERVICE}"
RUN echo "Building with REACT_APP_PAYMENT_SERVICE=${REACT_APP_PAYMENT_SERVICE}"
RUN echo "Building with REACT_APP_LEDGER_SERVICE=${REACT_APP_LEDGER_SERVICE}"   
RUN npm run build

# Stage 2: Serve
FROM node:18-alpine

WORKDIR /app

# Install serve to run the React app
RUN npm install -g serve

COPY --from=builder /app/build ./build

EXPOSE 3000

CMD ["serve", "-s", "build", "-l", "3000"]
