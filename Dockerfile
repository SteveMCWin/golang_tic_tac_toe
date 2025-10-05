FROM golang:1.23

# get dependencies
RUN apt-get update && apt-get install -y \
    sqlite3 \
    libsqlite3-dev \
    build-essential \
    file \
    libpcre3-dev \
    libpcre3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# get the sqlite extension
RUN mkdir -p extensions
ADD https://raw.githubusercontent.com/sqlite/sqlite/master/ext/misc/spellfix.c ./
RUN gcc -fPIC -shared -o extensions/spellfix.so spellfix.c -lsqlite3 -lpcre
RUN rm spellfix.c

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 GOOS=linux go build -o tic_tac_toe.fun

EXPOSE 5000

CMD ["./tic_tac_toe.fun"]
