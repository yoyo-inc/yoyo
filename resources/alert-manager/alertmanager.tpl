global:
{{ if .EmailEnable}}
  smtp_smarthost: {{ .SmtpServer }}
  smtp_from: {{ .SmtpSender }}
  smtp_auth_username: {{ .SmtpAuthUser }}
  smtp_auth_password: {{ .SmtpAuthPassword }}
  smtp_require_tls: false
{{ end }}

templates:
  - '/etc/alertmanager/template/*.tmpl'
route:
  receiver: "default-receiver"
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1d
  group_by: ["alertname", "id"]
  routes:
{{ if .WebHookEnable }}
    - receiver: "web.hook"
      match_re:
        severity: warning|critical
      continue: true
{{ end }}
{{ if .EmailEnable}}
    - receiver: "email"
      match_re:
        severity: critical
      continue: true
{{ end }}

receivers:
  - name: "web.hook"
    webhook_configs:
      - url: "http://localhost:8080/api/alert/webhook"
{{ if .EmailEnable}}
  - name: "email"
    email_configs:
      {{range .SmtpReceivers}}
      - to: "{{.Email}}"
        html: {{ `'{{template "email.to.html" .}}'` }}
        headers: { Subject: "YOYO系统" }
      {{end}}
{{ end }}

inhibit_rules:
  - target_matchers:
      - severity = "warning"
    source_matchers:
      - severity = "critical"
    equal: ["alertname", "instance", "id"]