FROM node:16.15.1-alpine AS development
WORKDIR /app
COPY package.json .
RUN npm i
COPY . .
CMD ["npm", "run", "dev"]