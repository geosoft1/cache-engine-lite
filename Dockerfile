FROM golang
WORKDIR /app
COPY . /app
RUN go get github.com/gorilla/mux
RUN go build -o main .
EXPOSE 8080
ENV XAUTHTOKEN 35A6E
CMD ["./main"]