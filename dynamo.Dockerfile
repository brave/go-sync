ARG DB_LOCATION=/home/dynamodblocal/db
FROM amazon/dynamodb-local:1.12.0 AS install

USER root
RUN yum -y install awscli

USER dynamodblocal
ENV AWS_ACCESS_KEY_ID=#
ENV AWS_SECRET_ACCESS_KEY=#
ARG AWS_ENDPOINT=http://localhost:8000
ARG AWS_REGION=us-west-2
ARG DB_LOCATION
ARG TABLE_NAME=client-entity-dev

COPY schema/dynamodb/ .
RUN mkdir -p ${DB_LOCATION} && \
      java -jar DynamoDBLocal.jar -sharedDb -dbPath ${DB_LOCATION} & \
      DYNAMO_PID=$! && \
      aws dynamodb create-table --cli-input-json file://table.json \
      --endpoint-url ${AWS_ENDPOINT} --region ${AWS_REGION} && \
      kill $DYNAMO_PID

FROM amazon/dynamodb-local:1.12.0

ARG DB_LOCATION
COPY --chown=dynamodblocal:dynamodblocal --from=install ${DB_LOCATION} /db

CMD ["-jar", "DynamoDBLocal.jar", "-sharedDb", "-dbPath", "/db"]
