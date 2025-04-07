# locust/Dockerfile
FROM locustio/locust:2.34.0

# install lib
RUN pip3 install psutil faker

# copy script
COPY tcp_test.py /home/locust/locustfile.py

RUN mkdir -p /home/locust/templates
COPY templates/*.j2 /home/locust/templates/

# set default run script
ENV LOCUSTFILE_PATH /home/locust/locustfile.py