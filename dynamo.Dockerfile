FROM amazon/dynamodb-local:2.6.1

USER dynamodblocal
ENV AWS_ACCESS_KEY_ID=GOSYNC
ENV AWS_SECRET_ACCESS_KEY=GOSYNC

COPY schema/dynamodb/ /schema/

VOLUME ["/db"]

USER root
RUN dnf -y install awscli util-linux \
 && dnf clean all \
 && rm -rf /var/cache/dnf
COPY dynamo.entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["-jar", "DynamoDBLocal.jar", "-sharedDb", "-dbPath", "/db"]
