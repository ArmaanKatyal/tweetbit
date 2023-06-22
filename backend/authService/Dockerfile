FROM node:16.20-alpine

WORKDIR /app

COPY . .

RUN yarn && yarn add typescript tsc ts-node && yarn build

EXPOSE 5001

CMD [ "yarn", "start" ]
