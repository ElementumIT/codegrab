# This Dockerfile is optional for custom build steps
# The main build is handled by xgo in docker-compose.yml
FROM techknowlogick/xgo:latest
WORKDIR /workspace
COPY . /workspace
CMD ["/bin/bash"]
