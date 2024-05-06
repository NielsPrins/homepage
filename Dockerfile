FROM alpine AS build

WORKDIR /usr/src/app

RUN apk add bash go npm

COPY . .

RUN chmod +x build.sh && ./build.sh --linux-only

FROM alpine

WORKDIR /usr/src/app

COPY --from=build /usr/src/app/dist .

EXPOSE 80
CMD ["./homepage-linux"]