import json
import random
import os
import string

def generate_random_string():
    return ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(10))

def generate_field(type_index):
    if type_index == 0:
        return (generate_random_string(), generate_random_string())
    if type_index == 1:
        return (generate_random_string(), random.randint(0, 100))
    if type_index == 2:
        return (generate_random_string(), round(random.uniform(0, 100), 2))
    if type_index == 3:
        return (generate_random_string(), {})
    if type_index == 4:
        return (generate_random_string(), [])

def generate_random_field(max_type_span, json_body):
    for k in range(max_type_span):
        pair = generate_field(k)
        json_body[pair[0]] = pair[1]

def generate_new_random_field(max_type_span):
    json_body = {}
    for k in range(max_type_span):
        pair = generate_field(k)
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

def generate_json_string(max_level, max_field_span, max_array_len, max_type_span):
    noise = random.randint(0, 1)
    json_body = {}        

    if random.randint(0, 1) == 0:
        json_body = {}
    else:
        json_body = []
    

    for i in range(max_level):
        
        if type(json_body) == dict:
            for j in range(max_field_span):
                generate_random_field(max_type_span, json_body)
        if type(json_body) == list:
            for j in range(max_array_len):
                json_body.append(generate_new_random_field(max_type_span))    
    
    return json_body

def main():
    #json_body = generate_json_string(10, 5, 5, 5)
    json_body = generate_json_string_recursion(3, 5, 5, 5, random.randint(0,1))
    with open("output.json", "w") as f:
        json.dump(json_body, f, indent=4)


if __name__ == "__main__":
    main()
