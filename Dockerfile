FROM golang:1.20 as build
WORKDIR /workspace/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o boson -v ./cmd/boson

FROM scratch
COPY --from=build /workspace/boson /boson
COPY --from=build /workspace/secrets/credentials.json /secrets/credentials.json
COPY --from=subtlepseudonym/healthcheck:0.1.1 /healthcheck /healthcheck

EXPOSE 9000/tcp
HEALTHCHECK --interval=60s --timeout=2s --retries=3 --start-period=2s \
	CMD ["/healthcheck", "localhost:9000", "/health"]

CMD ["/boson"]
