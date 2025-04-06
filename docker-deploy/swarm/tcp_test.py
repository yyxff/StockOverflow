from locust import User, task, events
import socket
import time
from jinja2 import Template 

class TcpClient:
    def __init__(self, host, port):
        self.host = host
        self.port = port
        self.sock = None

    def connect(self):
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.connect((self.host, self.port))

    def send_xml(self, xml_data):
        start_time = time.time()
        try:
            # 发送 XML 数据（示例）
            xml_bytes = xml_data.encode(encoding="UTF-8")
            xml_lenght = len(xml_bytes)
            head=f"{xml_lenght}\n".encode(encoding="UTF-8")
            xml_send = head+xml_bytes
            self.sock.sendall(xml_send)

            # 接收响应（需处理粘包，例如读取到 </response> 结束）
            response = self.sock.recv(4096)
            if b"</results>" in response:  # 根据实际 XML 结束标记调整
                events.request_success.fire(
                    request_type="TCP",
                    name="send_xml",
                    response_time=int((time.time() - start_time) * 1000),
                    response_length=len(response),
                )
            else:
                events.request_failure.fire(
                    request_type="TCP",
                    name="send_xml",
                    response_time=int((time.time() - start_time) * 1000),
                    exception="Invalid Response",
                )
        except Exception as e:
            events.request_failure.fire(
                request_type="TCP",
                name="send_xml",
                response_time=int((time.time() - start_time) * 1000),
                exception=str(e),
            )

class TcpUser(User):
    abstract = True
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.client = TcpClient(self.host, self.port)
        self.client.connect()

    @task
    def send_request(self):
        xml_request = "<request>test</request>"  # 替换为真实 XML
        self.client.send_xml(xml_request)

class MyUser(TcpUser):
    host = "app"  # service name in docker
    port = 12345