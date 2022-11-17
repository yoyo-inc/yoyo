{{define "email.to.html"  }}
{{- if gt (len .Alerts.Firing) 0 -}}
{{- range .Alerts -}}
告警类型: {{ .Labels.alertname }} <br>
告警主题: {{ .Annotations.summary }} <br>
告警级别: {{ }} <br>
触发时间: {{ .StartsAt.Local.Format "2006-01-02 15:04:05" }} <br>
{{- end }}
{{- end }}
{{end  }}
