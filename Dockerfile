FROM alpine

ADD jwtd /bin/jwtd
ADD jwtd-ctl/jwtd-ctl /bin/jwtd-ctl

EXPOSE 443

CMD ["/bin/jwtd" ]
