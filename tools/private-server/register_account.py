import time
import random
import string

import requests

register_url = "http://182.92.182.33:9090/gm/user/gameUserRegister"

if __name__ == '__main__':
    for i in range(95):
        letters = string.ascii_lowercase
        number = string.digits
        username = ''.join(random.choice(letters) for i in range(8))
        password = ''.join(random.choice(number) for i in range(8))
        data = {
            "account": username,
            "password": password,
            "shareUser": username,
        }
        print(data)
        response = requests.post(
            register_url,
            data=data,
            params={"PromotionCode": ""},
            headers={
                "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
                "Accept": "application/json, text/javascript, */*; q=0.01",
                "Accept-Encoding": "gzip, deflate",
                "Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
                "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
                "X-Requested-With": "XMLHttpRequest",
            }
        )
        print(response.text)
        if response.status_code == 200:
            print("注册成功")
            csv_file = open("account.csv", "a")
            csv_file.write('\n')
            csv_file.write(username + "," + password)
        else:
            print("注册失败")
        time.sleep(1)
