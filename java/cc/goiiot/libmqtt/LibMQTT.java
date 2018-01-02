package cc.goiiot.libmqtt;

import java.util.*;

public class LibMQTT {

    public static enum LogLevel {
        Silent, Verbose, Debug, Info, Warning, Error,
    }

    public static enum ConnAckCode {
        Accepted, Rejected, Unavaliable
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
        private Map<String, TopicMessageCallback> mTopicCallbacks;

        private Client(int id) throws Exception {
            mID = id;
            mTopicCallbacks = new HashMap<>();
        }

        public void setCallback(Callback callback) {
            mCallback = callback;
        }

        public void connect() {
            _conn(mID);
        }

        public void handle(String topic, TopicMessageCallback callback) {
            if (callback != null) {
                _handle(mID, topic, callback);
                mTopicCallbacks.put(topic, callback);
            }
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

            Client client = new Client(mID);
            sClientMap.put(mID, client);
            return client;
        }
    }

    public interface TopicMessageCallback {
        void onMessage(String topic, int qos, byte[] payload);
    }

    public static interface Callback {

        public void onConnResult(Exception e);
        
        public void onLost(Exception e);

        public void onSubResult(String topic, boolean ok);
        
        public void onPubResult(String topic, boolean ok);
        
        public void onUnSubResult(String topic, boolean ok);

        public void onPersistError(Exception e);
    }

    static {
        System.loadLibrary("mqtt-jni");
    }


    private static final Map<Integer, Client> sClientMap = new HashMap();

    private static void onConnMessage(int clientID, int code, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        if (err == null || "".equals(err)) {
            c.mCallback.onConnResult(null);
        }
        c.mCallback.onConnResult(new Exception(err));
    }

    private static void onNetMessage(int clientID, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        c.mCallback.onLost(new Exception(err));
    }

    private static void onPubMessage(int clientID, String topic, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        if (err == null || "".equals(err)) {
            c.mCallback.onPubResult(topic, true);
        }
        c.mCallback.onPubResult(topic, false);
    }

    private static void onSubMessage(int clientID, String topic, int qos, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        if (err == null || "".equals(err)) {
            c.mCallback.onSubResult(topic, true);
        }
        c.mCallback.onSubResult(topic, false);
    }

    private static void onUnsubMessage(int clientID, String topic, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        if (err == null || "".equals(err)) {
            c.mCallback.onUnSubResult(topic, true);
        }
        c.mCallback.onUnSubResult(topic, false);
    }

    private static void onPersistError(int clientID, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        c.mCallback.onPersistError(new Exception(err));
    }

    private static void onTopicMessage(int clientID, String topic, int qos, byte[] payload) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mTopicCallbacks == null) {
            return;
        }
        // TODO: deliver topic message to correct topic handler
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

    private static native String _setCallback(int id, Callback callback);

    private static native void _handle(int id, String topic, TopicMessageCallback  callback);

    private static native void _conn(int id);

    private static native void _pub(int id, String topic, int qos, byte[] payload);

    private static native void _sub(int id, String topic, int qos);

    private static native void _unsub(int id, String topic);

    private static native void _destroy(int id, boolean force);
}