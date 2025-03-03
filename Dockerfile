FROM node:23-alpine3.20
COPY . /usr/wolite

# Install dependencies
WORKDIR /usr/wolite
RUN npm install

RUN chmod +x /usr/wolite/docker_entrypoint.sh
ENTRYPOINT [ "/usr/wolite/docker_entrypoint.sh" ]

EXPOSE 3000

# Start the app
CMD ["npm", "run", "start"]

