FROM dopplerhq/cli:3.54
ADD ./pubgolf-api-server .
ADD ./web-app-content/ ./web-app-content/
RUN chmod +x ./pubgolf-api-server

CMD ["run", "--", "./pubgolf-api-server"]
