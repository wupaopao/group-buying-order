FROM ccr.ccs.tencentyun.com/qcloud/centos
WORKDIR /service/
ENV HOSTNAME ${NAME}
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
RUN mkdir -p /service/
COPY ./build/bin /service/

RUN mkdir -p /service/conf/
COPY ./conf_release /service/conf/

RUN mkdir -p /service/data
RUN mkdir -p /service/data/excel
COPY ./data/excel /service/data/excel

CMD ["/bin/sh"]
ENTRYPOINT "/service/bin"