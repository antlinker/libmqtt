#ifndef _LIB_MQTT_H_
#define _LIB_MQTT_H_

typedef int libmqtt_client_id;

typedef enum {
  libmqtt_connack_accepted = 0,
  libmqtt_connack_bad_proto = 1,
  libmqtt_connack_id_rejected = 2,
  libmqtt_connack_srv_unavail = 3,
  libmqtt_connack_bad_identity = 4,
  libmqtt_connack_auth_fail = 5,
} libmqtt_connack_t;

typedef enum {
  libmqtt_suback_max_qos0 = 0,
  libmqtt_suback_max_qos1 = 1,
  libmqtt_suback_max_qos2 = 2,
  libmqtt_suback_fail = 0x80,
} libmqtt_suback_t;

typedef enum {
  libmqtt_log_verbose = 0,
  libmqtt_log_debug = 1,
  libmqtt_log_info = 2,
  libmqtt_log_warning = 3,
  libmqtt_log_error = 4,
} libmqtt_log_level;

typedef void (*libmqtt_conn_handler)(char *server, libmqtt_connack_t code,
                                     char *err);

typedef void (*libmqtt_pub_handler)(char *topic, char *err);

typedef void (*libmqtt_sub_handler)(char *topic, int qos, char *err);

typedef void (*libmqtt_unsub_handler)(char *topic, char *err);

typedef void (*libmqtt_net_handler)(char *server, char *err);

typedef void (*libmqtt_topic_handler)(char *topic, int qos, char *msg,
                                      int size);

#endif // _LIB_MQTT_H_
