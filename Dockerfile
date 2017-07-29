FROM rmoriz/mvg-live

ADD proxymvg /opt/proxymvg

ENTRYPOINT ["/opt/proxymvg", "--listenport", "2233"]