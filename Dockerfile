# Docker para desarrollo
FROM golang:1.22.6-bullseye

WORKDIR /go/src/github.com/nmarsollier/authgo

ENV MONGO_URL=mongodb://host.docker.internal:27017
ENV RABBIT_URL=amqp://host.docker.internal
ENV FLUENT_URL=host.docker.internal:24224

# Puerto de Auth Service y debug
EXPOSE 3000

# Just a terminal, manual mode
# CMD ["bash"]

# To run in debug mode
CMD ["go" , "run" , "/go/src/github.com/nmarsollier/authgo"]
