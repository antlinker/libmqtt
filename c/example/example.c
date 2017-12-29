#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "libmqtth.h"
#include "libmqtt.h"

// example log level for libmqtt log
#define LOG_LEVEL           libmqtt_log_verbose

// example server address
#define SERVER              "localhost:8883"

// example connect config
#define USER                "username"
#define PASS                "password"
#define CLIENT_ID           "client_id"
#define CLEAN_SESSION       true
#define KEEPALIVE_INTERVAL  10
#define KEEPALIVE_FACTOR    1.2

// example topic and message
#define TOPIC               "foo"
#define DATA                "bar"

// example tls config for emqttd
#define USE_TLS             true
#define SSL_CLIENT_CERT     "../../testdata/client-cert.pem"
#define SSL_CLIENT_KEY      "../../testdata/client-key.pem"
#define SSL_CA_CERT         "../../testdata/ca-cert.pem"
#define SSL_SERVER_NAME     "../../testdata/MacBook-Air.local"
#define SSL_SKIP_VERIFY     true

// conn_handler for connect packet response
void conn_handler(char *server, libmqtt_connack_t code,
                  char *err) {
  if (err != NULL) {
    printf("connect to server: %s failed, error = %s\n", server, err);
  } else if (code != libmqtt_connack_accepted) {
    printf("connect to server rejected, code: %d\n", code);
  } else {
    // connected to the server, subscribe some topic
    Subscribe(TOPIC, 0);
  }

  free(server);
  free(err);
}

// sub_handler for subscribe packet response, non-null err means topic sub failed
void sub_handler(char *topic, int qos, char *err) {
  if (err != NULL) {
    printf("subscribe failed topic: %s, error = %s\n", topic, err);
  } else {
    printf("subscribe success topic: %s, qos = %d\n", topic, qos);
    // subscribe example topic success, then publish some message to it
    Publish(TOPIC, 0, DATA, sizeof(DATA));
    // wait some time for response
    sleep(1);
  }

  free(topic);
  free(err);
}

// pub_handler for publish packet response, non-null err means `topic` publish failed
void pub_handler(char *topic, char *err) {
  if (err != NULL) {
    printf("publish topic: %s failed, error = %s\n", topic, err);
  } else {
    printf("publish topic: %s success\n", topic);
    // publish topic message succeeded, then try unsubscribe that topic
    UnSubscribe(TOPIC);
  }

  free(topic);
  free(err);
}

// unsub_handler for unsubscribe response, non-null err means `topic` unsubscribe failed
void unsub_handler(char *topic, char *err) {
  if (err != NULL) {
    printf("unsub topic: %s failed, error = %s\n", topic, err);
  } else {
    printf("unsub topic: %s success\n", topic);
    // unsubscribe example topic success, destroy and exit now
    Destroy(true);
  }

  free(topic);
  free(err);
}

// net_handler for net connection information, called only when error happens
void net_handler(char *server, char *err) {
  printf("connection to server: %s broken, error = %s\n", server, err);

  free(server);
  free(err);
}

// topic_handler, just a example topic handler
void topic_handler(char *topic, int qos, char *msg, int size) {
  printf("received topic: %s, qos: %d, msg: %s\n", topic, qos, msg);

  free(topic);
  free(msg);
}

// set client options
void init() {
  SetServer(SERVER);
  SetLog(LOG_LEVEL);
  SetCleanSession(CLEAN_SESSION);
  SetKeepalive(KEEPALIVE_INTERVAL, KEEPALIVE_FACTOR);
  SetIdentity(USER, PASS);
  SetClientID(CLIENT_ID);
  if (USE_TLS) {
    SetTLS(SSL_CLIENT_CERT, SSL_CLIENT_KEY, SSL_CA_CERT, SSL_SERVER_NAME, SSL_SKIP_VERIFY);
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