INSERT INTO bookmarks (id, url, title) VALUES 
(gen_random_uuid(), 'https://go.dev', 'Go Programming Language'),
(gen_random_uuid(), 'https://www.postgresql.org', 'PostgreSQL Official'),
(gen_random_uuid(), 'https://github.com', 'GitHub'),
(gen_random_uuid(), 'https://stackoverflow.com', 'Stack Overflow'),
(gen_random_uuid(), 'https://pkg.go.dev', 'Go Packages'),
(gen_random_uuid(), 'https://hub.docker.com', 'Docker Hub'),
(gen_random_uuid(), 'https://12factor.net', 'The Twelve-Factor App'),
(gen_random_uuid(), 'https://golangweekly.com', 'Go Weekly Newsletter'),
(gen_random_uuid(), 'https://sqlformat.org', 'SQL Formatter'),
(gen_random_uuid(), 'https://json.org', 'JSON Introduction')
ON CONFLICT (id) DO NOTHING;
