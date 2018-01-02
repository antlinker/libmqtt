package cc.goiiot.libmqtt.example;

import cc.goiiot.libmqtt.*;

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

    private static void println(String ...strings) {
        String result = String.join(" ", strings);
        System.out.println(result);
    }

    public static void main(String[] args) {
        try {
            LibMQTT.Client client = LibMQTT.newBuilder("localhost:8883")
                    .setCleanSession(true)
                    .setDialTimeout(10)
                    .setKeepalive(10, 1.2)
                    .setIdentity(sUsername, sPassword)
                    .setLog(LibMQTT.LogLevel.Verbose)
                    .setSendBuf(100).setRecvBuf(100)
                    .setTLS(sClientCert, sClientKey, sCACert, sServerName, true)
                    .setClientID(sClientID)
                    .build();

            client.setCallback(new LibMQTT.Callback() {
                public void onConnResult(Exception e) {
                    if (e != null) {
                        println("connection error:", e.getMessage());
                    }
                    println("connected to server");
                    client.subscribe(sTopicName, 0);
                }
        
                public void onLost(Exception e) {
                    println("connection lost, err:", e.getMessage());
                }

                public void onSubResult(String topic, boolean ok) {
                    if (!ok) {
                        println("sub", topic ,"failed");
                        return;
                    }
                    println("sub", topic ,"success");
                    client.publish(sTopicName, 0, sTopicMsg.getBytes());
                }
                
                public void onPubResult(String topic, boolean ok) {
                    if (!ok) {
                        println("pub", topic ,"failed");
                        return;
                    }
                    println("pub", topic ,"success");
                    client.unsubscribe(sTopicName);
                }
                
                public void onUnSubResult(String topic, boolean ok) {
                    if (!ok) {
                        println("unsub", topic ,"failed:");
                        return;
                    }
                    println("unsub", topic ,"success");
                    client.destroy(true);
                }

                public void onPersistError(Exception e) {
                    println("persist err happened:", e.getMessage());
                }
            });

            client.connect();
            Thread.sleep(100 *1000);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}