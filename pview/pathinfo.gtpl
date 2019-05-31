{{ template "header" . }}
<div>
<h1>Hello,lengfs !</h1>
</div>
{{range .}}
{{.}}<BR>
{{end}}
<BR>
{{ template "footer" . }}
