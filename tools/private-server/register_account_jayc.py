import time
import random
import string
import json

import requests

register_url = "http://103.70.76.250/register/api.php"

if __name__ == '__main__':
    invite_code = input("请输入邀请码：")
    for i in range(10):
        letters = string.ascii_lowercase
        number = string.digits
        username = ''.join(random.choice(letters) for i in range(8))
        password = ''.join(random.choice(number) for i in range(8))
        data = {
            "account": username,
            "password": password,
            "invitecode": invite_code,
        }
        print(data)
        response = requests.post(
            register_url,
            json=data,
            params={"do": "1"},
            headers={
                "Content-Type": "application/json; charset=UTF-8",
                "Accept": "application/json, text/javascript, */*; q=0.01",
                "Accept-Encoding": "gzip, deflate",
                "Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
                "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
            }
        )
        print(response.status_code)
        body = response.text
        print(body)
        res_body = json.loads(body)
        print(res_body)
        if response.status_code == 200 and res_body.get("code") != 2:
            print("注册成功")
            csv_file = open("account.csv", "a")
            csv_file.write('\n')
            csv_file.write(username + "," + password)
            csv_file.close()
            data = {
                "account": username,
                "password": password,
            }
            response = requests.post(
                register_url,
                json=data,
                params={"do": "2"},
                headers={
                    "Content-Type": "application/json; charset=UTF-8",
                    "Accept": "application/json, text/javascript, */*; q=0.01",
                    "Accept-Encoding": "gzip, deflate",
                    "Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
                    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
                }
            )
            print(response.status_code)
            body = response.text
            print(body)
        else:
            print("注册失败")
        time.sleep(1)
