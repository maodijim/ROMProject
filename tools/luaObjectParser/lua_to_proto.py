import argparse
import json
import re

messages = {}
enum = {}
type_map = {
    "2": "float",
    "4": "uint64",
    "5": "int32",
    "8": "bool",
    "9": "string",
    "12": "bytes",
    "13": "uint32",
}
kv_pattern = re.compile(r'(.*?)\s*=\s*(.*)')


def parse_enum(key, val):
    if is_enum_descriptor(val):
        enum[key] = {}
    elif is_enum_val_descriptor(val):
        ks = key.split("_")
        enum_name = ks[0]
        val_start_point = 1
        if enum.get(enum_name, None) is None:
            enum_name = ks[0] + "_" + ks[1]
            val_start_point += 1
        enum_val = "_".join(ks[val_start_point:-1])
        cur_vals = enum.get(enum_name, dict())
        cur_vals[enum_val] = dict()
    elif is_enum_field_val(key):
        split_dot_parts = key.split(".")
        split_parts = split_dot_parts[0].split("_")
        field_name = split_dot_parts[-1]
        enum_name = split_parts[0]
        val_start_point = 1
        if enum.get(enum_name, None) is None:
            enum_name = split_parts[0] + "_" + split_parts[1]
            val_start_point += 1
        enum_val = "_".join(split_parts[val_start_point:-1])
        enum.get(enum_name, dict()).get(enum_val, dict())[field_name] = val.strip('"')


def parse_message(key, val):
    if is_message_descriptor(val):
        messages[key] = dict()
    elif is_message_field_descriptor(val):
        field_name = "_".join(key.split("_")[1:-1])
        message_name = key.split("_")[0]
        fields = messages.get(message_name, dict()).get("fields", dict())
        if len(fields) == 0:
            messages.get(message_name, dict())["fields"] = {field_name: dict()}
        messages.get(message_name, dict())["fields"][field_name] = dict()
    elif is_message_field(key):
        split_parts = key.split(".")
        if len(split_parts) < 2:
            return
        field_name = split_parts[-1]
        message_name = split_parts[0].split("_")[0]
        messages.get(message_name, dict())[field_name] = val.strip('"')
    elif is_message_field_val(key):
        field_name = "_".join(key.split(".")[0].split("_")[1:-1])
        message_name = key.split("_")[0]
        message_val = key.split(".")[-1]
        val = val.strip('"')
        messages.get(message_name, dict())["fields"][field_name][message_val] = val


def get_enum(enum_name, number):
    for _, v in enum.get(enum_name, dict()).items():
        if type(v) != dict:
            continue
        if v.get("number") == number:
            return v.get("name")
    return ""


def has_key_value(line):
    if kv_pattern.match(line):
        return True
    return False


def is_enum_descriptor(val):
    if val == "protobuf.EnumDescriptor();":
        return True
    return False


def is_enum_val_descriptor(val):
    if val == "protobuf.EnumValueDescriptor();":
        return True
    return False


def is_message_descriptor(val):
    if val == "protobuf.Descriptor();":
        return True
    return False


def is_message_field_descriptor(val):
    if val == "protobuf.FieldDescriptor();":
        return True
    return False


def is_enum_field(key, val):
    if re.match(r'.*?.(name|full_name|values)', key) and "ENUM" not in key and "FIELD" not in key:
        try:
            split_parts = key.split(".")
            name = split_parts[0]
            if len(split_parts) < 2:
                return False
            field = split_parts[1]
            if enum.get(name, dict()):
                enum[name][field] = val.strip('"')
                return True
        except Exception as e:
            print(e)
            print("parsing {} {} failed".format(key, val))
    return False


def is_enum_field_val(key):
    if re.match(r'.*?_ENUM.(name|index|number)', key):
        return True
    return False


def is_message_field(key):
    if re.match(r'.*?.(name|full_name|values)', key) and "ENUM" not in key and "FIELD" not in key:
        return True
    return False


def is_message_field_val(key):
    if re.match(r'.*?_FIELD.(name|full_name|number|index|label|has_default_value|default_value|type|cpp_type|message_type|enum_type)', key):
        return True
    return False


def parse_lua(lua):
    line_num = 1
    for line in lua.splitlines():
        if has_key_value(line):
            key, val = kv_pattern.match(line).groups()
            if is_enum_descriptor(val) or is_enum_val_descriptor(val) or is_enum_field_val(key) or is_enum_field(key, val):
                parse_enum(key, val)
            elif is_message_field_descriptor(val) or is_message_descriptor(val) or is_message_field(key) or is_message_field_val(key):
                parse_message(key, val)
        line_num += 1
    print("parsed {} lines".format(line_num))

    # print(json.dumps(messages, indent=4))
    # print(json.dumps(enum, indent=4))


def print_proto():
    result = """syntax = "proto2";

package Cmd;

"""
    for message_name, message in messages.items():
        result += "message {} {{\n".format(message["name"])
        for field_name, field in message.items():
            if field_name != "fields":
                continue
            for _, field_val in field.items():
                if field_val.get("default_value") == "{}":
                    msg_line = "\trepeated "
                else:
                    msg_line = "\toptional "
                if field_val.get("message_type", None):
                    msg_type = field_val.get("message_type")
                    if messages.get(msg_type, None) is None:
                        msg_line += msg_type + "_replace_this_with_imported_proto"
                    else:
                        msg_line += messages.get(msg_type).get("name")
                elif field_val.get("enum_type", None):
                    msg_line += enum.get(field_val["enum_type"], dict()).get("name") or field_val["enum_type"] + "_replace_this_with_imported_proto"
                else:
                    msg_line += type_map[field_val["type"]]
                msg_line += " {} = {}".format(field_val["name"], field_val["number"])
                if field_val.get("has_default_value", None) == "true":
                    if field_val.get("enum_type", None):
                        name = get_enum(field_val["enum_type"], field_val["default_value"])
                        if name == "":
                            msg_line += " [default = {}_replace_this_with_imported_proto_{}]".format(
                                field_val["enum_type"],
                                field_val["default_value"],
                            )
                        else:
                            msg_line += " [default = {}]".format(name)
                    else:
                        msg_line += " [default = {}]".format(field_val["default_value"])
                msg_line += ";"
                result += msg_line + "\n"
        result += "}\n\n"

    for enum_name, enum_val in enum.items():
        if enum_val.get("name", None) is None:
            continue
        result +=("enum {} {{".format(enum_val["name"]))
        for enum_val_name, enum_val_val in enum_val.items():
            if type(enum_val_val) == dict:
                if enum_val_val.get("name", None) is None or enum_val_val.get("number", None) is None:
                    continue
                result += "\n\t{} = {};".format(enum_val_val["name"], enum_val_val["number"])
        result += "\n}\n\n"
    print(result)


if __name__ == '__main__':
    # read lua proto file
    parser = argparse.ArgumentParser()
    parser.add_argument("--path", help="file path to the lua proto file")
    args = parser.parse_args()
    with open(args.path, 'r') as f:
        lua = f.read()
    # parse lua proto file
    parse_lua(lua)
    print_proto()
