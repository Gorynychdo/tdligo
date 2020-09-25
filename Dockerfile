FROM ubuntu:18.04

RUN apt-get -qq update &&  apt-get -qq install -y wget cmake build-essential gperf libssl-dev zlib1g-dev git

WORKDIR /srv/src
RUN wget -q https://golang.org/dl/go1.15.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.15.2.linux-amd64.tar.gz && \
    rm -rf /srv/src

WORKDIR /srv/src
RUN git clone https://github.com/tdlib/td.git --depth 1 && \
    cd td && mkdir build && cd build && \
    cmake -DCMAKE_BUILD_TYPE=Release .. && \
    cmake --build . -- -j3 && make install && \
    rm -rf /srv/src

COPY . /srv
WORKDIR /srv
ENV PATH=$PATH:/usr/local/go/bin
RUN GO111MODULE=on env CGO_ENABLED=1 go build -o build/service -mod=vendor ./cmd/service

EXPOSE 8000

CMD ["/srv/build/service"]
