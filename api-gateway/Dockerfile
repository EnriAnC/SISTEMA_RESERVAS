FROM devopsfaith/krakend:2.5

COPY krakend.json /etc/krakend/krakend.json

EXPOSE 8080

CMD [ "run", "-d", "-c", "/etc/krakend/krakend.json" ]
