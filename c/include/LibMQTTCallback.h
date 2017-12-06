#ifndef _GO_IIOT_LIB_MQTT_CALLBACK_
#define _GO_IIOT_LIB_MQTT_CALLBACK_

typedef enum {
    t,
} lib_mqtt_conn_code_t;

typedef void (*conn_callback)(lib_mqtt_conn_code_t);

#endif /*_GO_IIOT_LIB_MQTT_CALLBACK_*/