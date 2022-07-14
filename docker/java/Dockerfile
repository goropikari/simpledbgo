FROM openjdk:18.0.1.1-jdk-buster

COPY ./entrypoint.sh /
WORKDIR /app

# RUN wget http://www.cs.bc.edu/%7Esciore/simpledb/SimpleDB_3.4.zip && \
#     unzip SimpleDB_3.4.zip && \
#     cd SimpleDB_3.4 && \
#     sed -i -e "1i package simpleclient;" simpleclient/SimpleIJ.java

ENTRYPOINT bash /entrypoint.sh
