#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "libmqtth.h"
#include "libmqtt.h"

#define LOG_LEVEL           libmqtt_log_verbose

#define SERVER              "localhost:1883"
#define USER                "username"
#define PASS                "password"
#define CLIENT_ID           "client_id"

#define CLEAN_SESSION       true
#define KEEPALIVE_INTERVAL  10
#define KEEPALIVE_FACTOR    1.2

#define TOPIC               "foo"
#define DATA                "bar"

#define USE_TLS             false
#define SSL_CLIENT_CERT     "client_cert.pem"
#define SSL_CLIENT_KEY      "client_key.pem"
#define SSL_CA_CERT         "ca_cert.pem"
#define SSL_SERVER_NAME     "foo.bar"
#define SSL_SKIP_VERIFY     true

void conn_handler(char *server, libmqtt_connack_t code,
                  char *err) {
  if (err != NULL) {
    printf("connect to server: %s failed, error = %s\n", server, err);
  } else if (code != libmqtt_connack_accepted) {
    printf("connect to server rejected, code: %d\n", code);
  } else {
    Subscribe(TOPIC, 0);
  }

  free(server);
  free(err);
}

void pub_handler(char *topic, char *err) {
  if (err != NULL) {
    printf("publish topic: %s failed, error = %s\n", topic, err);
  } else {
    printf("publish topic: %s success\n", topic);
    UnSubscribe(TOPIC);
  }

  free(topic);
  free(err);
}

void sub_handler(char *topic, int qos, char *err) {
  if (err != NULL) {
    printf("subscribe failed topic: %s, error = %s\n", topic, err);
  } else {
    printf("subscribe success topic: %s, qos = %d\n", topic, qos);
    Publish(TOPIC, 0, DATA, sizeof(DATA));
    // wait some time for response
    sleep(1);
  }

  free(topic);
  free(err);
}

void unsub_handler(char *topic, char *err) {
  if (err != NULL) {
    printf("unsub topic: %s failed, error = %s\n", topic, err);
  } else {
    printf("unsub topic: %s success\n", topic);
    Destroy(true);
  }

  free(topic);
  free(err);
}

void net_handler(char *server, char *err) {
  printf("connection to server: %s broken, error = %s\n", server, err);

  free(server);
  free(err);
}

void topic_handler(char *topic, int qos, char *msg, int size) {
  printf("received topic: %s, qos: %d, msg: %s\n", topic, qos, msg);

  free(topic);
  free(msg);
}

void init() {
  SetServer(SERVER);
  SetLog(LOG_LEVEL);
  SetCleanSession(CLEAN_SESSION);
  SetKeepalive(KEEPALIVE_INTERVAL, KEEPALIVE_FACTOR);
  SetIdentity(USER, PASS);
  SetClientID(CLIENT_ID);
  if (USE_TLS) {
    SetTLS(SSL_CLIENT_CERT, SSL_CLIENT_KEY, SSL_CA_CERT, SSL_SERVER_NAME, SSL_SKIP_VERIFY)
  }
}

int main(int argc, char *argv[]) {
  init();
  char *err = SetUp();

  if (err != NULL) {
    printf("error happened when create client: %s\n", err);
  }

  // register client handlers
  SetConnHandler(&conn_handler);
  SetPubHandler(&pub_handler);
  SetSubHandler(&sub_handler);
  SetUnSubHandler(&unsub_handler);
  SetNetHandler(&net_handler);

  // the topic handler
  Handle(TOPIC, &topic_handler);

  // connect to server
  Conn();
  
  // wait until all connection exit
  Wait();
}