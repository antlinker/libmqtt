from ctypes import cdll, c_float

mqtt = cdll.LoadLibrary('libmqtt.so')

class Client(object):
    '''
    Modern MQTT Client build on top of libmqtt for Python
    '''
    def __init__(self, server=None, dial_timeout=10, clean_session=True, keepalive=60, keepalive_factor=1.5, client_id=None, username=None, password=None, client_cert=None, client_key=None, ca_cert=None, server_name=None, skip_verify=True, will_topic=None, will_retain=None, will_qos=None, will_msg=None):
        '''
        Setup the client with options
        '''
        if server is None:
            raise ValueError('Server should not be None')
        cid = mqtt.Libmqtt_new_client()
        if cid < 1:
            raise RuntimeError("Couldn't create client")
        self.cid = cid
        mqtt.Libmqtt_client_set_server(cid, server.encode('utf-8'))
        mqtt.Libmqtt_client_set_clean_session(cid, clean_session)
        if client_cert is str:
            mqtt.Libmqtt_client_set_tls(cid, client_cert.encode('utf-8'), \
            client_key.encode('utf-8'), ca_cert.encode('utf-8'), \
            server_name.encode('utf-8'), skip_verify)
        if client_id is str:
            mqtt.Libmqtt_client_set_client_id(cid, client_id.encode('utf-8'))
        if will_topic is str:
            mqtt.Libmqtt_client_set_will(cid, will_topic.encode('utf-8'), \
            will_retain, will_qos, will_msg)
        mqtt.Libmqtt_client_set_dial_timeout(cid, dial_timeout)
        mqtt.Libmqtt_client_set_keepalive(cid, keepalive, c_float(keepalive_factor))
        if client_id is not None:
            mqtt.Libmqtt_client_set_client_id(cid, client_id.encode('utf-8'))
        if username is not None:
            mqtt.Libmqtt_client_set_identity(cid, username.encode('utf-8'), None)
            if password is not None:
                mqtt.Libmqtt_client_set_identity(cid, username.encode('utf-8'), \
                password.encode('utf-8'))

        err = mqtt.Libmqtt_setup(cid)
        if err is not None and err != 0:
            raise ValueError(err)
    def set_callback(self, listener=None):
        '''
        Set client lifecycle listener
        '''
        pass
    def connect(self, handler=None):
        '''
        Connect to server
        '''
        mqtt.Libmqtt_conn(self.cid, handler)
    def handle(self, topic=None, handler=None):
        '''
        Register handler for topic(s)
        '''
        if topic is not str:
            raise ValueError("Topic should be string")
        mqtt.Libmqtt_handle(self.cid, topic, handler)
    def subscribe(self, topic=None, qos=0):
        '''
        Subscribe one topic
        '''
        if topic is not str or qos is not int:
            raise ValueError("Topic should be string, qos should be int")
        mqtt.Libmqtt_sub(self.cid, topic, qos)
    def publish(self, topic=None, qos=0, retain=False, msg=None):
        '''
        Publish one topic message
        '''
        if topic is not str or qos is not int \
            or retain is not bool or msg is not bytes:
            raise ValueError("Topic should be string, qos should be int, \
            retain should be bool, msg should be bytes")
        mqtt.Libmqtt_pub(self.cid, topic, qos, retain, msg)
    def unsubscribe(self, topic=None):
        '''
        Unsubscribe topic(s)
        '''
        if topic is not str:
            raise ValueError('Topic should be string')
        mqtt.Libmqtt_unsub(self.cid, topic)
    def destroy(self, force=False):
        '''
        Destroy the client
        '''
        if force is not bool:
            raise ValueError('Force should be bool')
        mqtt.Libmqtt_destroy(self.cid, force)
