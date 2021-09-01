import os
import re
import json


def parse_dict(input_str):
    elems = input_str.split(",")
    for elem in elems:
        key_val = elem.split("=")
        key = key_val[0]
        val = key_val[1]


def get_key_val(input_str):
    return re.findall(r"([a-zA-z]*?)=('.*?'|[0-9.]*)", input_str)


def parse_byte_to_dict(in_file):
    with open(in_file, "rb") as f:
        b_file = f.read()
        f.close()

    searches = re.findall(b"{id.*?,}", b_file)

    output = {}
    if searches:
        for search in searches:
            buf_str = search.decode("utf-8", "replace")
            # print(buf_str)
            kvs = get_key_val(buf_str)
            body = {}
            i_id = "0"
            is_first_id = True
            for kv in kvs:
                if kv[0] != "" and kv[1] != "":
                    if kv[0] == "id" and is_first_id:
                        is_first_id = False
                        i_id = str(kv[1])
                    body[kv[0]] = kv[1].replace("'", "")
            output[i_id] = body
    return output

#in_files = r"D:\Downloads\ROAPK\script2\TextAsset\Table_Item.bytes"
in_files = r"D:\Downloads\ROAPK\script2\TextAsset\Table_Buffer.bytes"
# in_files = r"D:\Downloads\ROAPK\script2\TextAsset"
output = {}

if os.path.isdir(in_files):
    output = {}
    file_list = os.listdir(in_files)
    for f in file_list:
        if os.path.isfile(os.path.join(in_files,f)):
            o = parse_byte_to_dict(os.path.join(in_files, f))
            if output is not None:
                output = {**output, **o}
            else:
                output = o
else:
    o = parse_byte_to_dict(in_files)
    output = o

out_file = "items.json"
if os.path.exists(out_file):
    with open(out_file, "r", encoding="utf-8") as f:
        in_file = json.load(f)
    for key, val in in_file.items():
        output[key] = val

with open(out_file, "w", encoding="utf-8") as f:
    print("saving to" + os.getcwd())
    json.dump(output, f, indent=2, ensure_ascii=False, sort_keys=True)
    f.close()
