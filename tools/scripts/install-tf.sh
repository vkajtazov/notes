#! /bin/sh
 
wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.14.0.tar.gz
sudo tar -C /usr/local -xzf libtensorflow-cpu-linux-x86_64-1.14.0.tar.gz
sudo ldconfig
rm -f libtensorflow-cpu-linux-x86_64-1.14.0.tar.gz