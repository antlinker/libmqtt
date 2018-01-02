package cc.goiiot.libmqtt;

public interface TopicMessageCallback {
    void onMessage(String topic, int qos, byte[] payload);
}