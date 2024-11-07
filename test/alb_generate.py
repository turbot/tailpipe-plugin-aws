import random
from datetime import datetime, timedelta
import ipaddress
import uuid

def generate_alb_logs(num_lines=10000, start_date=datetime(2024, 11, 1)):
    # ALB details
    alb_names = [
        f"app-my-load-balancer-{uuid.uuid4().hex[:16]}" for _ in range(3)
    ]
    
    target_groups = [
        f"arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup-app-{i}-73e2d6bc24d8a067" 
        for i in range(1, 4)
    ]
    
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

    # URLs including both normal and suspicious paths
    normal_urls = [
        '/api/v1/products',
        '/users/login',
        '/assets/main.css',
        '/images/logo.png',
        '/cart',
        '/checkout',
        '/health'
    ]
    
    suspicious_urls = [
        '/.env',
        '/wp-admin',
        '/phpinfo.php',
        '/admin/console',
        '/actuator/env',
        '/.git/config'
    ]

    # SSL Ciphers
    ssl_ciphers = [
        'ECDHE-RSA-AES128-GCM-SHA256',
        'ECDHE-RSA-AES256-GCM-SHA384',
        'TLS_AES_128_GCM_SHA256'
    ]

    # TLS versions
    tls_versions = ['TLSv1.2', 'TLSv1.3']

    logs = []
    current_time = start_date

    for _ in range(num_lines):
        # Decide if this is a scanner (10% chance)
        is_scanner = random.random() < 0.1
        
        # Basic request details
        http_type = random.choice(['http', 'https', 'h2'])
        timestamp = current_time.strftime('%Y-%m-%dT%H:%M:%S.%fZ')
        alb_name = random.choice(alb_names)
        
        # Client details
        client_ip = str(ipaddress.IPv4Address(random.randint(0, 2**32 - 1)))
        client_port = random.randint(10000, 65000)
        
        # Target details
        target_ip = f"10.0.{random.randint(1,255)}.{random.randint(1,255)}"
        target_port = random.choice([80, 443, 8080, 8443])

        # Processing times (in seconds)
        request_processing = round(random.uniform(0, 0.1), 6)
        target_processing = round(random.uniform(0, 0.5), 6)
        response_processing = round(random.uniform(0, 0.1), 6)

        # Choose URLs and status codes based on scanner or normal traffic
        if is_scanner:
            url = random.choice(suspicious_urls)
            status_code = random.choice([200, 403, 404, 404, 404, 500])
            user_agent = random.choice(scanner_agents)
        else:
            url = random.choice(normal_urls)
            status_code = random.choice([200, 200, 200, 200, 301, 302, 404])
            user_agent = random.choice(user_agents)

        # Request details
        method = random.choice(['GET', 'POST', 'PUT', 'DELETE']) if is_scanner else random.choice(['GET', 'GET', 'GET', 'POST'])
        request = f"{method} {url} HTTP/1.1"
        
        # Generate trace ID
        trace_id = f"Root=1-{hex(int(current_time.timestamp()))[2:]}-{uuid.uuid4().hex[:24]}"
        
        # Size of response (varies by status code)
        received_bytes = random.randint(0, 1000)
        sent_bytes = 0 if status_code == 304 else random.randint(200, 15000)

        # Domain name (sometimes missing for scanners)
        domain = "-" if is_scanner and random.random() < 0.5 else "example.com"
        
        # Build log line
        log_line = (
            f'{http_type} {timestamp} {alb_name} '
            f'{client_ip}:{client_port} {target_ip}:{target_port} '
            f'{request_processing:.6f} {target_processing:.6f} {response_processing:.6f} '
            f'{status_code} {status_code} {received_bytes} {sent_bytes} '
            f'"{request}" "{user_agent}" {random.choice(ssl_ciphers)} {random.choice(tls_versions)} '
            f'{random.choice(target_groups)} "{trace_id}" "{domain}" "-" '
            f'{random.randint(0,10)} {timestamp} "forward" "-" "-" "{target_ip}:{target_port}" '
            f'"{status_code}" "-" "-"'
        )
        
        logs.append(log_line)
        
        # Increment time randomly between 1 and 10 seconds
        current_time += timedelta(seconds=random.randint(1, 10))

    return logs

# Generate and write logs
if __name__ == "__main__":
    logs = generate_alb_logs()
    with open('alb-test.log', 'w') as f:
        for log in logs:
            f.write(log + '\n')