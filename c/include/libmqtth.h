#include <string.h>
#include <stdbool.h>

#ifndef _LIB_MQTT_H_
#define _LIB_MQTT_H_

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
    libmqtt_suback_subfail = 0x80,
} libmqtt_suback_t;

typedef enum {
    libmqtt_log_verbose = 0,
    libmqtt_log_debug = 1,
    libmqtt_log_info = 2,
    libmqtt_log_warning = 3,
    libmqtt_log_error = 4,
} libmqtt_log_level;

typedef void (*libmqtt_conn_handler)
    (const char * server, libmqtt_connack_t code, const char * err);

typedef void (*libmqtt_pub_handler)
    (const char * topic, const char * err);

typedef void (*libmqtt_sub_handler)
    (const char * topic, int qos, const char * err);

typedef void (*libmqtt_unsub_handler)
    (const char * topic, const char * err);

typedef void (*libmqtt_net_handler)
    (const char * server, const char * err);

typedef void (*libmqtt_topic_handler)
    (const char * topic, int qos , const char * msg);

#endif // _LIB_MQTT_H_
