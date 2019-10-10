FROM alpine as alpine 
RUN addgroup -S lapkin && adduser -S lapkin -G lapkin
 
FROM scratch 
LABEL maintainer="alexeybezrukov2@gmail.com" 
WORKDIR /home/lapkins-api 
COPY --from=alpine /etc/passwd /etc/passwd
COPY lapkins-api . 
USER lapkin
EXPOSE 8081
ENTRYPOINT [ "./lapkins-api" ] 