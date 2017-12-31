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
                public void onConnected() {
                    println("connected to server");
                    client.subscribe(sTopicName, 0);
                }

                public void onSubResult(String topic, String err){
                    if (err != null) {
                        println("sub", topic ,"failed:", err);
                    }
                    println("sub", topic ,"success");
                    client.publish(sTopicName, 0, sTopicMsg.getBytes());
                }
        
                public void onPubResult(String topic, String err){
                    if (err != null) {
                        println("pub", topic ,"failed:", err);
                    }
                    println("pub", topic ,"success");
                    client.unsubscribe(sTopicName);
                }
        
                public void onUnSubResult(String topic, String err){
                    if (err != null) {
                        println("unsub", topic ,"failed:", err);
                    }
                    println("unsub", topic ,"success");
                    client.destroy(true);
                }

                public void onPersistError(String err) {
                    println("persist err happened:", err);
                }
        
                public void onLost(String err) {
                    println("connection lost, err:", err);
                }
            });

            client.connect();
            Thread.sleep(100 *1000);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}