FROM node:16.15.1-alpine AS development
WORKDIR /app
COPY package.json /app/package.json
COPY package-lock.json /app/package-lock.json
RUN npm ci
COPY . /app
CMD [ "npm", "start" ]

FROM development AS prod
ENV NODE_ENV production
RUN addgroup -g 1001 appuser && adduser -u 1001 -G appuser -s /bin/sh -D appuser
USER appuser
RUN npm install -g pm2
CMD ["pm2-runtime", "start", "npm", "--", "start"]