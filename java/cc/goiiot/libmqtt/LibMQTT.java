package cc.goiiot.libmqtt;

public class LibMQTT {

    private static void onTopicMessage(int clientID, String topic, int qos, byte[] payload) {
        
    }

    public static enum LogLevel {
        Silent, Verbose, Debug, Info, Warning, Error,
    }

    public static Builder newBuilder(String server) throws Exception {
        int id = _newClient();
        if (id < 1) {
            throw new Exception("Can't create mqtt client");
        }
        _setServer(id, server);
        
        return new Builder(id);
    }

    public static class Client {
        private int mID;
        private Callback mCallback;

        private Client(int id) throws Exception {
            mID = id;
        }

        public void setCallback(Callback callback) {
            mCallback = callback;
        }

        public void connect() {
            _conn(mID);
        }

        public void handle(String topic, TopicMessageCallback callback) {
            _handle(mID, topic, callback);
        }

        public void publish(String topic, int qos, byte[] payload) {
            _pub(mID, topic, qos, payload);
        }

        public void subscribe(String topic, int qos) {
            _sub(mID, topic, qos);
        }

        public void unsubscribe(String topic) {
            _unsub(mID, topic);
        }

        public void destroy(boolean force) {
            _destroy(mID, force);
        }
    }

    public static class Builder {
        private int mID;

        private Builder(int id) {
            mID = id;
        }

        public Builder setCleanSession(boolean cleanSession) {
            _setCleanSession(mID, cleanSession);
            return this;
        }

        public Builder setKeepalive(int keepalive, double factor) {
            _setKeepalive(mID, keepalive, factor);
            return this;
        }

        public Builder setClientID(String clientID) {
            _setClientID(mID, clientID);
            return this;
        }

        public Builder setDialTimeout(int timeout) {
            _setDialTimeout(mID, timeout);
            return this;
        }

        public Builder setIdentity(String username, String password) {
            _setIdentity(mID, username, password);
            return this;
        }

        public Builder setLog(LogLevel level) {
            _setLog(mID, level.ordinal());
            return this;
        }

        public Builder setSendBuf(int size) {
            _setSendBuf(mID, size);
            return this;
        }

        public Builder setRecvBuf(int size) {
            _setRecvBuf(mID, size);
            return this;
        }

        public Builder setTLS(String cert, String key, String ca, String name, boolean skipVerify) {
            _setTLS(mID, cert, key, ca, name, skipVerify);
            return this;
        }

        public Builder setWill(String topic, int qos, boolean retain, byte[] payload) {
            _setWill(mID, topic, qos, retain, payload);
            return this;
        }

        public Client build() throws Exception {
            String s = _setup(mID);
            if (s != null) {
                throw new Exception(s);
            }

            return new Client(mID);
        }
    }

    public interface TopicMessageCallback {
        void onMessage(String topic, int qos, byte[] payload);
    }

    public static interface Callback {

        public void onConnected();

        public void onSubResult(String topic, String err);
        
        public void onPubResult(String topic, String err);
        
        public void onUnSubResult(String topic, String err);

        public void onPersistError(String err);
        
        public void onLost(String err);
    }

    static {
        System.loadLibrary("mqtt-jni");
    }

    private static native int _newClient();

    private static native void _setServer(int id, String server);

    private static native void _setCleanSession(int id, boolean flag);

    private static native void _setKeepalive(int id, int keepalive, double factor);

    private static native void _setClientID(int id, String clientID);

    private static native void _setDialTimeout(int id, int timeout);

    private static native void _setIdentity(int id, String username, String password);

    private static native void _setLog(int id, int logLevel);

    private static native void _setSendBuf(int id, int size);

    private static native void _setRecvBuf(int id, int size);

    private static native void _setTLS(int id, String cert, String key, String ca, String name, boolean skipVerify);

    private static native void _setWill(int id, String topic, int qos, boolean retain, byte[] payload);

    private static native String _setup(int id);

    private static native void _handle(int id, String topic, TopicMessageCallback  callback);

    private static native void _conn(int id);

    private static native void _pub(int id, String topic, int qos, byte[] payload);

    private static native void _sub(int id, String topic, int qos);

    private static native void _unsub(int id, String topic);

    private static native void _destroy(int id, boolean force);
}