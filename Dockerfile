FROM scratch
COPY go-greenhouse_pi /
COPY assets /assets
CMD ["/go-greenhouse_pi"]
## docker run -t -i --device=/dev/ttyACM0 -p 8080:8080 noodlensk/go-greenhouse