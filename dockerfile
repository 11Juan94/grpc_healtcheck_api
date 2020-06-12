FROM golang:alpine
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
WORKDIR /build
RUN go get -u github.com/gorilla/mux
#RUN go mod init github.com/11Juan94/grpc_healtcheck_api
COPY . .
# Build the application
RUN go build -o main .
# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist
# Copy binary from build to main folder
RUN cp /build/main .
# Export necessary port
EXPOSE 80
# Command to run when starting the container
CMD ["/dist/main"]