FROM rabbitmq:3-management
ENV RABBITMQ_DEFAULT_USER=guest
ENV RABBITMQ_DEFAULT_PASS=guest
RUN apt-get update
RUN apt-get install -y netcat
EXPOSE 15672 5672