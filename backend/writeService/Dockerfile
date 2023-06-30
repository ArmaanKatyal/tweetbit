FROM node:16.20

WORKDIR /app

COPY . .

RUN yarn && yarn add typescript tsc ts-node && yarn build

EXPOSE 5002

CMD [ "yarn", "start" ]
