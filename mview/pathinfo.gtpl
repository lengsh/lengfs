{{ template "header" . }}
<div>
<h1>Hello,lengfs !</h1>
</div>
{{range .}}
{{.}}
{{end}}
<BR>
{{ template "footer" . }}
