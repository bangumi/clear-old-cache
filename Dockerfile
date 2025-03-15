FROM gcr.io/distroless/base-debian11

ENTRYPOINT ["/app/clear-old-cache"]

COPY /dist/clear-old-cache /app/clear-old-cache
