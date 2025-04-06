# locust/Dockerfile
FROM locustio/locust:2.34.0

# install lib
RUN pip3 install psutil

# copy script
COPY tcp_test.py /home/locust/locustfile.py

# set default run script
ENV LOCUSTFILE_PATH /home/locust/locustfile.py