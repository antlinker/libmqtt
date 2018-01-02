package cc.goiiot.libmqtt;

final class Utils {

    static boolean isEmpty(String s) {
        return s == null || "".equals(s) || "null".equals(s);
    }
}