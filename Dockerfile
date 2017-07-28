FROM alpine:3.5
MAINTAINER Frank Bille-Stauner <frank@cohousing.nu>

ADD cohousing-tenant-api /
ADD config.yml /

EXPOSE 8080

CMD ["/cohousing-tenant-api"]