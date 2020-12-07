FROM alpine:3.7
RUN mkdir /app 
COPY sync-problem /app
EXPOSE 8080:8080
ENTRYPOINT ["/app/sync-problem"]
