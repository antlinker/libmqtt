#include "cc_goiiot_libmqtt_LibMQTT.h"
#include "libmqtt.h"

static JavaVM *jvm;

static jclass libmqtt_class;

static jmethodID on_conn_msg_mid;
static jmethodID on_sub_msg_mid;
static jmethodID on_pub_msg_mid;
static jmethodID on_unsub_msg_mid;
static jmethodID on_net_msg_mid;
static jmethodID on_persist_err_mid;
static jmethodID on_topic_msg_mid;

void conn_handler(int client, char *server, libmqtt_connack_t code, char *err) {
  JNIEnv *g_env;
  (*jvm)->AttachCurrentThread(jvm, (void **)&g_env, NULL);

  if (g_env == NULL || libmqtt_class == NULL || on_conn_msg_mid == 0) {
    return;
  } else {
    (*g_env)->CallStaticVoidMethod(g_env, libmqtt_class,
                                   on_conn_msg_mid, client,
                                   0,
                                   (*g_env)->NewStringUTF(g_env, err));
  }
}

void sub_handler(int client, char *topic, int qos, char *err) {
  JNIEnv *g_env;
  (*jvm)->AttachCurrentThread(jvm, (void **)&g_env, NULL);

  if (g_env == NULL || libmqtt_class == NULL || on_sub_msg_mid == 0) {
    return;
  } else {
    (*g_env)->CallStaticVoidMethod(g_env, libmqtt_class,
                                   on_sub_msg_mid, client,
                                   (*g_env)->NewStringUTF(g_env, topic),
                                   qos,
                                   (*g_env)->NewStringUTF(g_env, err));
  }
}

void pub_handler(int client, char *topic, char *err) {
  JNIEnv *g_env;
  (*jvm)->AttachCurrentThread(jvm, (void **)&g_env, NULL);

  if (g_env == NULL || libmqtt_class == NULL || on_pub_msg_mid == 0) {
    return;
  } else {
    (*g_env)->CallStaticVoidMethod(g_env, libmqtt_class,
                                   on_pub_msg_mid, client,
                                   (*g_env)->NewStringUTF(g_env, topic),
                                   (*g_env)->NewStringUTF(g_env, err));
  }
}

void unsub_handler(int client, char *topic, char *err) {
  JNIEnv *g_env;
  (*jvm)->AttachCurrentThread(jvm, (void **)&g_env, NULL);

  if (g_env == NULL || libmqtt_class == NULL || on_unsub_msg_mid == 0) {
    return;
  } else {
    (*g_env)->CallStaticVoidMethod(g_env, libmqtt_class,
                                   on_unsub_msg_mid, client,
                                   (*g_env)->NewStringUTF(g_env, topic),
                                   (*g_env)->NewStringUTF(g_env, err));
  }
}

void net_handler(int client, char *server, char *err) {
  JNIEnv *g_env;
  (*jvm)->AttachCurrentThread(jvm, (void **)&g_env, NULL);

  if (g_env == NULL || libmqtt_class == NULL || on_net_msg_mid == 0) {
    return;
  } else {
    (*g_env)->CallStaticVoidMethod(g_env, libmqtt_class,
                                   on_net_msg_mid, client,
                                   (*g_env)->NewStringUTF(g_env, err));
  }
}

void persist_handler(int client, char *err) {
  JNIEnv *g_env;
  (*jvm)->AttachCurrentThread(jvm, (void **)&g_env, NULL);

  if (g_env == NULL || libmqtt_class == NULL || on_topic_msg_mid == 0) {
    return;
  } else {
    (*g_env)->CallStaticVoidMethod(g_env, libmqtt_class,
                                   on_persist_err_mid, client,
                                   (*g_env)->NewStringUTF(g_env, err));
  }
}

