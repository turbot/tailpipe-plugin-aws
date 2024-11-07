import random
from datetime import datetime, timedelta
import ipaddress
import uuid

def generate_alb_logs(num_lines=10000, start_date=datetime(2024, 11, 1)):
    # Define fixed ALB names for better analysis
    alb_names = [
        "prod-web-alb",
        "prod-api-alb",
        "staging-alb"
    ]
    
    # Target groups associated with specific ALBs
    target_groups = {
        "prod-web-alb": [
            "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/prod-web-front/73e2d6bc24d8",
            "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/prod-web-back/92a1c4de56f7"
        ],
        "prod-api-alb": [
            "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/prod-api-v1/81b3e7cf92a1",
            "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/prod-api-v2/45d9f8g67h2"
        ],
        "staging-alb": [
            "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/staging-web/34f5g6h789a",
            "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/staging-api/12c3d4e567f"
        ]
    }

    # URLs mapped to specific ALBs
    urls = {
        "prod-web-alb": [
            '/assets/main.css',
            '/images/logo.png',
            '/cart',
            '/checkout',
            '/products'
        ],
        "prod-api-alb": [
            '/api/v1/products',
            '/api/v1/orders',
            '/api/v2/users',
            '/api/health'
        ],
        "staging-alb": [
            '/api/v1/products',
            '/web/test',
            '/staging/health'
        ]
    }
    
    # Common user agents
    user_agents = [
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
        'Mozilla/5.0 (iPhone; CPU iPhone OS 16_6_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1',
        'Python-urllib/3.8',
        'Apache-HttpClient/4.5.13 (Java/11.0.19)',
        'curl/7.64.1'
    ]

    # Scanner user agents (for threat simulation)
    scanner_agents = [
        'Nmap Scripting Engine',
        'zgrab/0.x',
        'Nuclei - Open-source project (github.com/projectdiscovery/nuclei)',
        'Qualys SSL Assessment Scanner',
        'WhatWeb/0.5.5'
    ]

    # Suspicious URLs that scanners might try
    suspicious_urls = [
        '/.env',
        '/wp-admin',
        '/phpinfo.php',
        '/admin/console',
        '/actuator/env',
        '/.git/config'
    ]

    # SSL/TLS configurations
    ssl_ciphers = [
        'ECDHE-RSA-AES128-GCM-SHA256',
        'ECDHE-RSA-AES256-GCM-SHA384',
        'TLS_AES_128_GCM_SHA256'
    ]
    tls_versions = ['TLSv1.2', 'TLSv1.3']

    logs = []
    current_time = start_date

    for _ in range(num_lines):
        # Select ALB and related configurations
        alb_name = random.choice(alb_names)
        is_scanner = random.random() < 0.1
        
        # Basic request details
        http_type = 'https' if alb_name.startswith('prod') else random.choice(['http', 'https'])
        timestamp = current_time.strftime('%Y-%m-%dT%H:%M:%S.%fZ')
        
        # Client details
        client_ip = str(ipaddress.IPv4Address(random.randint(0, 2**32 - 1)))
        client_port = random.randint(10000, 65000)
        
        # Target details - production uses private IPs, staging might use public
        if 'prod' in alb_name:
            target_ip = f"10.0.{random.randint(1,255)}.{random.randint(1,255)}"
        else:
            target_ip = str(ipaddress.IPv4Address(random.randint(0, 2**32 - 1)))
        target_port = 443 if alb_name.startswith('prod') else random.choice([80, 443, 8080])

        # Processing times vary by environment
        time_multiplier = 1.0 if alb_name.startswith('prod') else 1.5
        request_processing = round(random.uniform(0, 0.1) * time_multiplier, 6)
        target_processing = round(random.uniform(0, 0.5) * time_multiplier, 6)
        response_processing = round(random.uniform(0, 0.1) * time_multiplier, 6)

        # Choose URLs and status codes
        if is_scanner:
            url = random.choice(suspicious_urls)
            status_code = random.choice([403, 404, 404, 404, 500])
            user_agent = random.choice(scanner_agents)
        else:
            url = random.choice(urls[alb_name])
            status_code = random.choice([200, 200, 200, 200, 301, 302, 404])
            user_agent = random.choice(user_agents)

        # Request details
        method = 'GET' if not is_scanner else random.choice(['GET', 'POST', 'PUT', 'DELETE'])
        request = f"{method} {url} HTTP/1.1"
        
        # Generate trace ID
        trace_id = f"Root=1-{hex(int(current_time.timestamp()))[2:]}-{uuid.uuid4().hex[:24]}"
        
        # Response size varies by environment and status
        received_bytes = random.randint(0, 1000)
        sent_bytes = 0 if status_code == 304 else random.randint(200, 15000)

        # Domain names based on environment
        domain = f"{alb_name}.example.com" if not is_scanner else "-"
        
        # Select appropriate target group
        target_group = random.choice(target_groups[alb_name])

        log_line = f"{http_type} {timestamp} {alb_name} {client_ip}:{client_port} {target_ip}:{target_port} {request_processing:.6f} {target_processing:.6f} {response_processing:.6f} {status_code} {status_code} {received_bytes} {sent_bytes} \"{request}\" \"{user_agent}\" {random.choice(ssl_ciphers)} {random.choice(tls_versions)} {target_group} \"{trace_id}\" \"{domain}\" \"-\" {random.randint(0,10)} {timestamp} \"forward\" \"-\" \"-\" \"{target_ip}:{target_port}\" \"{status_code}\" \"-\" \"-\""
        
        logs.append(log_line)
        
        # Increment time with some randomness but ensure even distribution
        current_time += timedelta(seconds=random.uniform(0.1, 2))

    return logs

# Generate and write logs
if __name__ == "__main__":
    logs = generate_alb_logs()
    with open('alb-test.log', 'w') as f:
        for log in logs:
            f.write(log + '\n')