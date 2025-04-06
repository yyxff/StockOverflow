import xml.etree.ElementTree as ET
import xml.dom.minidom as minidom
import socket

def build_create_xml():
    """交互式构建<create>类型的XML结构"""
    root = ET.Element("create")
    
    while True:
        choice = input("Add 'account', 'symbol' or 'done'? ").strip().lower()
        if choice == 'done':
            break
        elif choice == 'account':
            acc_id = input("Account ID: ").strip()
            balance = input("Balance: ").strip()
            ET.SubElement(root, "account", id=acc_id, balance=balance)
        elif choice == 'symbol':
            sym = input("Symbol name: ").strip()
            symbol_elem = ET.SubElement(root, "symbol", sym=sym)
            
            while True:
                acc_id = input(f"  Add account ID under {sym} (or 'done'): ").strip()
                if acc_id == 'done':
                    break
                num = input(f"  Shares for {acc_id}: ").strip()
                acc_elem = ET.SubElement(symbol_elem, "account", id=acc_id)
                acc_elem.text = num
        else:
            print("Invalid choice")
    
    return root

def build_transaction_xml():
    """交互式构建<transactions>类型的XML结构"""
    acc_id = input("Account ID for transaction: ").strip()
    root = ET.Element("transactions", id=acc_id)
    
    while True:
        action = input("Add 'order', 'query', 'cancel' or 'done'? ").strip().lower()
        if action == 'done':
            break
        elif action == 'order':
            sym = input("Symbol: ").strip()
            amount = input("Amount: ").strip()
            limit = input("Limit: ").strip()
            ET.SubElement(root, "order", sym=sym, amount=amount, limit=limit)
        elif action == 'query':
            trans_id = input("Query ID: ").strip()
            ET.SubElement(root, "query", id=trans_id)
        elif action == 'cancel':
            trans_id = input("Cancel ID: ").strip()
            ET.SubElement(root, "cancel", id=trans_id)
        else:
            print("Invalid action")
    
    return root

def send_xml(xml_root, host, port):
    """将XML结构转换为字节流并通过TCP发送"""
    # 生成带XML声明的完整字节流
    xml_data = ET.tostring(xml_root, encoding='UTF-8')
    xml_with_declaration = b'<?xml version="1.0" encoding="UTF-8"?>\n' + xml_data
    # print(xml_with_declaration)
    xml_length = len(xml_with_declaration)
    header = f"{xml_length}\n".encode('utf-8')  # 标头格式：数字+换行
    full_data = header + xml_with_declaration        

    dom = minidom.parseString(xml_with_declaration)
    pretty_xml = dom.toprettyxml(indent='  ', encoding='UTF-8')
    print(pretty_xml.decode(encoding='UTF-8'))

    # 建立TCP连接并发送
    try:
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.connect((host, port))
            s.sendall(full_data)
            print(f"Sent {len(full_data)} bytes to {host}:{port}")

            # Read the response from the server
            response = s.recv(4096)  # Adjust buffer size as needed
            response_str = response.decode('utf-8')
            print("Response from server:")
            print(response_str)
            
            # Write the response to a local file
            with open("server_response.xml", "w", encoding="utf-8") as f:
                f.write(response_str)
    except Exception as e:
        print(f"Connection failed: {str(e)}")

if __name__ == "__main__":
    # 选择XML类型
    while 1:
        xml_type = input("Create which XML (create/transaction)? ").strip().lower()
        
        # 构建对应XML结构
        if xml_type == 'create':
            xml_root = build_create_xml()
        elif xml_type == 'transaction':
            xml_root = build_transaction_xml()
        elif xml_type == 'exit':
            exit(0)
        else:
            print("Invalid XML type")
            continue
                
        # 获取服务器信息
        # host = input("Server IP (default localhost): ").strip() or 'localhost'
        # port = int(input("Server port (default 12345): ").strip() or 12345
        host="127.0.0.1"
        port=12345
        
        # 发送数据
        send_xml(xml_root, host, port)