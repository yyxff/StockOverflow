from locust import User, task, events, between
from locust.env import Environment
import socket
import time
from jinja2 import Environment, FileSystemLoader
import random
from faker import Faker


env = Environment(loader=FileSystemLoader("templates"))
templates = {
    "template_create": env.get_template("create.xml.j2"),
    "template_transaction": env.get_template("transaction.xml.j2")
}


class TcpClient:
    def __init__(self, host, port, environment: Environment):
        self.host = host
        self.port = port
        self.sock = None
        self.fake = Faker()
        self.environment = environment

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
                self.environment.events.request.fire(
                    request_type="TCP",
                    name="send_xml",
                    response_time=int((time.time() - start_time) * 1000),
                    response_length=len(response),
                    exception=None,
                    context={"host": self.host}  # 可选上下文
            )
            else:
                raise Exception("no </results> in response")
        except Exception as e:
            # print("fail to send data:",e)
            self.environment.events.request.fire(
                    request_type="TCP",
                    name="send_xml",
                    response_time=int((time.time() - start_time) * 1000),
                    response_length=len(response),
                    exception=e,
                    context={"host": self.host}  # 可选上下文
            )


    def _generate_create_data(self):
        """生成 <create> 请求的测试数据"""
        # 生成 0~3 个 account 元素
        accounts = [
            {
                "id": self.fake.unique.bothify("ACCT_####"),  # 唯一账户ID
                "balance": round(random.uniform(1000.0, 100000.0), 2)  # 保留两位小数
            }
            for _ in range(random.randint(0, 3))
        ]
        
        # 记录生成的 account ID 到全局池
        # for acc in accounts:
        #     self.account_pool.append(acc["id"])

        # 生成 0~2 个 symbol 元素
        symbols = []
        for _ in range(random.randint(0, 2)):
            symbol = {
                "name": random.choice(["BTC", "ETH", "USD", "EUR"]),  # 符号名称
                "accounts": [
                    {"id": self.fake.unique.bothify("ACCT_####")}  # 从已有账户池中随机选择
                    for _ in range(random.randint(1, 3))  # 每个 symbol 包含 1~3 个 account
                ]
            }
            symbols.append(symbol)
        
        # 合并并打乱所有元素顺序
        all_elements = []
        all_elements.extend([{"type": "account", "data": a} for a in accounts])
        all_elements.extend([{"type": "symbol", "data": s} for s in symbols])
        random.shuffle(all_elements)  # 关键步骤：随机化顺序
        
        return {"elements": all_elements}

    def _generate_trans_data(self):
        """生成 <transactions> 请求的测试数据"""
        # 生成原始动作列表（未打乱顺序）
        action_types = ["order", "query", "cancel"]
        actions = []
        for _ in range(random.randint(1, 5)):
            action_type = random.choice(action_types)
            if action_type == "order":
                trans_id =  str(random.randint(1, 10000))
                actions.append({
                    "type": action_type,
                    "sym": self.fake.currency_code(),
                    "amount": random.randint(-100, 100),
                    "limit": round(random.uniform(10.0, 100.0), 2),
                    "trans_id": trans_id  # 生成新ID并记录
                })
                # self.trans_pool.append(trans_id)  # 添加到全局交易池
            else:
                # 从已有交易池中随机选择（确保即使顺序打乱也能引用有效ID）
                # trans_id = random.choice(self.trans_pool) if self.trans_pool else None
                trans_id = str(random.randint(1, 10000))
                if trans_id:
                    actions.append({
                        "type": action_type,
                        "trans_id": trans_id
                    })
        
        # 关键步骤：打乱动作顺序
        random.shuffle(actions)
        
        return {
            "account_id": self.fake.unique.bothify("ACCT_####"),
            "actions": actions
        }



class TcpUser(User):
    abstract = True
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.client = TcpClient(self.host, self.port, self.environment)
        self.client.connect()

    # @task
    # def send_request(self):
    #     xml_request = '<create><account id="1234567" balance="1000"/><symbol sym="SYM"><account id="1234567">233</account></symbol></create>' # 替换为真实 XML
    #     self.client.send_xml(xml_request)


    @task(weight=1)
    def send_create(self):
        xml_request = templates["template_create"].render(self.client._generate_create_data(),random=random)
        self.client.send_xml(xml_request)

    @task(weight=3)
    def send_transaction(self):
        xml_request = templates["template_transaction"].render(self.client._generate_trans_data())
        self.client.send_xml(xml_request)

        

class MyUser(TcpUser):
    host = "app"  # service name in docker
    port = 12345
    wait_time = between(1, 2)