void topic_handler(int client, char *topic, int qos, char *payload, int size) {
  JNIEnv *g_env;
  (*jvm)->AttachCurrentThread(jvm, (void **)&g_env, NULL);

  if (g_env == NULL || libmqtt_class == NULL || on_topic_msg_mid == 0) {
    return;
  } else {
    jbyteArray result = (*g_env)->NewByteArray(g_env, size);
    (*g_env)->SetByteArrayRegion(g_env, result, 0, size, (jbyte *)payload);

    (*g_env)->CallStaticVoidMethod(g_env, libmqtt_class,
                                   on_topic_msg_mid, client,
                                   (*g_env)->NewStringUTF(g_env, topic),
                                   qos, result);
  }
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _init
 * Signature: ()V
 */
JNIEXPORT void JNICALL
Java_cc_goiiot_libmqtt_LibMQTT__1init
(JNIEnv *env, jclass c) {

  (*env)->GetJavaVM(env, &jvm);
  libmqtt_class = (jclass) (*env)->NewGlobalRef(env, (*env)->FindClass(env, "cc/goiiot/libmqtt/LibMQTT"));;

  on_conn_msg_mid = (*env)->GetStaticMethodID(env, c,
                    "onConnMessage", "(IILjava/lang/String;)V");
  on_net_msg_mid = (*env)->GetStaticMethodID(env, c,
                   "onNetMessage", "(ILjava/lang/String;)V");
  on_pub_msg_mid = (*env)->GetStaticMethodID(env, c,
                   "onPubMessage", "(ILjava/lang/String;Ljava/lang/String;)V");
  on_sub_msg_mid = (*env)->GetStaticMethodID(env, c,
                   "onSubMessage", "(ILjava/lang/String;ILjava/lang/String;)V");
  on_unsub_msg_mid = (*env)->GetStaticMethodID(env, c,
                     "onUnsubMessage", "(ILjava/lang/String;Ljava/lang/String;)V");
  on_persist_err_mid = (*env)->GetStaticMethodID(env, c,
                       "onPersistError", "(ILjava/lang/String;)V");
  on_topic_msg_mid = (*env)->GetStaticMethodID(env, c,
                     "onTopicMessage", "(ILjava/lang/String;I[B)V");
}

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
(JNIEnv *env, jclass c, jint id, jstring server) {

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
(JNIEnv *env, jclass c, jint id, jboolean flag) {

  Libmqtt_client_set_clean_session(id, flag);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setKeepalive
 * Signature: (IID)V
 */
JNIEXPORT void JNICALL
Java_cc_goiiot_libmqtt_LibMQTT__1setKeepalive
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
(JNIEnv *env, jclass c, jint id, jstring client_id) {

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
(JNIEnv *env, jclass c, jint id, jint timeout) {

  Libmqtt_client_set_dial_timeout(id, timeout);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setIdentity
 * Signature: (ILjava/lang/String;Ljava/lang/String;)V
 */
JNIEXPORT void JNICALL
Java_cc_goiiot_libmqtt_LibMQTT__1setIdentity
(JNIEnv *env, jclass c, jint id, jstring username, jstring password) {

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
(JNIEnv *env, jclass c, jint id, jint log_level) {

  Libmqtt_client_set_log(id, log_level);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setSendBuf
 * Signature: (II)V
 */
JNIEXPORT void JNICALL
Java_cc_goiiot_libmqtt_LibMQTT__1setSendBuf
(JNIEnv *env, jclass c, jint id, jint size) {

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
(JNIEnv *env, jclass c, jint id, jstring cert,
 jstring key, jstring ca, jstring srv_name, jboolean skip_verify) {

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
(JNIEnv *env, jclass c, jint id, jstring topic,
 jint qos, jboolean retain, jbyteArray payload) {

  const char *c_topic = (*env)->GetStringUTFChars(env, topic, 0);
  jsize len = (*env)->GetArrayLength(env, payload);
  jbyte *body = (*env)->GetByteArrayElements(env, payload, 0);

  Libmqtt_client_set_will(id, c_topic, qos, retain, body, len);

  (*env)->ReleaseStringUTFChars(env, topic, c_topic);
  (*env)->ReleaseByteArrayElements(env, payload, body, 0);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setNonePersist
 * Signature: (I)V
 */
JNIEXPORT void JNICALL
Java_cc_goiiot_libmqtt_LibMQTT__1setNonePersist
(JNIEnv *env, jclass c, jint id) {

  Libmqtt_client_set_none_persist(id);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setMemPersist
 * Signature: (IIZZ)V
 */
JNIEXPORT void JNICALL
Java_cc_goiiot_libmqtt_LibMQTT__1setMemPersist
(JNIEnv *env, jclass c, jint id, jint max_count,
 jboolean ex_drop, jboolean dup_replace) {

  Libmqtt_client_set_mem_persist(id, max_count, ex_drop, dup_replace);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _setFilePersist
 * Signature: (ILjava/lang/String;IZZ)V
 */
JNIEXPORT void JNICALL
Java_cc_goiiot_libmqtt_LibMQTT__1setFilePersist
(JNIEnv *env, jclass c, jint id, jstring dir,
 jint max_count, jboolean ex_drop, jboolean dup_replace) {

  const char *c_dir = (*env)->GetStringUTFChars(env, dir, 0);

  Libmqtt_client_set_file_persist(id, c_dir, max_count, ex_drop, dup_replace);

  (*env)->ReleaseStringUTFChars(env, dir, c_dir);
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

  Libmqtt_set_pub_handler(id, &pub_handler);
  Libmqtt_set_sub_handler(id, &sub_handler);
  Libmqtt_set_net_handler(id, &net_handler);
  Libmqtt_set_unsub_handler(id, &unsub_handler);
  Libmqtt_set_persist_handler(id, &persist_handler);

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

  if (callback == NULL) {
    return;
  }

  const char *c_topic = (*env)->GetStringUTFChars(env, topic, 0);

  Libmqtt_handle(id, c_topic, &topic_handler);

  (*env)->ReleaseStringUTFChars(env, topic, c_topic);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _conn
 * Signature: (I)V
 */
JNIEXPORT void JNICALL
Java_cc_goiiot_libmqtt_LibMQTT__1connect
(JNIEnv *env, jclass c, jint id) {

  Libmqtt_connect(id, &conn_handler);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _wait
 * Signature: (I)V
 */
JNIEXPORT void JNICALL
Java_cc_goiiot_libmqtt_LibMQTT__1wait
(JNIEnv *env, jclass c, jint id) {

  (*env)->MonitorEnter(env, c);
  Libmqtt_wait(id);
  (*env)->MonitorExit(env, c);
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

  Libmqtt_publish(id, c_topic, qos, body, len);

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

  Libmqtt_subscribe(id, c_topic, qos);

  (*env)->ReleaseStringUTFChars(env, topic, c_topic);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _unsub
 * Signature: (ILjava/lang/String;)V
 */
JNIEXPORT void JNICALL
Java_cc_goiiot_libmqtt_LibMQTT__1unsub
(JNIEnv *env, jclass c, jint id, jstring topic) {

  const char *c_topic = (*env)->GetStringUTFChars(env, topic, 0);

  Libmqtt_unsubscribe(id, c_topic);

  (*env)->ReleaseStringUTFChars(env, topic, c_topic);
}

/*
 * Class:     cc_goiiot_libmqtt_LibMQTT
 * Method:    _destroy
 * Signature: (IZ)V
 */
JNIEXPORT void JNICALL
Java_cc_goiiot_libmqtt_LibMQTT__1destroy
(JNIEnv *env, jclass c, jint id, jboolean force) {

  Libmqtt_destroy(id, force);
}