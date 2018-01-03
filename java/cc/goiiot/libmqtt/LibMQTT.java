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

import java.util.Map;
import java.util.HashMap;

import cc.goiiot.libmqtt.Utils;
import cc.goiiot.libmqtt.Client;

final class LibMQTT {

    static Client setupClient(int cid) throws Exception {
        String s = LibMQTT._setup(cid);
        if (!Utils.isEmpty(s)) {
            throw new Exception(s);
        }
        Client c = new Client(cid);
        sClientMap.put(cid, c);
        return c;
    }

    private static final Map<Integer, Client> sClientMap = new HashMap<>();

    private static void onConnMessage(int clientID, int code, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        if (!Utils.isEmpty(err)) {
            c.mCallback.onConnResult(false, err);
            return;
        }
        c.mCallback.onConnResult(true, "ok");
    }

    private static void onNetMessage(int clientID, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        c.mCallback.onLost(err);
    }

    private static void onPubMessage(int clientID, String topic, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        if (!Utils.isEmpty(err)) {
            c.mCallback.onPubResult(topic, false, err);
            return;
        }
        c.mCallback.onPubResult(topic, true, "ok");
    }

    private static void onSubMessage(int clientID, String topic, int qos, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        if (!Utils.isEmpty(err)) {
            c.mCallback.onSubResult(topic, false, err);
            return;
        }
        c.mCallback.onSubResult(topic, true, "ok");
    }

    private static void onUnsubMessage(int clientID, String topic, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        if (!Utils.isEmpty(err)) {
            c.mCallback.onUnSubResult(topic, false, err);
            return;
        }
        c.mCallback.onUnSubResult(topic, true, "ok");
    }

    private static void onPersistError(int clientID, String err) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mCallback == null) {
            return;
        }
        c.mCallback.onPersistError(err);
    }

    private static void onTopicMessage(int clientID, String topic, int qos, byte[] payload) {
        Client c = sClientMap.get(clientID);
        if (c == null || c.mMainTopicCallback == null) {
            return;
        }
        c.mMainTopicCallback.onMessage(topic, qos, payload);
        // TODO: deliver topic message to correct topic handler
    }

    static {
        System.loadLibrary("mqtt-jni");
        _init();
    }

    private static native void _init();

    private static native String _setup(int id);

    static native int _newClient();

    static native void _setServer(int id, String server);

    static native void _setCleanSession(int id, boolean flag);

    static native void _setKeepalive(int id, int keepalive, double factor);

    static native void _setClientID(int id, String clientID);

    static native void _setDialTimeout(int id, int timeout);

    static native void _setIdentity(int id, String username, String password);

    static native void _setLog(int id, int logLevel);

    static native void _setSendBuf(int id, int size);

    static native void _setRecvBuf(int id, int size);

    static native void _setTLS(int id, String cert, String key, String ca, String name, boolean skipVerify);

    static native void _setWill(int id, String topic, int qos, boolean retain, byte[] payload);

    static native void _setNonePersist(int id);

    static native void _setMemPersist(int id, int maxCount, boolean dropOnExceed, boolean duplicateReplace);

    static native void _setFilePersist(int id, String dirPath, int maxCount, boolean dropOnExceed, boolean duplicateReplace);

    static native void _handle(int id, String topic, TopicMessageCallback callback);

    static native void _connect(int id);

    static native void _wait(int id);

    static native void _pub(int id, String topic, int qos, byte[] payload);

    static native void _sub(int id, String topic, int qos);

    static native void _unsub(int id, String topic);

    static native void _destroy(int id, boolean force);
}