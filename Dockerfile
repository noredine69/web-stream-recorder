FROM alpine:3.16

ARG BINARY_NAME_ARG
CMD echo BINARY_NAME_ARG ${BINARY_NAME_ARG}
ENV BINARY_NAME=$BINARY_NAME_ARG
CMD echo BINARY_NAME ${BINARY_NAME}
CMD apk update && apk upgrade && apk add bash
ENV APPDIR=/opt/${BINARY_NAME}
CMD mkdir -p ${APPDIR}

COPY build/$BINARY_NAME ${APPDIR}/

ARG USERNAME=running-user
RUN addgroup -g 1000 ${USERNAME}
RUN adduser --system -D -u 1000 -G ${USERNAME} ${USERNAME}
RUN chown -R ${USERNAME}: ${APPDIR}
USER ${USERNAME}

CMD ["/opt/web-stream-recorder/web-stream-recorder", "-config", "/opt/web-stream-recorder/config/config.json"]

