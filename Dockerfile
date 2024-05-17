FROM public.ecr.aws/lambda/provided:al2023 AS builder

RUN dnf -y update
RUN dnf install -y gzip 
RUN dnf install -y tar 
RUN dnf install -y gcc
RUN dnf -y install gcc-c++
RUN dnf install -y libstdc++-devel
ENV LD_LIBRARY_PATH /usr/lib64
# Install GoLang
RUN curl -LO https://golang.org/dl/go1.18.1.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz
RUN rm go1.18.1.linux-amd64.tar.gz

# Set environment variables for Go
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go

# Verify Go installation
RUN go version

WORKDIR /

COPY . .

RUN go build -o /output/built_file

FROM ubuntu
COPY --from=builder /output/built_file /built_file
ENTRYPOINT ["tail"]
CMD ["-f","/dev/null"]