FROM centos:7
MAINTAINER liujiaqiang "liujiaqiang@126.com"
ENV container docker
ENV REFRESHED_AT "2018/01/26 14:50"
RUN rpm -Uvh http://nginx.org/packages/centos/7/noarch/RPMS/nginx-release-centos-7-0.el7.ngx.noarch.rpm
RUN yum -y update; yum -y install openssh-server nginx epel-release; yum -y install supervisor; yum clean all;
RUN mkdir -p /var/www/var/www/html /var/run/sshd /opt/logs
COPY config/global.conf /etc/nginx/conf.d/
COPY config/nginx.conf /etc/nginx/nginx.conf
COPY config/supervisord.ini /etc/supervisord.d/nginx.ini
EXPOSE 22 80
CMD ["/usr/bin/supervisord"]