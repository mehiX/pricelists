FROM openjdk:21-jdk-buster

WORKDIR /h2

COPY . .
RUN ./downloadH2.sh

ENTRYPOINT [ "./runStandalone.sh" ]

EXPOSE 9092 8082 5435