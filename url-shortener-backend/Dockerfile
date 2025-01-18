# Step 1: Use a lightweight Go image
FROM golang:1.23-alpine

# Step 2: Set the working directory
WORKDIR /app

# Step 3: Copy Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 4: Copy the source code
COPY . .

# Step 5: Build the application
RUN go build -o url-shortener

# Step 6: Expose the application port
EXPOSE 8080

# Step 7: Command to run the application
CMD ["./url-shortener"]
