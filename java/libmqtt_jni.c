#include "cc_goiiot_libmqtt_LibMQTT.h"
#include "libmqtt.h"

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _newClient
 * Signature: ()I
 */
JNIEXPORT jint JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1newClient
    (JNIEnv *env, jclass c) {

    return Libmqtt_new_client();
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setServer
 * Signature: (ILjava/lang/String;)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1setServer
    (JNIEnv * env, jclass c, jint id, jstring server) {
    
    const char *srv = (*env)->GetStringUTFChars(env, server, 0);
    
    Libmqtt_client_set_server(id, srv);
    
    (*env)->ReleaseStringUTFChars(env, server, srv);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setCleanSession
 * Signature: (IZ)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1setCleanSession
    (JNIEnv * env, jclass c, jint id, jboolean flag) {
    
    Libmqtt_client_set_clean_session(id, flag);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setKeepalive
 * Signature: (IID)V
 */
JNIEXPORT void JNICALL Java_cc_goiiot_libmqtt_LibMQTT__1setKeepalive
  (JNIEnv *env, jclass c, jint id, jint keepalive, jdouble factor) {

      Libmqtt_client_set_keepalive(id, keepalive, factor);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setClientID
 * Signature: (ILjava/lang/String;)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1setClientID
    (JNIEnv * env, jclass c, jint id, jstring client_id) {

    const char *cid = (*env)->GetStringUTFChars(env, client_id, 0);
    
    Libmqtt_client_set_client_id(id, cid);

    (*env)->ReleaseStringUTFChars(env, client_id, cid);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setDialTimeout
 * Signature: (II)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1setDialTimeout
    (JNIEnv * env, jclass c, jint id, jint timeout) {

    Libmqtt_client_set_dial_timeout(id, timeout);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setIdentity
 * Signature: (ILjava/lang/String;Ljava/lang/String;)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1setIdentity
    (JNIEnv * env, jclass c, jint id, jstring username, jstring password) {
    
    const char *user = (*env)->GetStringUTFChars(env, username, 0);
    const char *pass = (*env)->GetStringUTFChars(env, password, 0);
    
    Libmqtt_client_set_identity(id, user, pass);

    (*env)->ReleaseStringUTFChars(env, username, user);
    (*env)->ReleaseStringUTFChars(env, password, pass);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setLog
 * Signature: (II)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1setLog
    (JNIEnv * env, jclass c, jint id, jint log_level) {

    Libmqtt_client_set_log(id, log_level);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setSendBuf
 * Signature: (II)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1setSendBuf
    (JNIEnv * env, jclass c, jint id, jint size) {

    Libmqtt_client_set_send_buf(id, size);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setRecvBuf
 * Signature: (II)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1setRecvBuf
    (JNIEnv *env, jclass c, jint id, jint size) {

    Libmqtt_client_set_send_buf(id, size);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setTLS
 * Signature: (ILjava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Z)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1setTLS
    (JNIEnv * env, jclass c, jint id, jstring cert, jstring key, jstring ca, jstring srv_name, jboolean skip_verify) {
    
    const char *c_cert = (*env)->GetStringUTFChars(env, cert, 0);
    const char *c_key = (*env)->GetStringUTFChars(env, key, 0);
    const char *c_ca = (*env)->GetStringUTFChars(env, ca, 0);
    const char *c_srv = (*env)->GetStringUTFChars(env, srv_name, 0);

    Libmqtt_client_set_tls(id, c_cert, c_key, c_ca, c_srv, skip_verify);

    (*env)->ReleaseStringUTFChars(env, cert, c_cert);
    (*env)->ReleaseStringUTFChars(env, key, c_key);
    (*env)->ReleaseStringUTFChars(env, ca, c_ca);
    (*env)->ReleaseStringUTFChars(env, srv_name, c_srv);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setWill
 * Signature: (ILjava/lang/String;IZ[B)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1setWill
    (JNIEnv *env, jclass c, jint id, jstring topic, jint qos, jboolean retain, jbyteArray payload) {

    const char *c_topic = (*env)->GetStringUTFChars(env, topic, 0);
    jsize len = (*env)->GetArrayLength(env, payload);
    jbyte *body = (*env)->GetByteArrayElements(env, payload, 0);

    Libmqtt_client_set_will(id, c_topic, qos, retain, body, len);

    (*env)->ReleaseStringUTFChars(env, topic, c_topic);
    (*env)->ReleaseByteArrayElements(env, payload, body, 0);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setup
 * Signature: (I)Ljava/lang/String;
 */
JNIEXPORT jstring JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1setup
    (JNIEnv *env, jclass c, jint id) {
    
    char *err = Libmqtt_setup(id);
    if (err != NULL) {
        return (*env)->NewStringUTF(env, err);
    }

    return NULL;
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _handle
 * Signature: (ILjava/lang/String;Lcc/goiiot/libmqtt/LibMQTT/TopicMessageCallback;)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1handle
    (JNIEnv *env, jclass c, jint id, jstring topic, jobject callback) {

    const char *c_topic = (*env)->GetStringUTFChars(env, topic, 0);
    
    Libmqtt_handle(id, c_topic, NULL);

    (*env)->ReleaseStringUTFChars(env, topic, c_topic);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _conn
 * Signature: (I)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1conn
    (JNIEnv *env, jclass c, jint id) {
    
    Libmqtt_conn(id, NULL);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _pub
 * Signature: (ILjava/lang/String;I[B)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1pub
    (JNIEnv *env, jclass c, jint id, jstring topic, jint qos, jbyteArray payload) {

    const char *c_topic = (*env)->GetStringUTFChars(env, topic, 0);
    jsize len = (*env)->GetArrayLength(env, payload);
    jbyte *body = (*env)->GetByteArrayElements(env, payload, 0);

    Libmqtt_pub(id, c_topic, qos, body, len);

    (*env)->ReleaseStringUTFChars(env, topic, c_topic);
    (*env)->ReleaseByteArrayElements(env, payload, body, 0);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _sub
 * Signature: (ILjava/lang/String;I)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1sub
    (JNIEnv *env, jclass c, jint id, jstring topic, jint qos) {

    const char *c_topic = (*env)->GetStringUTFChars(env, topic, 0);
    Libmqtt_sub(id, c_topic, qos);
    (*env)->ReleaseStringUTFChars(env, topic, c_topic);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _unsub
 * Signature: (ILjava/lang/String;)V
 */
JNIEXPORT void JNICALL 
Java_cc_goiiot_libmqtt_LibMQTT__1unsub
    (JNIEnv *env, jclass c, jint id, jstring topic){

    const char *c_topic = (*env)->GetStringUTFChars(env, topic, 0);
    Libmqtt_unsub(id, c_topic);
    (*env)->ReleaseStringUTFChars(env, topic, c_topic);
}