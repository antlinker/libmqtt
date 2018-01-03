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

package cc.goiiot.libmqtt.example;

import cc.goiiot.libmqtt.Client;
import cc.goiiot.libmqtt.Callback;
import cc.goiiot.libmqtt.LogLevel;
import cc.goiiot.libmqtt.TopicMessageCallback;

public class Example {

    private static final String sTopicName = "foo";
    private static final String sTopicMsg = "bar";

    private static final String sClientID = "foo";
    private static final String sUsername = "foo";
    private static final String sPassword = "bar";

    private static final String sClientCert = "../../testdata/client-cert.pem";
    private static final String sClientKey = "../../testdata/client-key.pem";
    private static final String sCACert = "../../testdata/ca-cert.pem";
    private static final String sServerName = "MacBook-Air.local";

    public static void main(String[] args) throws InterruptedException {
        Client client = createClient();
        if (client == null) {
            return;
        }
        setCallback(client);
        client.connect();
        client.waitClient();
    }

    private static Client createClient() {
        try {
            return Client.newBuilder("localhost:8883")
                .setTLS(sClientCert, sClientKey, sCACert, sServerName, true)
                .setCleanSession(true)
                .setDialTimeout(10)
                .setKeepalive(10, 1.2)
                .setIdentity(sUsername, sPassword)
                .setLog(LogLevel.Verbose)
                .setSendBuf(100)
                .setRecvBuf(100)
                .setClientID(sClientID)
                .build();

        } catch (Exception e) {
            // Handle client create fail
            println("create client failed:", e.getMessage());
        }
        
        println("created client!");

        return null;
    }

    private static Client setCallback(Client client) {
        client.handle(sTopicName, new TopicMessageCallback(){
            public void onMessage(String topic, int qos, byte[] payload) {
                println("received message:", new String(payload), "from topic:", topic, "qos = " + qos);
            }
        });

        client.setCallback(new Callback() {
            public void onConnResult(boolean ok, String descp) {
                if (!ok) {
                    println("connection error:", descp);
                    return;
                }
                println("connected to server");
                client.subscribe(sTopicName, 0);
            }

            public void onLost(String descp) {
                println("connection lost, err:", descp);
            }

            public void onSubResult(String topic, boolean ok, String descp) {
                if (!ok) {
                    println("sub", topic, "failed:", descp);
                    return;
                }
                println("sub", topic, "success");
                client.publish(sTopicName, 0, sTopicMsg.getBytes());
            }

            public void onPubResult(String topic, boolean ok, String descp) {
                if (!ok) {
                    println("pub", topic, "failed:", descp);
                    return;
                }
                println("pub", topic, "success");
                client.unsubscribe(sTopicName);
            }

            public void onUnSubResult(String topic, boolean ok, String descp) {
                if (!ok) {
                    println("unsub", topic, "failed:", descp);
                    return;
                }
                println("unsub", topic, "success");
                client.destroy(true);
            }

            public void onPersistError(String descp) {
                println("persist err happened:", descp);
            }
        });

        return client;
    }

    private static void println(String... strings) {
        String result = String.join(" ", strings);
        System.out.println(result);
    }
}