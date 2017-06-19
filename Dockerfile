FROM scratch
COPY hydroponics_pi /
COPY assets /assets
CMD ["/hydroponics_pi"]
## docker run -t -i --device=/dev/ttyACM0 -p 8080:8080 noodlensk/hydroponics