#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "libmqtt.h"

// example log level for libmqtt log
#define LOG_LEVEL           libmqtt_log_silent

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
#define SSL_SERVER_NAME     "MacBook-Air.local"
#define SSL_SKIP_VERIFY     true

static int client;

// conn_handler for connect packet response
void conn_handler(char *server, libmqtt_connack_t code, char *err) {
  if (err != NULL) {
    printf("connect to server: %s failed, error = %s\n", server, err);
  } else if (code != libmqtt_connack_accepted) {
    printf("connect to server rejected, code: %d\n", code);
  } else {
    // connected to the server, subscribe some topic
    Libmqtt_subscribe(client, TOPIC, 0);
  }
}

// sub_handler for subscribe packet response, non-null err means topic sub failed
void sub_handler(char *topic, int qos, char *err) {
  if (err != NULL) {
    printf("sub failed topic: %s, error = %s\n", topic, err);
  } else {
    printf("sub success topic: %s, qos = %d\n", topic, qos);
    // subscribe example topic success, then publish some message to it
    Libmqtt_publish(client, TOPIC, 0, DATA, sizeof(DATA));
  }
}

// pub_handler for publish packet response, non-null err means `topic` publish failed
void pub_handler(char *topic, char *err) {
  if (err != NULL) {
    printf("pub topic: %s failed, error = %s\n", topic, err);
  } else {
    printf("pub topic: %s success\n", topic);
    // publish topic message succeeded, then unsubscribe that topic
    Libmqtt_unsubscribe(client, TOPIC);
  }
}

// unsub_handler for unsubscribe response, non-null err means `topic` unsubscribe failed
void unsub_handler(char *topic, char *err) {
  if (err != NULL) {
    printf("unsub topic: %s failed, error = %s\n", topic, err);
  } else {
    printf("unsub topic: %s success\n", topic);
    // unsubscribe example topic success, destroy and exit now
    Libmqtt_destroy(client, true);
  }
}

// net_handler for net connection information, called only when error happens
void net_handler(char *server, char *err) {
  printf("conn to server: %s broken, error = %s\n", server, err);
}

// topic_handler, just a example topic handler
void topic_handler(int client, char *topic, int qos, char *msg, int size) {
  printf("client: %d recv topic: %s, qos: %d, msg: %s\n", client, topic, qos, msg);
}

// set client options
void init() {
  client = Libmqtt_new_client();
  if (client < 1) {
    printf("create client failed\n");
    exit(1);
  }

  Libmqtt_client_set_server(client, SERVER);
  Libmqtt_client_set_log(client, LOG_LEVEL);
  Libmqtt_client_set_clean_session(client, CLEAN_SESSION);
  Libmqtt_client_set_keepalive(client, KEEPALIVE_INTERVAL, KEEPALIVE_FACTOR);
  Libmqtt_client_set_identity(client, USER, PASS);
  Libmqtt_client_set_client_id(client, CLIENT_ID);
  if (USE_TLS) {
    Libmqtt_client_set_tls(client, SSL_CLIENT_CERT, SSL_CLIENT_KEY, SSL_CA_CERT, SSL_SERVER_NAME, SSL_SKIP_VERIFY);
  }
}

int main(int argc, char *argv[]) {
  init();
  char *err = Libmqtt_setup(client);
  if (err != NULL) {
    printf("error happened when create client: %s\n", err);
    return 1;
  }

  Libmqtt_set_pub_handler(client, &pub_handler);
  Libmqtt_set_sub_handler(client, &sub_handler);
  Libmqtt_set_net_handler(client, &net_handler);
  Libmqtt_set_unsub_handler(client, &unsub_handler);
  Libmqtt_handle(client, TOPIC, &topic_handler);

  // connect to server
  Libmqtt_connect(client, &conn_handler);
  
  // wait until all connection exit
  Libmqtt_wait(client);
}