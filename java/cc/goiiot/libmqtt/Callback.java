package cc.goiiot.libmqtt;

public interface Callback {

    public void onSubResult(String topic, boolean ok, String description);
    
    public void onPubResult(String topic, boolean ok, String description);
    
    public void onUnSubResult(String topic, boolean ok, String description);

    public void onConnResult(boolean ok, String description);
    
    public void onLost(String description);

    public void onPersistError(String description);
}