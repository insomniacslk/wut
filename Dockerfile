from golang:1.22

LABEL BUILD="docker build -t insomniacslk/wut -f Dockerfile ."
LABEL RUN="docker run --rm -it insomniacslk/wut"

WORKDIR /app

ADD . .

RUN cd cmd/wut && go build
RUN mv cmd/wut/wut /app
RUN mv cmd/wut/acronyms.json.example /app/acronyms.json

ENTRYPOINT ["/app/wut", "-f", "/app/acronyms.json"]
