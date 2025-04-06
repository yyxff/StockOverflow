#!/bin/bash

# 函数：计算文件字节长度并发送到服务器
send_to_server() {
    local file=$1
    local length=$(wc -c < "$file")
    (echo "$length"; cat "$file") | nc localhost 12345
    echo "\n"
    sleep 2
}

# 确保文件存在
check_files() {
    for file in 1_create.xml 2_buyorder1.xml 3_sellorder1.xml 4_multiorder.xml 5_query.xml 6_cancel.xml 7_mix.xml; do
        if [ ! -f "$file" ]; then
            echo "Error: File $file not found!"
            exit 1
        fi
    done
}

# 主函数
main() {
    echo "Starting test sequence...\n"
    check_files
    
    # 按顺序发送每个文件
    send_to_server "1_create.xml"
    send_to_server "2_buyorder1.xml"
    send_to_server "3_sellorder1.xml"
    send_to_server "4_multiorder.xml"
    send_to_server "5_query.xml"
    send_to_server "6_cancel.xml"
    send_to_server "7_mix.xml"
    
    echo "All files have been sent to the server."
}

# 执行主函数
main