# Step 1: Use an official Golang image as a build stage
FROM golang:1.21 AS build

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 4: Copy the rest of the application code
COPY . .

# Step 5: Build the Go app for the container
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Step 6: Use a smaller base image for the final container
FROM alpine:latest

# Step 7: Set the working directory for the final image
WORKDIR /app

# Step 8: Copy the compiled app from the build stage
COPY lib/resources/ /app/lib/resources/
COPY --from=build /app/main .
RUN touch /app/.env

RUN chmod -R 755 ./lib

# Step 9: Expose the port the app will run on
EXPOSE 8080

# Step 10: Run the app
CMD ["./main"]
