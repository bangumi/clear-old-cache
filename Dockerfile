FROM gcr.io/distroless/static@sha256:9197324ba51d9cd071af8505989365c006adf9d6d2067eada25aef00abbb5278

ENTRYPOINT ["/app/clear-old-cache"]

COPY /dist/clear-old-cache /app/clear-old-cache
