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

package cc.goiiot.libmqtt;

import cc.goiiot.libmqtt.LogLevel;
import cc.goiiot.libmqtt.PersistMethod;

public class Client {

    private int mID;
    Callback mCallback;
    TopicMessageCallback mMainTopicCallback;

    // TODO: complete separate topic message callback
    // private Map<String, TopicMessageCallback> mTopicCallbacks;

    Client(int id) {
        mID = id;
        // mTopicCallbacks = new HashMap<>();
    }

    public void setCallback(Callback callback) {
        mCallback = callback;
    }

    public void connect() {
        LibMQTT._connect(mID);
    }

    public void waitClient() {
        LibMQTT._wait(mID);
    }

    public void handle(String topic, TopicMessageCallback callback) {
        if (callback != null) {
            LibMQTT._handle(mID, topic, callback);
            mMainTopicCallback = callback;
            // mTopicCallbacks.put(topic, callback);
        }
    }

    public void publish(String topic, int qos, byte[] payload) {
        LibMQTT._pub(mID, topic, qos, payload);
    }

    public void subscribe(String topic, int qos) {
        LibMQTT._sub(mID, topic, qos);
    }

    public void unsubscribe(String topic) {
        LibMQTT._unsub(mID, topic);
    }

    public void destroy(boolean force) {
        LibMQTT._destroy(mID, force);
    }

    public static Builder newBuilder(String server) throws Exception {
        int id = LibMQTT._newClient();
        if (id < 1) {
            throw new Exception("Can't create mqtt client");
        }
        LibMQTT._setServer(id, server);

        return new Builder(id);
    }

    public static class Builder {
        private int mID;

        private Builder(int id) {
            mID = id;
        }

        public Builder setCleanSession(boolean cleanSession) {
            LibMQTT._setCleanSession(mID, cleanSession);
            return this;
        }

        public Builder setKeepalive(int keepalive, double factor) {
            LibMQTT._setKeepalive(mID, keepalive, factor);
            return this;
        }

        public Builder setClientID(String clientID) {
            LibMQTT._setClientID(mID, clientID);
            return this;
        }

        public Builder setDialTimeout(int timeout) {
            LibMQTT._setDialTimeout(mID, timeout);
            return this;
        }

        public Builder setIdentity(String username, String password) {
            LibMQTT._setIdentity(mID, username, password);
            return this;
        }

        public Builder setLog(LogLevel level) {
            LibMQTT._setLog(mID, level.ordinal());
            return this;
        }

        public Builder setSendBuf(int size) {
            LibMQTT._setSendBuf(mID, size);
            return this;
        }

        public Builder setRecvBuf(int size) {
            LibMQTT._setRecvBuf(mID, size);
            return this;
        }

        public Builder setTLS(String cert, String key, String ca, String name, boolean skipVerify) {
            LibMQTT._setTLS(mID, cert, key, ca, name, skipVerify);
            return this;
        }

        public Builder setWill(String topic, int qos, boolean retain, byte[] payload) {
            LibMQTT._setWill(mID, topic, qos, retain, payload);
            return this;
        }

        public Builder setPersist(PersistMethod method) {
            // switch (method) {
            // case PersistMethod.None:
            //     // LibMQTT._setNonePersist(mID);
            //     break;
            // case PersistMethod.File:
            //     // LibMQTT._setFilePersist(mID);
            //     break;
            // case PersistMethod.Memory:
            //     // LibMQTT._setMemPersist(mID);
            //     break;
            // }
            return this;
        }

        public Client build() throws Exception {
            return LibMQTT.setupClient(mID);
        }
    }

}