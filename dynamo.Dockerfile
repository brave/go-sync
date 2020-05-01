ARG DB_LOCATION=/home/dynamodblocal/db
FROM amazon/dynamodb-local:1.12.0 AS install

USER root
RUN yum -y install awscli

USER dynamodblocal
ENV AWS_ACCESS_KEY_ID=#
ENV AWS_SECRET_ACCESS_KEY=#
ARG DB_LOCATION
ARG TABLE_NAME=client-entity-token-dev

COPY dynamo_local/ .
RUN mkdir -p ${DB_LOCATION} && \
      java -jar DynamoDBLocal.jar -sharedDb -dbPath ${DB_LOCATION} & \
      DYNAMO_PID=$! && \
      aws dynamodb create-table --cli-input-json file://table.json \
      --endpoint-url http://localhost:8000 --region us-west-2 && \
      kill $DYNAMO_PID

FROM amazon/dynamodb-local:1.12.0

ARG DB_LOCATION
COPY --chown=dynamodblocal:dynamodblocal --from=install ${DB_LOCATION} /db

CMD ["-jar", "DynamoDBLocal.jar", "-sharedDb", "-dbPath", "/db"]
