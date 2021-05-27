FROM golang:1.15

RUN DEBIAN_FRONTEND=noninteractive apt-get install -y tzdata
ENV TZ Asia/Istanbul

#Create group and user named aybush, do not create home directory, do not assign password to user.
RUN groupadd aybush && useradd -m -g aybush aybush

WORKDIR /go/src/github.com/skarakasoglu/aybush-bot-api/
#Change owner and group of files as aybush.
ADD --chown=aybush:aybush . /go/src/github.com/skarakasoglu/aybush-bot-api/

#Give the owner to execute the run.sh shell script.
RUN chown aybush:aybush /go/src/github.com/skarakasoglu/aybush-bot-api/ && chmod 744 run.sh

#Log in as aybush.
#USER aybush

# is going to be used with twitch webhooks.
EXPOSE 8080:443

CMD ["/bin/bash", "-c", "/go/src/github.com/skarakasoglu/aybush-bot-api/run.sh"]