#include "LibMQTTConst.h"
#include "LibMQTTCallback.h"
#include "libmqtt.h"

int main(int argc, char *argv[]) {
    Publish("hello", goiiot_lib_mqtt_qos_0, "data", 4);
}
