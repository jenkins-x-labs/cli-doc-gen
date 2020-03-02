FROM centos:7
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/cli-doc-gen"]
COPY ./build/linux/cli-doc-gen /usr/local/bin/cli-doc-gen

