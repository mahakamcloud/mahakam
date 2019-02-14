# This file creates standard build environment for building Mahakam
FROM mahakamcloud/maha-cross:v0.1.0-1

ADD dist/bin/mahakam-server-linux/mahakam-server /mahakam-server
RUN chmod +x /mahakam-server

ENTRYPOINT ["/mahakam-server"]
CMD ["--config", "/etc/mahakam/config.yaml"]