FROM node:12.13.0

COPY ./proto ./proto
COPY ./web-app ./web-app
WORKDIR ./web-app
RUN npm install -g npx && \
    npm install && \
    npm run sapper:build

CMD ["node", "__sapper__/build/index.js"]
