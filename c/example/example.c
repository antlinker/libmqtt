/*
 * Copyright GoIIoT (https://github.com/goiiot)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
void conn_handler(int client, char *server, libmqtt_connack_t code, char *err) {
  if (err != NULL) {
    printf("client: %d connect to server: %s failed, error = %s\n", client, server, err);
  } else if (code != libmqtt_connack_accepted) {
    printf("client: %d connect to server rejected, code: %d\n", client, code);
  } else {
    // connected to the server, subscribe some topic
    Libmqtt_subscribe(client, TOPIC, 0);
  }
}

// sub_handler for subscribe packet response, non-null err means topic sub failed
void sub_handler(int client, char *topic, int qos, char *err) {
  if (err != NULL) {
    printf("client: %d sub failed topic: %s, error = %s\n", client, topic, err);
  } else {
    printf("client: %d sub success topic: %s, qos = %d\n", client, topic, qos);
    // subscribe example topic success, then publish some message to it
    Libmqtt_publish(client, TOPIC, 0, DATA, sizeof(DATA));
  }
}

// pub_handler for publish packet response, non-null err means `topic` publish failed
void pub_handler(int client, char *topic, char *err) {
  if (err != NULL) {
    printf("client: %d pub topic: %s failed, error = %s\n", client, topic, err);
  } else {
    printf("client: %d pub topic: %s success\n", client, topic);
    // publish topic message succeeded, then unsubscribe that topic
    Libmqtt_unsubscribe(client, TOPIC);
  }
}

// unsub_handler for unsubscribe response, non-null err means `topic` unsubscribe failed
void unsub_handler(int client, char *topic, char *err) {
  if (err != NULL) {
    printf("client: %d unsub topic: %s failed, error = %s\n", client, topic, err);
  } else {
    printf("client: %d unsub topic: %s success\n", client, topic);
    // unsubscribe example topic success, destroy and exit now
    Libmqtt_destroy(client, true);
  }
}

// net_handler for net connection information, called only when error happens
void net_handler(int client, char *server, char *err) {
  printf("client: %d conn to server: %s broken, error = %s\n", client, server, err);
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
  printf("create client success, client = %d\n", client);

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
    printf("client: %d error happened when setup client: %s\n", client, err);
    return 1;
  }
  printf("setup client success\n");

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