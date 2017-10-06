FROM rmoriz/mvg-live

ADD proxymvg /opt/proxymvg

EXPOSE 2121

ENTRYPOINT ["/opt/proxymvg", "--listenport", "2121"]