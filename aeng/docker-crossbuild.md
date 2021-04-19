
1. 启用 docker 的 buildx 特性

    更新 docker 配置文件 `~/.docker/config.json`，添加以下属性: 
    ```json
    {
      "experimental": "enabled"
    }
    ```

2. [启用主机的arm架构容器运行支持](https://github.com/multiarch/qemu-user-static#getting-started)

3. 创建 builder 实例

    ```shell
    docker buildx create --name multiarch --driver docker-container --use
    docker buildx inspect multiarch --bootstrap
    docker buildx ls # 确认当前使用的builder支持多平台
    ```

4. 构建镜像

    ```shell
    cat <<EOF > Dockerfile.me

    FROM golang:1.17.10-alpine3.15 as builder

    RUN apk update && apk add --no-cache git && \
        apk add --no-cache make && \
        apk add --no-cache bash && \
        apk add --no-cache gcc && \
        apk add --no-cache libc-dev && \
        apk add --no-cache binutils-gold

    WORKDIR $GOPATH/src/github.com/thanos-io/thanos
    RUN git clone -b v0.18.0 --depth 1 https://github.com/thanos-io/thanos.git .
    RUN git update-index --refresh
    RUN make build

    FROM quay.io/prometheus/busybox

    COPY --from=builder /go/bin/thanos /bin/thanos

    ENTRYPOINT [ "/bin/thanos" ]
    EOF

    # Makefile中的arm64修改为aarch64

    docker buildx build -f Dockerfile.me --output=type=registry --platform linux/amd64,linux/arm64 -t junotx/thanos:v0.18.0 .
    ```

参考：

> - [multiarch/qemu-user-static](https://github.com/multiarch/qemu-user-static): 用来启用支持容器的多架构运行
> - [mitchellh/gox](https://github.com/mitchellh/gox.git): golang交叉编译工具
> - [Golang 交叉编译中的那些坑](https://blog.csdn.net/Three_dog/article/details/94640507)
> - [docker manifest](https://docs.docker.com/engine/reference/commandline/manifest/)