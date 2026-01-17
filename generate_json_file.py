import json
import random
import os
import string
from enum import Enum, auto

class FieldType(Enum):
    STRING_PAIR = 0
    STRING_INT = 1
    STRING_FLOAT = 2
    STRING_DICT = 3
    STRING_LIST = 4

def generate_random_string():
    return ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(10))

def generate_field(field_type: FieldType):
    key = generate_random_string()
    
    match field_type:
        case FieldType.STRING_PAIR:
            value = generate_random_string()
        case FieldType.STRING_INT:
            value = random.randint(0, 100)
        case FieldType.STRING_FLOAT:
            value = round(random.uniform(0, 100), 2)
        case FieldType.STRING_DICT:
            value = {}
        case FieldType.STRING_LIST:
            value = []
        case _:
            raise ValueError(f"Unknown field type: {field_type}")
    return (key, value)

def generate_random_field(max_type_span, json_body):
    for k in range(max_type_span):
        pair = generate_field(FieldType(k))
        json_body[pair[0]] = pair[1]

def generate_new_random_field(max_type_span):
    json_body = {}
    for k in range(max_type_span):
        pair = generate_field(FieldType(k))
        json_body[pair[0]] = pair[1]
    return json_body

def generate_new_random_json_array(max_array_len, max_type_span):
    json_array = []
    for i in range(max_array_len):
        json_array.append(generate_new_random_field(max_type_span))
    return json_array

def generate_json_string_recursion(max_level, max_field_span, max_array_len, max_type_span, random_noise):
    if max_level == 0 and random_noise == 0:
        return {}
    elif max_level == 0 and random_noise == 1:
        return []
    elif max_level == 1 and random_noise == 0:
        return generate_new_random_field(max_type_span)
    elif max_level == 1 and random_noise == 1:
        return generate_new_random_json_array(max_array_len, max_type_span)

    if random_noise == 0:
        return {generate_random_string(): generate_json_string_recursion(max_level - 1, max_field_span, max_array_len, max_type_span, random.randint(0, 1))}
    else:
        return [generate_json_string_recursion(max_level - 1, max_field_span, max_array_len, max_type_span, random.randint(0, 1)) for i in range(max_array_len)]

def main():
    json_body = generate_json_string_recursion(3, 5, 5, 5, random.randint(0,1))
    with open("output.json", "w") as f:
        json.dump(json_body, f, indent=4)


if __name__ == "__main__":
    main()
