FROM gcr.io/distroless/static@sha256:3d0f463de06b7ddff27684ec3bfd0b54a425149d0f8685308b1fdf297b0265e9

ENTRYPOINT ["/app/clear-old-cache"]

COPY /dist/clear-old-cache /app/clear-old-cache
