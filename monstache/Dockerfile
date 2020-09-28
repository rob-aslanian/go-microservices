FROM rwynn/monstache:4.15.1

WORKDIR /bin/

ADD ./data/ .

RUN chmod +x init.sh

ENTRYPOINT ["sh", "init.sh"]
