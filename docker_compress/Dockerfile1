FROM ubuntu:focal

ENV REDIS_VERSION=6.0.5
ENV REDIS_URL=http://download.redis.io/releases/redis-$REDIS_VERSION.tar.gz

# update source and install tools
RUN sed -i "s/archive.ubuntu.com/mirrors.aliyun.com/g; s/security.ubuntu.com/mirrors.aliyun.com/g" /etc/apt/sources.list 
RUN apt update 
RUN apt install -y curl make gcc

# download source code and install redis
RUN curl -L $REDIS_URL | tar xzv
WORKDIR redis-$REDIS_VERSION
RUN make
RUN make install

# clean up
RUN rm  -rf /var/lib/apt/lists/* 

CMD ["redis-server"]
