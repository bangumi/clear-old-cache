FROM gcr.io/distroless/static@sha256:3592aa8171c77482f62bbc4164e6a2d141c6122554ace66e5cc910cadb961ff0

ENTRYPOINT ["/app/clear-old-cache"]

COPY /dist/clear-old-cache /app/clear-old-cache
