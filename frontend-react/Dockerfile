# Stage 1: Build
FROM node:18-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY . .

# Pass environment variable during build
ARG REACT_APP_API_URL
# Debug: check if var is set
RUN echo "DEBUG: REACT_APP_API_URL is '${REACT_APP_API_URL}'"
RUN if [ -z "$REACT_APP_API_URL" ]; then echo "ERROR: REACT_APP_API_URL is missing"; exit 1; fi

ENV REACT_APP_API_URL=$REACT_APP_API_URL

RUN echo "Building with REACT_APP_API_URL=${REACT_APP_API_URL}"
RUN npm run build

# Stage 2: Serve
FROM node:18-alpine

WORKDIR /app

# Install serve to run the React app
RUN npm install -g serve

COPY --from=builder /app/build ./build

EXPOSE 3000

CMD ["serve", "-s", "build", "-l", "3000"]
