yum search java-1.8.0-openjdk
yum install -y java-1.8.0-openjdk-devel
javac com/Sum.java
jar cvfm sum.jar Manifest.MF com/Sum.class
java -jar sum.jar 10 abc-def
